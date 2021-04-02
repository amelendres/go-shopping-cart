package cart

import (
	"errors"
	"fmt"
)

func ErrProductNotFound(id string) error {
	return errors.New(fmt.Sprintf("Product <%s> not found", id))
}

func ErrCartNotFound(id string) error {
	return errors.New(fmt.Sprintf("Cart <%s> not found", id))
}

func ErrCartAlreadyExists(id string) error {
	return errors.New(fmt.Sprintf("Cart <%s> already exists", id))
}