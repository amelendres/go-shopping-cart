package shopping

import (
	"github.com/google/uuid"
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
	product := Product{"uuid", "Pepsi", 0.7,  1}
	cartID := uuid.New().String()

	server.ServeHTTP(httptest.NewRecorder(), newPostProductRequest(cartID, product))
	server.ServeHTTP(httptest.NewRecorder(), newPostProductRequest(cartID, product))
	server.ServeHTTP(httptest.NewRecorder(), newPostProductRequest(cartID, product))

	t.Run("get Products", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newGetProductsRequest(cartID))
		assertStatus(t, response.Code, http.StatusOK)

		got := getProductsFromResponse(t, response.Body)
		want := []Product{
			{"uuid", "Pepsi", 0.7,  3},
		}
		assertProducts(t, got, want)
	})
}
