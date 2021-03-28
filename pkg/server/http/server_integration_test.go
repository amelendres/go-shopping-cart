package http

import (
	cart "github.com/amelendres/go-shopping-cart/pkg"
	"github.com/amelendres/go-shopping-cart/pkg/storage/fs"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAddingProductsAndRetrievingThem(t *testing.T) {
	database, cleanDatabase := createTempFile(t, `[]`)
	defer cleanDatabase()
	store, err := fs.NewCartStore(database)

	cart.AssertNoError(t, err)

	cs := NewCartServer(store)
	product := cart.Product{"uuid", "Pepsi", 0.7,  1}
	cartID := uuid.New().String()

	cs.ServeHTTP(httptest.NewRecorder(), newPostProductRequest(cartID, product))
	cs.ServeHTTP(httptest.NewRecorder(), newPostProductRequest(cartID, product))
	cs.ServeHTTP(httptest.NewRecorder(), newPostProductRequest(cartID, product))

	t.Run("get Products", func(t *testing.T) {
		response := httptest.NewRecorder()
		cs.ServeHTTP(response, newGetProductsRequest(cartID))
		assertStatus(t, response.Code, http.StatusOK)

		got := getProductsFromResponse(t, response.Body)
		want := []cart.Product{
			{"uuid", "Pepsi", 0.7,  3},
		}
		assertProducts(t, got, want)
	})

	t.Run("get Products Bad Request", func(t *testing.T) {
		response := httptest.NewRecorder()
		cs.ServeHTTP(response, newGetProductsRequest(uuid.New().String()))
		assertStatus(t, response.Code, http.StatusBadRequest)
	})
}

func createTempFile(t *testing.T, initialData string) (*os.File, func()){
	t.Helper()

	tmpfile, err := ioutil.TempFile("", "db")

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tmpfile.Write([]byte(initialData))

	removeFile := func() {
		os.Remove(tmpfile.Name())
	}

	return tmpfile, removeFile
}
