package server

import (
	"github.com/amelendres/go-shopping-cart"
	"github.com/amelendres/go-shopping-cart/fs"
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
	store, err := fs.NewFileSystemCartStore(database)

	shopping.AssertNoError(t, err)

	cs := NewCartServer(store)
	product := shopping.Product{"uuid", "Pepsi", 0.7,  1}
	cartID := uuid.New().String()

	cs.ServeHTTP(httptest.NewRecorder(), newPostProductRequest(cartID, product))
	cs.ServeHTTP(httptest.NewRecorder(), newPostProductRequest(cartID, product))
	cs.ServeHTTP(httptest.NewRecorder(), newPostProductRequest(cartID, product))

	t.Run("get Products", func(t *testing.T) {
		response := httptest.NewRecorder()
		cs.ServeHTTP(response, newGetProductsRequest(cartID))
		assertStatus(t, response.Code, http.StatusOK)

		got := getProductsFromResponse(t, response.Body)
		want := []shopping.Product{
			{"uuid", "Pepsi", 0.7,  3},
		}
		assertProducts(t, got, want)
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
