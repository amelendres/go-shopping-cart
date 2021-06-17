package cart

import (
	"errors"
	"fmt"
)

func ErrProductNotFound(ID UUID) error {
	return errors.New(fmt.Sprintf("Product <%s> not found", ID))
}

func ErrCartNotFound(ID string) error {
	return errors.New(fmt.Sprintf("Cart <%s> not found", ID))
}

func ErrCartAlreadyExists(ID UUID) error {
	return errors.New(fmt.Sprintf("Cart <%s> already exists", ID.String()))
}

func ErrAddingOtherProductWithSameId(ID UUID) error {
	return errors.New(fmt.Sprintf("Error adding product id <%s> with different name or price", ID))
}