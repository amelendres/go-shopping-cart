package cart

import (
	"github.com/google/uuid"
	"testing"
)

func TestCartAddProduct(t *testing.T) {
	var productID = uuid.New().String()
	var cart = NewCart(UUID(uuid.New().String()), UUID(uuid.New().String()))

	t.Run("Cart Add First Product", func(t *testing.T) {
		cart.AddProduct(productID, "Pants", 10, 2 )
		AssertCartProduct(t, cart, Product{productID, "Pants", 10, 2})
	})

	t.Run("Cart Add same Product increase product units", func(t *testing.T) {
		cart.AddProduct(productID, "Pants", 10, 2)
		AssertCartProductQty(t, cart, productID, 4)
		AssertCartProductLines(t, cart, 1)
	})

	var otherProductID = uuid.New().String()
	t.Run("Cart Add other Product", func(t *testing.T) {
		cart.AddProduct(otherProductID, "Short", 10, 5)
		AssertCartProductQty(t, cart, otherProductID, 5)
		AssertCartProductLines(t, cart, 2)
	})

	t.Run("Cart Get Products", func(t *testing.T) {
		AssertEquals(t, len(cart.GetProducts()), 2)
	})
}
