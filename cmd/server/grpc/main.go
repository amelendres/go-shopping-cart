package main

import (
	cart "github.com/amelendres/go-shopping-cart/pkg"
	"github.com/amelendres/go-shopping-cart/pkg/adding"
	"github.com/amelendres/go-shopping-cart/pkg/creating"
	"github.com/amelendres/go-shopping-cart/pkg/listing"
	"github.com/amelendres/go-shopping-cart/pkg/server"
	"github.com/amelendres/go-shopping-cart/pkg/server/grpc"
	"github.com/amelendres/go-shopping-cart/pkg/storage/fs"
	"log"
	"os"
)

const (
	ServerProtocol = "tcp"
	ServerHost     = "localhost"
	ServerPort     = "3333"
	dbFileName     = "cart.db.json"
)

func main() {

	var (
		protocol = getEnv("SERVER_PROTOCOL", ServerProtocol)
		host     = getEnv("SERVER_HOST", ServerHost)
		port     = getEnv("SERVER_PORT", ServerPort)
	)
	srvCfg := server.Config{Protocol: protocol, Host: host, Port: port}
	store, closeStore := getStore()
	srv := grpc.NewCartServer(
		srvCfg,
		creating.NewCartCreator(store),
		adding.NewProductAdder(store),
		listing.NewProductLister(store))

	log.Printf("gRPC server running at %s://%s:%s ...\n", protocol, host, port)
	log.Fatal(srv.Serve())
	closeStore()
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

func getStore() (cart.Repository, func()) {
	store, closeStore, err := fs.CartStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	return store, closeStore
}
