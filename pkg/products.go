package shopping

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

