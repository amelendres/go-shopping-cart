package cart

type Price float64
type Quantity int

// Product stores a product line information
type Product struct {
	ID, Name string
	Price    Price
	Units    Quantity
}

func NewProduct(ID, name string, price Price, units Quantity) Product{
	return Product{ID, name, price, units}
}

func (p *Product) Equal(o Product) bool {
	return p.ID == o.ID && p.Name == o.Name && p.Price == o.Price
}

func (p *Product) IncQty(qty Quantity) Product {
	return NewProduct(p.ID, p.Name, p.Price, p.Units + qty)
}