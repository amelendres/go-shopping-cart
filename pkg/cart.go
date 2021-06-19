package cart

type UUID string

type Products []Product

type Repository interface {
	Save(c Cart) error
	Get(ID string) (*Cart, error)
}

type Cart struct {
	id, buyerID UUID
	products    Products
}

func NewCart(ID, buyerID UUID) *Cart {
	return &Cart{ID, buyerID, nil}
}

func (c *Cart) ID() UUID {
	return c.id
}

func (c *Cart) BuyerID() UUID {
	return c.id
}

func (c *Cart) AddProduct(productID UUID, name string, price Price, units Quantity) error {
	prod := NewProduct(productID, name, price, units)
	p := c.find(productID)
	if p == nil {
		c.products = append(c.products, *prod)
		return nil
	}

	if p.Equal(*prod) {
		*p = *p.IncQty(prod.Qty())
		return nil
	}

	return ErrAddingOtherProductWithSameId(productID)
}

func (c *Cart) Products() Products {
	return c.products
}

//Product
func (c *Cart) Product(productID UUID) (*Product, error) {
	p := c.find(productID)
	if p != nil {
		p1 := *p
		return &p1, nil
	}
	return nil, ErrProductNotFound(productID)
}

//func (c *Cart) IncProductQty(productID UUID, qty Quantity) error {
//
//	p := c.find(productID)
//	if p != nil {
//		*p = *p.IncQty(qty)
//		return nil
//	}
//	return ErrProductNotFound(productID)
//}

func (c *Cart) find(productID UUID) *Product {
	for i, p := range c.products {
		if p.id == productID {
			return &c.products[i]
		}
	}
	return nil
}

//Value Objects
func (u UUID) String() string {
	return string(u)
}
