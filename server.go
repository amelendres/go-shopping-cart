package shopping

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

// CartServer is a HTTP interface for Cart
type CartServer struct {
	repository CartRepository
	http.Handler
}

const jsonContentType = "application/json"

// NewCartServer creates a CartServer with routing configured
func NewCartServer(repository CartRepository) *CartServer {
	c := new(CartServer)
	c.repository = repository

	router := mux.NewRouter()
	router.Handle("/carts/{cartID}/products", http.HandlerFunc(c.productsHandler))

	c.Handler = router

	return c
}

func (c *CartServer) productsHandler(w http.ResponseWriter, r *http.Request) {
	cartID := mux.Vars(r)["cartID"]

	switch r.Method {
	case http.MethodPost:
		var product Product
		json.NewDecoder(r.Body).Decode(&product)
		c.processAddProduct(w, cartID, product)
	case http.MethodGet:
		c.processGetProducts(w, cartID)
	}
}

func (c *CartServer) processGetProducts(w http.ResponseWriter, cartID string) {
	cart, err := c.repository.Get(cartID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(cart.GetProducts())
}

func (c *CartServer) processAddProduct(w http.ResponseWriter, cartID string, product Product) {
	cart, _ := c.repository.Get(cartID)

	//WIP: move to CreateCart endpoint
	if cart == nil {
		cart = NewCart(UUID(cartID), UUID(uuid.New().String()))
	}

	cart.AddProduct(product)
	c.repository.Save(cart)
	w.WriteHeader(http.StatusAccepted)
}
