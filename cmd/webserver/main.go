package main

import (
	"log"
	"net/http"

	shop "github.com/amelendres/go-shopping-cart"
)

const dbFileName = "cart.db.json"

func main() {
	//store, close, err := shop.FileSystemCartStoreFromFile(dbFileName)
	repository, close, err := shop.FileSystemCartStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	server := shop.NewCartServer(repository)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
