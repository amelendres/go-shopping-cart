package rest

import (
	"encoding/json"
	cart "github.com/amelendres/go-shopping-cart/pkg"
	"github.com/google/uuid"
	"net/http"

	"github.com/amelendres/go-shopping-cart/internal/logging"
	"github.com/amelendres/go-shopping-cart/pkg/adding"
	"github.com/amelendres/go-shopping-cart/pkg/creating"
	"github.com/amelendres/go-shopping-cart/pkg/listing"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// CartServer is a HTTP interface for Cart
type CartServer struct {
	cartCreator   creating.CartCreator
	productAdder  adding.ProductAdder
	productLister listing.ProductLister
	http.Handler
}

const jsonContentType = "application/json"

// NewCartServer creates a CartServer with routing configured
func NewCartServer(
	logger *logrus.Logger,
	cc creating.CartCreator,
	pa adding.ProductAdder,
	pl listing.ProductLister,
) *CartServer {

	cs := new(CartServer)

	cs.cartCreator = cc
	cs.productAdder = pa
	cs.productLister = pl

	router := mux.NewRouter()
	router.HandleFunc("/carts", cs.createCart).Methods(http.MethodPost)
	router.HandleFunc("/carts/{cartID}/products", cs.addProduct).Methods(http.MethodPost)
	router.HandleFunc("/carts/{cartID}/products", cs.getProducts).Methods(http.MethodGet)

	loggingMiddleware := logging.Middleware(logger)
	loggedRouter := loggingMiddleware(router)

	cs.Handler = loggedRouter

	return cs
}

func (cs *CartServer) createCart(w http.ResponseWriter, r *http.Request) {
	req := creating.CreateCartReq{}
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil || req.CartID == "" || req.BuyerID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errors.Errorf("request must not be nil"))
		return
	}

	err = cs.cartCreator.Create(req)

	if err != nil {
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(err)
		return
	}
}

func (cs *CartServer) addProduct(w http.ResponseWriter, r *http.Request) {
	req := adding.AddProductReq{CartID: mux.Vars(r)["cartID"]}
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errors.Errorf("request must not be nil"))
		return
	}
	req.CartID = mux.Vars(r)["cartID"]

	err = cs.productAdder.Add(adding.AddProductReq{
		CartID:    req.CartID,
		ProductID: req.ProductID,
		Name:      req.Name,
		Price:     req.Price,
		Units:     req.Units,
	})
	if err != nil && err.Error() == cart.ErrCartNotFound(req.CartID).Error() {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(err)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (cs *CartServer) getProducts(w http.ResponseWriter, r *http.Request) {
	cartID := mux.Vars(r)["cartID"]

	_, err := uuid.Parse(cartID)
	if  err!= nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	lp, err := cs.productLister.ListProducts(cartID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("content-type", jsonContentType)
	_ = json.NewEncoder(w).Encode(lp)
}
