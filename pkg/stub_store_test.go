package cart_test

import (
	. "github.com/amelendres/go-shopping-cart/pkg"
)

type StubCartStore struct {
	Carts []Cart
}

func NewStubCartStore(carts []Cart) Repository {
	return &StubCartStore{
		carts,
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

