package creating

import cart "github.com/amelendres/go-shopping-cart/pkg"

type CartCreator interface {
	Create(cartID, buyerID string) error
}

type service struct {
	repository cart.Repository
}

func (s service) Create(cartID, buyerID string) error {
	c, _ := s.repository.Get(cartID)
	if c != nil {
		return cart.ErrCartAlreadyExists(c.ID.String())
	}

	w := cart.NewCart(cart.UUID(cartID), cart.UUID(buyerID))
	return s.repository.Save(*w)
}

func NewCartCreator(repository cart.Repository) CartCreator {
	return &service{repository: repository}
}
