package listing

import (
	. "github.com/amelendres/go-shopping-cart/pkg"
)

type ProductLister interface {
	ListProducts(string) ([]ProductResponse, error)
}

type service struct {
	repository Repository
}

type ProductResponse struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Units int     `json:"units"`
}

func (s service) ListProducts(cartID string) ([]ProductResponse, error) {
	c, err := s.repository.Get(cartID)
	if err != nil {
		return nil, err
		//return nil, adding.ErrCartNotFound(cartID)
	}

	return mapToResponse(c.Products()), nil
}

func NewProductLister(repository Repository) ProductLister {
	return service{repository: repository}
}

func mapToResponse(prods []Product) []ProductResponse {
	var resp = []ProductResponse{}
	for _, p := range prods {
		resp = append(resp, ProductResponse{p.ID().String(), p.Name(), float64(p.Price()), p.Qty().Int()})
	}
	return resp
}
