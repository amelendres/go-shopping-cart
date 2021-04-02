package cart

import (
	"testing"
)

type StubCartStore struct {
	Carts []Cart
}

func NewStubCartStore() Repository {
	return &StubCartStore{
		[]Cart{},
	}
}


func (s *StubCartStore) Save(c Cart) error {
	s.Carts = append(s.Carts, c)
	return nil
}

func (s *StubCartStore) Get(ID string) (*Cart, error) {
	if len(s.Carts) < 1 {
		return nil, ErrCartNotFound(ID)
	}
	return &s.Carts[0], nil
}

func AssertCartProduct(t *testing.T, c *Cart, p Product) {
	t.Helper()

	if len(c.Products) != 1 {
		t.Fatalf("got %d calls to AddProduct want %d", len(c.Products), 1)
	}

	if c.Products[0] != p {
		t.Errorf("did not cart the correct product got %+v want %+v", c.Products[0], p)
	}
}

func AssertCartProductQty(t *testing.T, c *Cart, productID string, qty Quantity) {
	t.Helper()

	p, err := c.GetProduct(productID)
	if err != nil {
		t.Fatalf("Product <%s> Not Found", productID)
	}

	if p.Units != qty {
		t.Errorf("did not cart the correct product Units, got %d want %d", p.Units, qty)
	}
}

func AssertCartProductLines(t *testing.T, c *Cart, qty int) {
	t.Helper()
	lines := len(c.Products)
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