package shopping

import (
	"testing"
)

type StubCartStore struct {
	Carts []Cart
}

func (s *StubCartStore) Save(cart *Cart) error {
	s.Carts = append(s.Carts, *cart)
	return nil
}

func (s *StubCartStore) Get(id string) (*Cart, error) {
	if len(s.Carts) != 1 {
		return nil, ErrUnknownCart
	}
	return &s.Carts[0], nil
}

func AssertProduct(t *testing.T, store *StubCartStore, product Product) {
	t.Helper()

	if len(store.Carts) != 1 {
		t.Fatalf("got %d calls to AddProduct want %d", len(store.Carts), 1)
	}
	if store.Carts[0].Products[0] != product {
		t.Errorf("did not store the correct product got %+v want %+v", store.Carts[0], product)
	}
}

func AssertCartProduct(t *testing.T, cart *Cart, product Product) {
	t.Helper()

	if len(cart.Products) != 1 {
		t.Fatalf("got %d calls to AddProduct want %d", len(cart.Products), 1)
	}

	if cart.Products[0] != product {
		t.Errorf("did not cart the correct product got %+v want %+v", cart.Products[0], product)
	}
}

func AssertCartProductQty(t *testing.T, cart *Cart, product Product, qty Quantity) {
	t.Helper()

	_, prod := cart.Products.Find(product)

	if nil == prod {
		t.Fatalf("Proudct Not Found, got %q  want %d", product.Name, qty)
	}

	if prod.Units != qty {
		t.Errorf("did not cart the correct product Units, got %d want %d", prod.Units, qty)
	}
}

func AssertCartProductLines(t *testing.T, cart *Cart, qty int) {
	t.Helper()
	lines := len(cart.Products)
	if lines != qty{
		t.Errorf("did not cart the correct product lines, got %q want %q", lines, qty)
	}
}


//BASIC ASSERTS
func AssertEquals(t *testing.T, got int, want int) {
	t.Helper()

	if got != want{
		t.Errorf("got %q want %q", got, want)
	}
}

func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one, %v", err)
	}
}