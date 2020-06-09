package shop

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddingProductsAndRetrievingThem(t *testing.T) {
	database, cleanDatabase := createTempFile(t, `[]`)
	defer cleanDatabase()
	store, err := NewFileSystemCartStore(database)

	assertNoError(t, err)

	server := NewCartServer(store)
	product := Product{"Pepsi", 1}

	server.ServeHTTP(httptest.NewRecorder(), newPostProductRequest(product))
	server.ServeHTTP(httptest.NewRecorder(), newPostProductRequest(product))
	server.ServeHTTP(httptest.NewRecorder(), newPostProductRequest(product))

	t.Run("get products", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetProductsRequest())
		assertStatus(t, response.Code, http.StatusOK)

		got := getProductsFromResponse(t, response.Body)
		fmt.Printf("%+v\n", got)
		want := []Product{
			{product.Name, 3},
		}
		assertProducts(t, got, want)
	})
}
