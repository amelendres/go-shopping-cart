package cart

type Price float64
type Quantity int

// Product line information
type Product struct {
	id    UUID
	name  string
	price Price
	qty Quantity
}

func NewProduct(ID UUID, name string, price Price, units Quantity) *Product {
	return &Product{ID, name, price, units}
}

func (p *Product) ID() UUID {
	return p.id
}

func (p *Product) Name() string {
	return p.name
}

func (p *Product) Price() Price {
	return p.price
}

func (p *Product) Qty() Quantity {
	return p.qty
}

func (p *Product) Equal(o Product) bool {
	return Product{p.id, p.name, p.price, 0} == Product{o.id, o.name, o.price, 0}
}

func (p *Product) IncQty(qty Quantity) *Product {
	return NewProduct(p.id, p.name, p.price, p.qty+qty)
}



//Value Objects
func (q Quantity) Int() int {
	return int(q)
}
