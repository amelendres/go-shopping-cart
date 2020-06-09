package shop

import (
	"encoding/json"
	"fmt"
	"io"
)

// Products stores a collection of products
type Products []Product

//func (p Products) Add(product Product){
//	current := p.Find(product.Name)
//	if current != nil {
//		current.Qty++
//	} else {
//		p = append(p, product)
//	}
//}

func (p Products) Find(name string) *Product {
	for i, pro := range p {
		if pro.Name == name {
			return &p[i]
		}
	}
	return nil
}

// NewProducts creates a products from JSON
func NewProducts(rdr io.Reader) ([]Product, error) {
	var products []Product
	err := json.NewDecoder(rdr).Decode(&products)

	if err != nil {
		err = fmt.Errorf("problem parsing products, %v", err)
	}

	return products, err
}
