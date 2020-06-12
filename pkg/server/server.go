package server

import (
	"encoding/json"
	shopping "github.com/amelendres/go-shopping-cart/pkg"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

// CartServer is a HTTP interface for Cart
type CartServer struct {
	repository shopping.CartRepository
	http.Handler
}

const jsonContentType = "application/json"

// NewCartServer creates a CartServer with routing configured
func NewCartServer(repository shopping.CartRepository) *CartServer {
	cs := new(CartServer)
	cs.repository = repository

	router := mux.NewRouter()
	router.Handle("/carts/{cartID}/products", http.HandlerFunc(cs.productsHandler))

	cs.Handler = router

	return cs
}

func (cs *CartServer) productsHandler(w http.ResponseWriter, r *http.Request) {
	cartID := mux.Vars(r)["cartID"]

	switch r.Method {
	case http.MethodPost:
		var product shopping.Product
		json.NewDecoder(r.Body).Decode(&product)
		cs.processAddProduct(w, cartID, product)
	case http.MethodGet:
		cs.processGetProducts(w, cartID)
	}
}

func (cs *CartServer) processGetProducts(w http.ResponseWriter, cartID string) {
	cart, err := cs.repository.Get(cartID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(cart.GetProducts())
}

func (cs *CartServer) processAddProduct(w http.ResponseWriter, cartID string, product shopping.Product) {
	cart, _ := cs.repository.Get(cartID)

	//WIP: move to CreateCart endpoint
	if cart == nil {
		cart = shopping.NewCart(shopping.UUID(cartID), shopping.UUID(uuid.New().String()))
	}

	cart.AddProduct(product)
	cs.repository.Save(cart)
	w.WriteHeader(http.StatusAccepted)
}
