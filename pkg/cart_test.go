package cart_test

import (
	. "github.com/amelendres/go-shopping-cart/pkg"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

var(
	cart *Cart
)

func TestCreateCart(t *testing.T) {

	t.Run("it creates cart successfully", func(t *testing.T) {
		cartID := UUID(uuid.New().String())
		cart = NewCart(cartID, UUID(uuid.New().String()))

		assert.NotNil(t, cart)
		assert.Equal(t, cartID, cart.ID())
		assert.NotEmpty(t, cart.BuyerID().String())
		assert.Equal(t,0, len(cart.Products()))
	})
}

func TestAddProduct(t *testing.T) {
	productID := UUID(uuid.New().String())
	p := NewProduct(productID, "Pants", Price(10), Quantity(2))

	t.Run("Add First Product", func(t *testing.T) {
		cart.AddProduct(p.ID(), p.Name(), p.Price(), p.Qty())

		prod, _ := cart.Product(productID)
		assert.Equal(t, prod, p)
	})

	t.Run("Add same Product increase product units", func(t *testing.T) {
		cart.AddProduct(p.ID(), p.Name(), p.Price(), 2)

		prod, _ := cart.Product(productID)
		assert.Equal(t, 4, prod.Qty().Int())
		assert.Equal(t, 1, len(cart.Products()))

	})

	t.Run("Fails Adding same Product with different price", func(t *testing.T) {
		err := cart.AddProduct(p.ID(), p.Name(), Price(20), p.Qty())

		assert.Error(t, err)
	})

	t.Run("Add other Product", func(t *testing.T) {
		otherProductID := UUID(uuid.New().String())

		cart.AddProduct(otherProductID, "Short", 10, 5)

		prod, _ := cart.Product(otherProductID)
		assert.Equal(t, 5, prod.Qty().Int())
		assert.Equal(t, 2, len(cart.Products()))
	})
}

func TestGetProducts(t *testing.T) {
	t.Run("Get Products", func(t *testing.T) {
		assert.Equal(t, 2, len(cart.Products()))
	})
}

func TestProductNotFound(t *testing.T) {
	t.Run("Get a non existing Product", func(t *testing.T) {
		_, err := cart.Product(UUID(uuid.New().String()))
		assert.Error(t, err)
	})
}
