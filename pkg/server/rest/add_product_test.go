package rest_test

import (
	. "github.com/amelendres/go-shopping-cart/pkg"
	cart "github.com/amelendres/go-shopping-cart/pkg"
	"github.com/amelendres/go-shopping-cart/pkg/adding"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	_ "github.com/mattn/go-sqlite3"
)

func TestAddingProduct(t *testing.T) {
	c := cart.NewCart(UUID(uuid.New().String()), UUID(uuid.New().String()))
	oc := cart.NewCart(UUID(uuid.New().String()), UUID(uuid.New().String()))

	p := NewProduct(UUID(uuid.New().String()), "Pants", Price(7.5), Quantity(1))
	p1 := NewProduct(p.ID(), p.Name(), Price(1), Quantity(1))
	op := NewProduct(UUID(uuid.New().String()), "Socks", Price(0.5), Quantity(10))

	repo, closeRepo := getRepo(t, SQL)
	cartService := newHTTPServer(t, repo)

	//GIVEN
	cartService.ServeHTTP(httptest.NewRecorder(), newCreateCartHTTPRequest(newCreateCartReq(c)))
	cartService.ServeHTTP(httptest.NewRecorder(), newCreateCartHTTPRequest(newCreateCartReq(oc)))

	var testCases = []TestCase{
		{
			name:   "A new valid product",
			req:    newAddProductReq(c.ID(), p),
			resp:   "",
			status: 200,
		},
		{
			name:   "Same product",
			req:    newAddProductReq(c.ID(), p),
			resp:   "",
			status: 200,
		},
		{
			name: "Same product with different price",
			req:    newAddProductReq(c.ID(), p1),
			status: 409,
		},
		{
			name:   "Same product to another cart",
			req:    newAddProductReq(oc.ID(), p),
			resp:   "",
			status: 200,
		},
		{
			name:   "Another valid product",
			req:    newAddProductReq(c.ID(), op),
			resp:   "",
			status: 200,
		},
		{
			name:   "A non existing cart",
			req:    newAddProductReq(UUID(uuid.New().String()), p),
			status: 404,
		},
		//{
		//	name:   "An invalid request",
		//	req:    newAddProductReq(UUID(uuid.NewCart().String()), &cart.Product{}),
		//	status: 400,
		//},
	}

	//Run test scenarios
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			//WHEN
			response := httptest.NewRecorder()
			cartService.ServeHTTP(response, newAddProductHTTPRequest(testCase.req.(adding.AddProductReq)))

			//THEN
			assert.Equal(t, testCase.status, response.Code)
			if testCase.resp != nil {
				assert.Equal(t, testCase.resp, response.Body.String())
			}
		})
	}

	t.Cleanup(func() {
		closeRepo()
	})

}

func newAddProductReq(cartID UUID, p *Product) adding.AddProductReq {
	return adding.AddProductReq{
		CartID:    cartID.String(),
		ProductID: p.ID().String(),
		Name:      p.Name(),
		Price:     float64(p.Price()),
		Units:     p.Qty().Int(),
	}
}
