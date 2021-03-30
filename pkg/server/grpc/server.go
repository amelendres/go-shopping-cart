package grpc

import (
	"fmt"
	"net"
	"os"

	"github.com/amelendres/go-shopping-cart/pkg/adding"
	"github.com/amelendres/go-shopping-cart/pkg/creating"
	"github.com/amelendres/go-shopping-cart/pkg/listing"
	"github.com/amelendres/go-shopping-cart/pkg/server"
	"github.com/amelendres/go-shopping-cart/pkg/server/grpc/interceptor"
	cartgrpc "github.com/amelendres/go-shopping-cart/proto"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

type cartServer struct {
	config        server.Config
	cartCreator   creating.CartCreator
	productAdder  adding.ProductAdder
	productLister listing.ProductLister
}

func NewCartServer(
	config server.Config,
	cc creating.CartCreator,
	pa adding.ProductAdder,
	pl listing.ProductLister,
) server.Server {
	return &cartServer{config: config, cartCreator: cc, productAdder: pa, productLister: pl}
}

func (s *cartServer) Serve() error {
	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
	listener, err := net.Listen(s.config.Protocol, addr)
	if err != nil {
		return err
	}

	grpcLog := grpclog.NewLoggerV2(os.Stdout, os.Stderr, os.Stderr)
	grpclog.SetLoggerV2(grpcLog)

	srv := grpc.NewServer(withUnaryInterceptor())
	serviceServer := NewCartServiceServer(
		s.cartCreator,
		s.productAdder,
		s.productLister,
	)
	cartgrpc.RegisterCartServiceServer(srv, serviceServer)

	if err := srv.Serve(listener); err != nil {
		return err
	}

	return nil
}

func withUnaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		interceptor.LoggingServerInterceptor,
		interceptor.AuthorizationServerInterceptor,
	))
}
