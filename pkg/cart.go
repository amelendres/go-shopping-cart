package shopping

type UUID string

type Cart struct {
	ID, BuyerID UUID
	Products    Products
}

func NewCart(id, buyerId UUID) *Cart {
	return &Cart{id, buyerId, nil}
}

func (c *Cart) AddProduct(product Product) {
	key, prod := c.Products.Find(product)

	if prod != nil {
		c.Products[key] = NewProduct(prod.ID, prod.Name, prod.Price, prod.Units+product.Units)
	} else {
		c.Products = append(c.Products, product)
	}
}

func (c *Cart) GetProducts() Products {
	return c.Products
}
