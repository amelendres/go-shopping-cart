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

func (c *Cart) AddProduct(productID, name string, price float64, units int) {
	_, err := c.GetProduct(productID)
	if err != nil {
		c.Products = append(c.Products, Product{productID, name, Price(price), Quantity(units)})
		return
	}

	_ = c.IncProductQty(productID, units)
}

func (c *Cart) GetProducts() Products {
	return c.Products
}

//GetProduct
func (c *Cart) GetProduct(productID string) (*Product, error) {
	for _, p := range c.Products {
		if p.ID == productID {
			return &p, nil
		}
	}
	return nil, ErrProductNotFound(productID)
}

func (c *Cart) IncProductQty(productID string, qty int) error {
	for i, p := range c.Products {
		if p.ID == productID {
			c.Products[i] = p.IncQty(Quantity(qty))
			return nil
		}
	}
	return ErrProductNotFound(productID)
}
