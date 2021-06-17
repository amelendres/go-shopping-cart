package grpc

import (
	"context"
	"github.com/amelendres/go-shopping-cart/pkg/adding"
	"github.com/amelendres/go-shopping-cart/pkg/creating"
	"github.com/amelendres/go-shopping-cart/pkg/listing"
	cartgrpc "github.com/amelendres/go-shopping-cart/proto"
	"github.com/pkg/errors"
)

type cartHandler struct {
	cartCreator   creating.CartCreator
	productAdder  adding.ProductAdder
	productLister listing.ProductLister
}

// NewCartServiceServer provides Cart gRPC operations
func NewCartServiceServer(
	cc creating.CartCreator,
	pa adding.ProductAdder,
	pl listing.ProductLister,
) cartgrpc.CartServiceServer {
	return &cartHandler{cartCreator: cc, productAdder: pa, productLister: pl}
}

func (s cartHandler) Create(ctx context.Context, req *cartgrpc.CreateCartReq) (*cartgrpc.CreateCartResp, error) {
	if req == nil {
		return nil, errors.Errorf("request must not be nil")
	}

	if req.Cart == nil {
		return nil, errors.Errorf("Cart but not be empty in the request")
	}

	err := s.cartCreator.Create(creating.CreateCartReq{req.Cart.Id, req.Cart.BuyerId})
	if err != nil {
		return nil, err
	}
	return &cartgrpc.CreateCartResp{}, nil
}

func (s cartHandler) Add(ctx context.Context, req *cartgrpc.AddProductReq) (*cartgrpc.AddProductResp, error) {
	if req == nil {
		return nil, errors.Errorf("request must not be nil")
	}

	if req.Product == nil {
		return nil, errors.Errorf("Product but not be empty in the request")
	}

	err := s.productAdder.Add(
		adding.AddProductReq{
			req.CartId,
			req.Product.Id,
			req.Product.Name,
			req.Product.UnitPrice,
			int(req.Product.Units),
		},
	)
	if err != nil {
		//log.Printf("error: %+v", err)
		return nil, err
	}
	return &cartgrpc.AddProductResp{}, nil
}

func (s cartHandler) List(ctx context.Context, req *cartgrpc.ListCartReq) (*cartgrpc.ListCartResp, error) {
	if req == nil {
		return nil, errors.Errorf("request must not be nil")
	}

	products, err := s.productLister.ListProducts(req.CartId)
	if err != nil {
		return nil, err
	}
	return &cartgrpc.ListCartResp{Products: mapSliceOfProducts(products)}, nil
}

func mapSliceOfProducts(products []listing.ProductResponse) (grpcProducts []*cartgrpc.Product) {
	for _, p := range products {
		grpcProducts = append(grpcProducts, &cartgrpc.Product{
			Id:        p.ID,
			Name:      p.Name,
			UnitPrice: p.Price,
			Units:     int32(p.Units),
		})
	}
	return
}
