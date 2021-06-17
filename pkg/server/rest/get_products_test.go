package rest_test

import (
	. "github.com/amelendres/go-shopping-cart/pkg"
	cart "github.com/amelendres/go-shopping-cart/pkg"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestListingCartProducts(t *testing.T) {
	c := cart.NewCart(UUID(uuid.New().String()), UUID(uuid.New().String()))
	emptyCart := cart.NewCart(UUID(uuid.New().String()), UUID(uuid.New().String()))

	p := NewProduct(UUID(uuid.New().String()), "Pants", Price(7.5), Quantity(1))
	op := NewProduct(UUID(uuid.New().String()), "Socks", Price(0.5), Quantity(10))

	repo, closeRepo := getRepo(t, SQL)
	cartService := newHTTPServer(t, repo)

	//GIVEN
	response := httptest.NewRecorder()
	cartService.ServeHTTP(response, newCreateCartHTTPRequest(newCreateCartReq(c)))
	cartService.ServeHTTP(response, newAddProductHTTPRequest(newAddProductReq(c.ID(), p)))
	cartService.ServeHTTP(response, newAddProductHTTPRequest(newAddProductReq(c.ID(), p)))
	cartService.ServeHTTP(response, newAddProductHTTPRequest(newAddProductReq(c.ID(), op)))

	cartService.ServeHTTP(response, newCreateCartHTTPRequest(newCreateCartReq(emptyCart)))

	var testCases = []TestCase{
		{
			name:   "An existing cart",
			req:    c.ID().String(),
			resp:   2,
			status: 200,
		},
		{
			name:   "A non exisiting cart",
			req:    uuid.New().String(),
			status: 404,
		},
		{
			name:   "An empty cart",
			req:    emptyCart.ID().String(),
			resp:   0,
			status: 200,
		},
		{
			name:   "An invalid request",
			req:    "invalid-cart-id",
			status: 400,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			//WHEN
			response := httptest.NewRecorder()
			cartService.ServeHTTP(response, newGetProductsRequest(testCase.req.(string)))

			//THEN
			assert.Equal(t, testCase.status, response.Code)
			if testCase.resp != nil {
				products, err := newProductResponseFromJSON(response.Body)
				if err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, testCase.resp, len(products))
			}
		})
	}
	t.Cleanup(func() {
		closeRepo()
	})
}
