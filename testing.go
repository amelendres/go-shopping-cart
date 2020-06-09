package shop

import "testing"

type StubCartStore struct {
	Products []Product
}

func (s *StubCartStore) AddProduct(product Product) {
	s.Products = append(s.Products, product)
	//s.Products.Add(Product{name, 1})
}

func (s *StubCartStore) GetProducts() Products {
	return s.Products
}

func AssertProduct(t *testing.T, store *StubCartStore, product Product) {
	t.Helper()

	if len(store.Products) != 1 {
		t.Fatalf("got %d calls to AddProduct want %d", len(store.Products), 1)
	}

	if store.Products[0] != product {
		t.Errorf("did not store the correct product got %q want %q", store.Products[0], product)
	}
}
