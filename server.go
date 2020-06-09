package shop

import (
	"encoding/json"
	"net/http"
)

// CartServer is a HTTP interface for Cart
type CartServer struct {
	store CartStore
	http.Handler
}

const jsonContentType = "application/json"

// NewCartServer creates a CartServer with routing configured
func NewCartServer(store CartStore) *CartServer {
	c := new(CartServer)

	c.store = store

	router := http.NewServeMux()
	router.Handle("/products", http.HandlerFunc(c.productsHandler))

	c.Handler = router

	return c
}

func (c *CartServer) productsHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		var product Product
		json.NewDecoder(r.Body).Decode(&product)
		c.processAddProduct(w, product)
	case http.MethodGet:
		c.processGetProducts(w)
	}
}

func (c *CartServer) processGetProducts(w http.ResponseWriter) {
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(c.store.GetProducts())
}

func (c *CartServer) processAddProduct(w http.ResponseWriter, product Product) {
	c.store.AddProduct(product)
	w.WriteHeader(http.StatusAccepted)
}
