package shopping

import "errors"

var ErrUnknownCart = errors.New("unknown cart")

type CartRepository interface {
	Save(cart *Cart) error
	Get(id string) (*Cart, error)
}
