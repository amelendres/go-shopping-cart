package shop

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestAddProduct(t *testing.T) {
	store := StubCartStore{
		[]Product{},
	}
	server := NewCartServer(&store)

	t.Run("it add Product on POST", func(t *testing.T) {
		product := Product{"Dr. Pepper", 2}

		request := newPostProductRequest(product)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)
		AssertProduct(t, &store, product)
	})

	//Add same product increment qty
	//Add different products
}

func TestGetProducts(t *testing.T) {

	t.Run("it returns the products table as JSON", func(t *testing.T) {
		wantedProducts := []Product{
			{"Te", 32},
			{"Bread", 20},
			{"Coffee", 14},
		}

		store := StubCartStore{wantedProducts}
		server := NewCartServer(&store)

		request := newGetProductsRequest()
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := getProductsFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		assertProducts(t, got, wantedProducts)
		assertContentType(t, response, jsonContentType)

	})
}

func assertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Header().Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.HeaderMap)
	}
}

func getProductsFromResponse(t *testing.T, body io.Reader) []Product {
	t.Helper()
	products, err := NewProducts(body)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Product, '%v'", body, err)
	}

	return products
}

func assertProducts(t *testing.T, got, want []Product) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func newGetProductsRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/products", nil)
	return req
}

func newPostProductRequest(product Product) *http.Request {
	body, _ := json.Marshal(product)
	req, _ := http.NewRequest(http.MethodPost,"/products", bytes.NewBuffer(body))
	return req
}

func assertResponseBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}
