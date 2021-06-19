package rest_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/amelendres/go-shopping-cart/pkg"
	"github.com/amelendres/go-shopping-cart/pkg/adding"
	"github.com/amelendres/go-shopping-cart/pkg/creating"
	"github.com/amelendres/go-shopping-cart/pkg/listing"
	"github.com/amelendres/go-shopping-cart/pkg/storage/fs"
	"github.com/amelendres/go-shopping-cart/pkg/storage/pgsql"
	_ "github.com/mattn/go-sqlite3"

	"github.com/amelendres/go-shopping-cart/pkg/server/rest"
)

type TestCase struct {
	name   string
	req    interface{}
	resp   interface{}
	status int
}

type dbType int

const (
	TMP dbType = iota
	JSON
	SQL
)

func newHTTPServer(t *testing.T, r cart.Repository) *rest.CartServer {
	t.Helper()
	return rest.NewCartServer(
		logrus.New(),
		creating.NewCartCreator(r),
		adding.NewProductAdder(r),
		listing.NewProductLister(r))
}

func getRepo(t *testing.T, dt dbType) (cart.Repository, func()) {
	t.Helper()
	switch dt {
	case TMP:
		store, closeStore, err := createTempFile(t, `[]`)
		if err != nil {
			log.Fatal(err)
		}

		repo, err := fs.NewCartStore(store)
		if err != nil {
			log.Fatal(err)
		}
		return repo, closeStore

	case JSON:
		repo, closeStore, err := fs.CartStoreFromFile("cart_test.db.json")
		if err != nil {
			log.Fatal(err)
		}
		return repo, closeStore

	case SQL:
		conn, err := NewTestConn(t)
		if err != nil {
			log.Fatal(err)
		}
		resetDB(t, conn)
		repo := pgsql.NewCartRepository(conn)
		if err != nil {
			log.Fatal(err)
		}
		return repo, func() {
			conn.Close()
		}
	default:
		log.Fatal("Unsupported DB Type")
		return nil, nil
	}
}

func createTempFile(t *testing.T, initialData string) (*os.File, func(), error) {
	t.Helper()

	tf, err := ioutil.TempFile("", "db")

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	_, err = tf.Write([]byte(initialData))
	if err != nil {
		return nil, nil, err
	}

	removeFile := func() {
		os.Remove(tf.Name())
	}

	return tf, removeFile, nil
}

func resetDB(t *testing.T, db *sql.DB) error {
	t.Helper()

	q := `	DROP TABLE IF EXISTS product_line;
			DROP TABLE IF EXISTS cart;

			CREATE TABLE IF NOT EXISTS cart
			(
				id       UUID NOT NULL,
				buyer_id UUID NOT NULL,
				CONSTRAINT pk_cart PRIMARY KEY (id)
			);
			
			CREATE TABLE IF NOT EXISTS product_line
			(
				id      INT auto_increment,
				product_id UUID  NOT NULL,
				cart_id UUID  NOT NULL,
				name    text  NOT NULL,
				price   FLOAT NOT NULL,
				qty     INT   NOT NULL,
				CONSTRAINT pk_product_line PRIMARY KEY (id),
				CONSTRAINT fk_product_cart FOREIGN KEY (cart_id) REFERENCES cart (id)
			);
			`
	_, err := db.Exec(q)
	if err != nil {
		return err
	}
	return nil
}

func NewTestConn(t *testing.T) (*sql.DB, error) {
	t.Helper()

	return sql.Open("sqlite3", "db_test.sqlite")
}

func newGetProductsRequest(cartID string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/carts/%s/products", cartID), nil)
	return req
}

func newCreateCartHTTPRequest(cc creating.CreateCartReq) *http.Request {
	body, err := json.Marshal(cc)
	if err != nil {
		log.Fatalln(err)
	}
	httpReq, _ := http.NewRequest(http.MethodPost,"/carts", bytes.NewBuffer(body))
	//fmt.Println(httpReq)
	return httpReq
}

func newAddProductHTTPRequest(ap adding.AddProductReq) *http.Request {
	body, err := json.Marshal(ap)
	if err != nil {
		log.Fatalln(err)
	}
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/carts/%s/products", ap.CartID), bytes.NewBuffer(body))
	//fmt.Println(req)

	return req
}

func newProductResponseFromJSON(rdr io.Reader) ([]listing.ProductResponse, error) {
	var products []listing.ProductResponse
	err := json.NewDecoder(rdr).Decode(&products)

	if err != nil {
		err = fmt.Errorf("problem parsing ProductResponse, %v", err)
	}

	return products, err
}
