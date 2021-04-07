package cart

type UUID string

type Products []Product

type Repository interface {
	Save(c Cart) error
	Get(ID string) (*Cart, error)
}

type Cart struct {
	ID, BuyerID UUID
	Products    Products
}

func NewCart(ID, buyerID UUID) *Cart {
	return &Cart{ID, buyerID, nil}
}

func (c *Cart) AddProduct(productID, name string, price float64, units int) error {
	prod := NewProduct(productID, name, Price(price), Quantity(units))
	p := c.find(productID)
	if p == nil {
		c.Products = append(c.Products, *prod)
		return nil
	}

	if p.Equal(*prod) {
		*p = *p.IncQty(prod.Units)
		return nil
	}

	return ErrAddingOtherProductWithSameId(productID)
}

func (c *Cart) GetProducts() Products {
	return c.Products
}

//GetProduct
func (c *Cart) GetProduct(productID string) (*Product, error) {
	p := c.find(productID)
	if p != nil {
		p1 := *p
		return &p1, nil
	}
	return nil, ErrProductNotFound(productID)
}

func (c *Cart) IncProductQty(productID string, qty int) error {

	p := c.find(productID)
	if p != nil {
		*p = *p.IncQty(Quantity(qty))
		return nil
	}
	return ErrProductNotFound(productID)
}

func (c *Cart) find(productID string) *Product {
	for i, p := range c.Products {
		if p.ID == productID {
			return &c.Products[i]
		}
	}
	return nil
}

//Value Objects
func (u *UUID) String() string {
	return string(*u)
}
