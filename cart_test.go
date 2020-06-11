package shopping

import (
	"testing"
)

func TestCartAddProduct(t *testing.T) {
	var cart = NewCart("uuid", "uuid")

	t.Run("Cart Add First Product", func(t *testing.T) {
		product := Product{"uuid", "Pants", 10, 2 }
		cart.AddProduct(product)
		AssertCartProduct(t, cart, product)
	})

	t.Run("Cart Add same Product increase product units", func(t *testing.T) {
		product := Product{"uuid", "Pants", 10, 2 }
		cart.AddProduct(product)
		AssertCartProductQty(t, cart, product, 4)
		AssertCartProductLines(t, cart, 1)
	})

	t.Run("Cart Add other Product", func(t *testing.T) {
		product := Product{"uuid", "Short", 10, 5 }
		cart.AddProduct(product)
		AssertCartProductQty(t, cart, product, 5)
		AssertCartProductLines(t, cart, 2)
	})

	t.Run("Cart Get Products", func(t *testing.T) {
		AssertEquals(t, len(cart.GetProducts()), 2)
	})
}
