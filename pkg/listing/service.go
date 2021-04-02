package listing

import cart "github.com/amelendres/go-shopping-cart/pkg"

type ProductLister interface {
	ListProducts(string) ([]cart.Product, error)
}

type service struct {
	repository cart.Repository
}

func (s service) ListProducts(cartID string) ([]cart.Product, error) {
	c, err := s.repository.Get(cartID)
	if err != nil {
		return nil, err
	}

	return c.Products, nil
}

func NewProductLister(repository cart.Repository) ProductLister {
	return service{repository: repository}
}
