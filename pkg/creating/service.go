package creating

import (
	. "github.com/amelendres/go-shopping-cart/pkg"
)

type CartCreator interface {
	Create(req CreateCartReq) error
}

type service struct {
	repository Repository
}

type CreateCartReq struct {
	CartID  string `json:"cart_id"`
	BuyerID string `json:"buyer_id"`
}

func (s service) Create(req CreateCartReq) error {
	c, _ := s.repository.Get(req.CartID)
	if c != nil {
		return ErrCartAlreadyExists(c.ID())
	}

	w := NewCart(UUID(req.CartID), UUID(req.BuyerID))
	return s.repository.Save(*w)
}

func NewCartCreator(repository Repository) CartCreator {
	return &service{repository: repository}
}
