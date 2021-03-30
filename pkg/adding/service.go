package adding

import cart "github.com/amelendres/go-shopping-cart/pkg"

type ProductAdder interface {
	AddCartProduct(cartID, productID, name string, price float64, units int) error
}

type service struct {
	repository cart.Repository
}

func (s service) AddCartProduct(cartID, productID, name string, price float64, units int) error {
	c, err := s.repository.Get(cartID)
	if err != nil {
		return err
	}

	c.AddProduct(productID, name, price, units)

	return s.repository.Save(*c)
}

func NewProductAdder(repository cart.Repository) ProductAdder {
	return service{repository: repository}
}
