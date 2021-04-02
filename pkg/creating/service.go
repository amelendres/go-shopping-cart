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
		return cart.ErrCartAlreadyExists(string(c.ID))
	}

	w := cart.NewCart(cart.UUID(cartID), cart.UUID(buyerID))
	err := s.repository.Save(*w)
	if err != nil {
		return err
	}
	return nil
}

func NewCartCreator(repository cart.Repository) CartCreator {
	return &service{repository: repository}
}
