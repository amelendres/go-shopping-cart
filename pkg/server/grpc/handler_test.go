package grpc

import (
	"context"
	"database/sql"
	cart "github.com/amelendres/go-shopping-cart/pkg"
	"github.com/amelendres/go-shopping-cart/pkg/adding"
	"github.com/amelendres/go-shopping-cart/pkg/creating"
	"github.com/amelendres/go-shopping-cart/pkg/listing"
	"github.com/amelendres/go-shopping-cart/pkg/storage/fs"
	"github.com/amelendres/go-shopping-cart/pkg/storage/pgsql"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	. "github.com/onsi/gomega"

	cartgrpc "github.com/amelendres/go-shopping-cart/proto"
)

type dbType int
const (
	TMP dbType = iota
	JSON
	SQL
)

type serviceType int
const (
	CREATE_CART serviceType = iota
	ADD_PRODUCT
	LIST_CART_PRODUCTS
)

func TestCreatingCart(t *testing.T) {
	c := &cartgrpc.Cart{
		Id:      uuid.New().String(),
		BuyerId: uuid.New().String(),
	}
	repo, closeRepo := getRepo(t, SQL)
	cartService := newCartGrpcServer(t, repo)

	testCases := []struct {
		name        string
		req         *cartgrpc.CreateCartReq
		message     string
		expectedErr bool
	}{
		{
			name:        "A new valid cart",
			req:         &cartgrpc.CreateCartReq{Cart: c},
			message:     "",
			expectedErr: false,
		},
		{
			name:        "An existing cart",
			req:         &cartgrpc.CreateCartReq{Cart: c},
			expectedErr: true,
		},
		{
			name:        "An invalid request",
			req:         &cartgrpc.CreateCartReq{},
			expectedErr: true,
		},
		{
			name:        "nil request",
			req:         nil,
			expectedErr: true,
		},
	}

	//Run test scenarios
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			g := NewGomegaWithT(t)

			ctx := context.Background()

			//WHEN
			response, err := cartService.Create(ctx, testCase.req)
			t.Log("Got : ", response)

			//THEN
			if testCase.expectedErr {
				g.Expect(err).ToNot(BeNil(), "Error shouldn't be nil")
			} else {
				g.Expect(response).To(Equal(testCase.message))
			}
		})
	}
	t.Cleanup(func() {
		closeRepo()
	})
}

func TestAddingCartProduct(t *testing.T) {
	c := &cartgrpc.Cart{
		Id:      uuid.New().String(),
		BuyerId: uuid.New().String(),
	}
	oc := &cartgrpc.Cart{
		Id:      uuid.New().String(),
		BuyerId: uuid.New().String(),
	}
	p := &cartgrpc.Product{
		Id:        uuid.New().String(),
		Name:      "Pants",
		UnitPrice: 7.5,
		Units:     1,
	}
	p1 := &cartgrpc.Product{
		Id:        p.Id,
		Name:      p.Name,
		UnitPrice: 1,
		Units:     1,
	}

	op := &cartgrpc.Product{
		Id:        uuid.New().String(),
		Name:      "Socks",
		UnitPrice: 0.5,
		Units:     10,
	}

	repo, closeRepo := getRepo(t, SQL)
	cartService := newCartGrpcServer(t, repo)
	ctx := context.Background()
	cartService.Create(ctx, &cartgrpc.CreateCartReq{Cart: c})
	cartService.Create(ctx, &cartgrpc.CreateCartReq{Cart: oc})

	testCases := []struct {
		name        string
		req         *cartgrpc.AddProductReq
		message     string
		expectedErr bool
	}{
		{
			name:        "A new valid product",
			req:         &cartgrpc.AddProductReq{CartId: c.Id, Product: p},
			message:     "",
			expectedErr: false,
		},
		{
			name:        "Same product",
			req:         &cartgrpc.AddProductReq{CartId: c.Id, Product: p},
			message:     "",
			expectedErr: false,
		},
		{
			name:        "Same product with different price",
			req:         &cartgrpc.AddProductReq{CartId: c.Id, Product: p1},
			expectedErr: true,
		},
		{
			name:        "Same product to another cart",
			req:         &cartgrpc.AddProductReq{CartId: oc.Id, Product: p},
			message:     "",
			expectedErr: false,
		},
		{
			name:        "Another valid product",
			req:         &cartgrpc.AddProductReq{CartId: c.Id, Product: op},
			message:     "",
			expectedErr: false,
		},
		{
			name:        "A non existing cart",
			req:         &cartgrpc.AddProductReq{CartId: uuid.New().String(), Product: p},
			expectedErr: true,
		},
		{
			name:        "An invalid request",
			req:         &cartgrpc.AddProductReq{},
			expectedErr: true,
		},
		{
			name:        "An invalid request",
			req:         &cartgrpc.AddProductReq{Product: p},
			expectedErr: true,
		},
		{
			name:        "nil request",
			req:         nil,
			expectedErr: true,
		},
	}

	//Run test scenarios
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			g := NewGomegaWithT(t)

			//WHEN
			response, err := cartService.Add(ctx, testCase.req)
			t.Log("Got : ", response)

			//THEN
			if testCase.expectedErr {
				g.Expect(err).ToNot(BeNil(), "Error shouldn't be nil")
			} else {
				g.Expect(response).To(Equal(testCase.message))
			}
		})
	}

	t.Cleanup(func() {
		closeRepo()
	})

}

func TestListingCartProducts(t *testing.T) {
	c := &cartgrpc.Cart{
		Id:      uuid.New().String(),
		BuyerId: uuid.New().String(),
	}
	ec := &cartgrpc.Cart{
		Id:      uuid.New().String(),
		BuyerId: uuid.New().String(),
	}
	p := &cartgrpc.Product{
		Id:        uuid.New().String(),
		Name:      "Pants",
		UnitPrice: 7.5,
		Units:     1,
	}
	op := &cartgrpc.Product{
		Id:        uuid.New().String(),
		Name:      "Socks",
		UnitPrice: 0.5,
		Units:     10,
	}

	repo, closeRepo := getRepo(t, SQL)
	cartService := newCartGrpcServer(t, repo)
	ctx := context.Background()

	//GIVEN
	cartService.Create(ctx, &cartgrpc.CreateCartReq{Cart: c})
	cartService.Add(ctx, &cartgrpc.AddProductReq{CartId: c.Id, Product: p})
	cartService.Add(ctx, &cartgrpc.AddProductReq{CartId: c.Id, Product: p})
	cartService.Add(ctx, &cartgrpc.AddProductReq{CartId: c.Id, Product: op})

	cartService.Create(ctx, &cartgrpc.CreateCartReq{Cart: ec})

	testCases := []struct {
		name        string
		req         *cartgrpc.ListCartReq
		message     int
		expectedErr bool
	}{
		{
			name:        "An existing cart",
			req:         &cartgrpc.ListCartReq{CartId: c.Id},
			message:     2,
			expectedErr: false,
		},
		{
			name:        "A non exisiting cart",
			req:         &cartgrpc.ListCartReq{CartId: uuid.New().String()},
			expectedErr: true,
		},
		{
			name:        "An empty cart",
			req:         &cartgrpc.ListCartReq{CartId: ec.Id},
			message:     0,
			expectedErr: false,
		},
		{
			name:        "An invalid request",
			req:         &cartgrpc.ListCartReq{},
			expectedErr: true,
		},
		{
			name:        "nil request",
			req:         nil,
			expectedErr: true,
		},
	}

	//Run test scenarios
	for _, testCase := range testCases {
	//for _, tc := range testCases {
		//testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			g := NewGomegaWithT(t)

			//WHEN
			response, err := cartService.List(ctx, testCase.req)
			t.Log("Got : ", response)

			//THEN
			if testCase.expectedErr {
				g.Expect(err).ToNot(BeNil(), "Error shouldn't be nil")
			} else {
				g.Expect(len(response.Products)).To(Equal(testCase.message))
			}
		})
	}
	t.Cleanup(func() {
		closeRepo()
	})
}

type TestCase struct {
	name        string
	req         *interface{}
	message     int
	expectedErr bool
}

func newCartGrpcServer(t *testing.T, r cart.Repository) cartgrpc.CartServiceServer {
	t.Helper()
	return NewCartServiceServer(creating.NewCartCreator(r), adding.NewProductAdder(r), listing.NewProductLister(r))
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
		//repo, closeStore, err := fs.CartStoreFromFile("cart_test.db.json")
		repo, _, err := fs.CartStoreFromFile("cart_test.db.json")
		if err != nil {
			log.Fatal(err)
		}
		//return repo, closeStore
		return repo, func() {}

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
