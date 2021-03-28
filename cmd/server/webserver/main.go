package main

import (
	server "github.com/amelendres/go-shopping-cart/pkg/server/http"
	"github.com/amelendres/go-shopping-cart/pkg/storage/fs"
	"log"
	"net/http"
)

const dbFileName = "cart.db.json"

func main() {
	repository, close, err := fs.CartStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	cs := server.NewCartServer(repository)

	if err := http.ListenAndServe(":5000", cs); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
