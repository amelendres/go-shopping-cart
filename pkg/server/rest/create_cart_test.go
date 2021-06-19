package rest_test

import (
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"

	. "github.com/amelendres/go-shopping-cart/pkg"
	cart "github.com/amelendres/go-shopping-cart/pkg"
	"github.com/amelendres/go-shopping-cart/pkg/creating"
	"github.com/google/uuid"

	_ "github.com/mattn/go-sqlite3"
)

func TestCreatingCart(t *testing.T) {
	c := cart.NewCart(UUID(uuid.New().String()), UUID(uuid.New().String()))

	repo, closeRepo := getRepo(t, SQL)
	cartService := newHTTPServer(t, repo)

	var testCases = []TestCase{
		{
			name:   "A new valid cart",
			req:    newCreateCartReq(c),
			resp:   "",
			status: 200,
		},
		{
			name:   "An existing cart",
			req:    newCreateCartReq(c),
			status: 409,
		},
		{
			name:   "An invalid request",
			req:    creating.CreateCartReq{},
			status: 400,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			//WHEN
			response := httptest.NewRecorder()
			cartService.ServeHTTP(response, newCreateCartHTTPRequest(testCase.req.(creating.CreateCartReq)))

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

func newCreateCartReq(c *Cart) creating.CreateCartReq {
	return creating.CreateCartReq{
		CartID:  c.ID().String(),
		BuyerID: c.BuyerID().String(),
	}
}
