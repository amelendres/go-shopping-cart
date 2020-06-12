package shopping

type Price float32
type Quantity int

// Product stores a product line information
type Product struct {
	ID, Name string
	Price    Price
	Units    Quantity
}
func NewProduct(ID, Name string, Price Price, Units Quantity) Product{
	return Product{ID, Name, Price, Units}
}

func (p Product) Equal(other Product) bool {
	return p.ID == other.ID && p.Name == other.Name && p.Price == other.Price
}
