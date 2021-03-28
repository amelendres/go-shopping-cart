package cart

import (
	"github.com/google/uuid"
	"testing"
)

func TestCartAddProduct(t *testing.T) {
	var productId = uuid.New().String()
	var cart = NewCart(UUID(uuid.New().String()), UUID(uuid.New().String()))

	t.Run("Cart Add First Product", func(t *testing.T) {
		cart.AddProduct(productId, "Pants", 10, 2 )
		AssertCartProduct(t, cart, Product{productId, "Pants", 10, 2})
	})

	t.Run("Cart Add same Product increase product units", func(t *testing.T) {
		cart.AddProduct(productId, "Pants", 10, 2)
		AssertCartProductQty(t, cart, productId, 4)
		AssertCartProductLines(t, cart, 1)
	})

	var otherProductId = uuid.New().String()
	t.Run("Cart Add other Product", func(t *testing.T) {
		cart.AddProduct(otherProductId, "Short", 10, 5)
		AssertCartProductQty(t, cart, otherProductId, 5)
		AssertCartProductLines(t, cart, 2)
	})

	t.Run("Cart Get Products", func(t *testing.T) {
		AssertEquals(t, len(cart.GetProducts()), 2)
	})
}
