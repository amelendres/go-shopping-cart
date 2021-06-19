package adding

import (
	. "github.com/amelendres/go-shopping-cart/pkg"
)

type ProductAdder interface {
	Add(req AddProductReq) error
}

type service struct {
	repository Repository
}

type AddProductReq struct {
	CartID    string  `json:"cart_id"`
	ProductID string  `json:"product_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Units     int     `json:"units"`
}

func (s service) Add(req AddProductReq) error {
	c, err := s.repository.Get(req.CartID)
	if err != nil {
		//return err
		//fmt.Printf("fail fetching cart %+v", err)
		return ErrCartNotFound(req.CartID)
	}

	err = c.AddProduct(UUID(req.ProductID), req.Name, Price(req.Price), Quantity(req.Units))
	if err != nil {
		return err
	}
	//log.Printf("cart: %+v", c)
	return s.repository.Save(*c)
}

func NewProductAdder(repository Repository) ProductAdder {
	return service{repository: repository}
}
