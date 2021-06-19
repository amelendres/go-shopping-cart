package main

import (
	"github.com/amelendres/go-shopping-cart/pkg/adding"
	"github.com/amelendres/go-shopping-cart/pkg/creating"
	"github.com/amelendres/go-shopping-cart/pkg/listing"
	"github.com/amelendres/go-shopping-cart/pkg/server"
	"github.com/amelendres/go-shopping-cart/pkg/server/grpc"
	"github.com/amelendres/go-shopping-cart/pkg/storage/pgsql"
	"log"
	"os"
)

const (
	ServerProtocol = "tcp"
	ServerHost     = "localhost"
	ServerPort     = "3333"
)

func main() {

	var (
		protocol = getEnv("SERVER_PROTOCOL", ServerProtocol)
		host     = getEnv("SERVER_HOST", ServerHost)
		port     = getEnv("SERVER_PORT", ServerPort)
	)
	srvCfg := server.Config{Protocol: protocol, Host: host, Port: port}
	dbURI := os.Getenv("DATABASE_URI")
	conn, err := pgsql.NewConn(dbURI)
	if err != nil {
		log.Fatal(err)
	}
	store := pgsql.NewCartRepository(conn)

	srv := grpc.NewCartServer(
		srvCfg,
		creating.NewCartCreator(store),
		adding.NewProductAdder(store),
		listing.NewProductLister(store))

	log.Printf("gRPC server running at %s://%s:%s ...\n", protocol, host, port)
	log.Fatal(srv.Serve())
	conn.Close()
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
