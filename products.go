package shopping

import (
	"encoding/json"
	"fmt"
	"io"
)

// Products stores a collection of Products
type Products []Product

func (p Products) Find(product Product) (int, *Product) {
	for i, pro := range p {
		if pro.Equal(product) {
			return i, &p[i]
		}
	}
	return 0, nil
}

// NewProducts creates a Products from JSON
func NewProducts(rdr io.Reader) ([]Product, error) {
	var products []Product
	err := json.NewDecoder(rdr).Decode(&products)

	if err != nil {
		err = fmt.Errorf("problem parsing Products, %v", err)
	}

	return products, err
}
