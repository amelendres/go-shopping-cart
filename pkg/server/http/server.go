package http

import (
	"encoding/json"
	cart "github.com/amelendres/go-shopping-cart/pkg"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

// CartServer is a HTTP interface for Cart
type CartServer struct {
	repository cart.Repository
	http.Handler
}

const jsonContentType = "application/json"

// NewCartServer creates a CartServer with routing configured
func NewCartServer(repository cart.Repository) *CartServer {
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
		var product cart.Product
		_ = json.NewDecoder(r.Body).Decode(&product)
		cs.processAddProduct(w, cartID, product)
	case http.MethodGet:
		cs.processGetProducts(w, cartID)
	}
}

func (cs *CartServer) processGetProducts(w http.ResponseWriter, cartID string) {
	c, err := cs.repository.Get(cartID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", jsonContentType)
	_ = json.NewEncoder(w).Encode(c.GetProducts())
}

func (cs *CartServer) processAddProduct(w http.ResponseWriter, cartID string, p cart.Product) {
	c, err := cs.repository.Get(cartID)

	//WIP: move to CreateCart endpoint
	if err != nil {
		c = cart.NewCart(cart.UUID(cartID), cart.UUID(uuid.New().String()))
	}
	c.AddProduct(p.ID, p.Name, float64(p.Price), int(p.Units))

	_ = cs.repository.Save(*c)
	w.WriteHeader(http.StatusAccepted)
}
