package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	cart "github.com/amelendres/go-shopping-cart/pkg"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestAddProduct(t *testing.T) {

	store := cart.NewStubCartStore()
	cs := NewCartServer(store)

	t.Run("it add Product on POST", func(t *testing.T) {
		product := cart.Product{"uuid", "Dr. Pepper", 0.5, 2}

		request := newPostProductRequest(uuid.New().String(), product)
		response := httptest.NewRecorder()
		cs.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)
	})
}

func TestGetProducts(t *testing.T) {

	t.Run("it returns the Products as JSON", func(t *testing.T) {
		cartID := uuid.New().String()
		wantedCarts := []cart.Cart{
			{
				"cartId",
				"buyerId",
				[]cart.Product{
					{"uuid1", "Te", 1, 32},
					{"uuid2", "Bread", 0.3, 20},
					{"uuid3", "Coffee", 1, 14},
				},
			},
		}

		store := cart.StubCartStore{wantedCarts}
		server := NewCartServer(&store)

		request := newGetProductsRequest(cartID)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := getProductsFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		assertProducts(t, got, wantedCarts[0].Products)
		assertContentType(t, response, jsonContentType)

	})
}

func assertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Header().Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.HeaderMap)
	}
}

func getProductsFromResponse(t *testing.T, body io.Reader) []cart.Product {
	t.Helper()
	products, err := newProductsFromJSON(body)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Product, '%v'", body, err)
	}

	return products
}

func assertProducts(t *testing.T, got, want []cart.Product) {
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

func newGetProductsRequest(cartID string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/carts/%s/products", cartID), nil)
	return req
}

func newPostProductRequest(cartID string, product cart.Product) *http.Request {
	body, err := json.Marshal(product)
	if err != nil {
		log.Fatalln(err)
	}
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/carts/%s/products", cartID), bytes.NewBuffer(body))
	return req
}

func newProductsFromJSON(rdr io.Reader) ([]cart.Product, error) {
	var products []cart.Product
	err := json.NewDecoder(rdr).Decode(&products)

	if err != nil {
		err = fmt.Errorf("problem parsing Products, %v", err)
	}

	return products, err
}
