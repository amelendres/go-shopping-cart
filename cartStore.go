package shop

// CartStore stores information about cart
type CartStore interface {
	AddProduct(product Product)
	GetProducts() Products
}