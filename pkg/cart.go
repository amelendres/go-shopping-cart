package cart

type UUID string

type Products []Product

type Repository interface {
	Save(cart Cart) error
	Get(id string) (*Cart, error)
}

type Cart struct {
	ID, BuyerID UUID
	Products    Products
}

func NewCart(id, buyerId UUID) *Cart {
	return &Cart{id, buyerId, nil}
}

func (c *Cart) AddProduct(productId, name string, price float64, units int) {
	_, err := c.GetProduct(productId)
	if err != nil {
		c.Products = append(c.Products, Product{productId, name, Price(price), Quantity(units)})
		return
	}

	_ = c.IncProductQty(productId, units)
}

func (c *Cart) GetProducts() Products {
	return c.Products
}

//GetProduct
func (c *Cart) GetProduct(productId string) (*Product, error) {
	for _, p := range c.Products {
		if p.ID == productId {
			return &p, nil
		}
	}
	return nil, ErrProductNotFound(productId)
}

func (c *Cart) IncProductQty(productId string, qty int) error {
	for i, p := range c.Products {
		if p.ID == productId {
			c.Products[i] = p.IncQty(Quantity(qty))
			return nil
		}
	}
	return ErrProductNotFound(productId)
}
