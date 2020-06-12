package main

import (
	"github.com/amelendres/go-shopping-cart/fs"
	"github.com/amelendres/go-shopping-cart/server"
	"log"
	"net/http"
)

const dbFileName = "cart.db.json"

func main() {
	repository, close, err := fs.FileSystemCartStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	cs := server.NewCartServer(repository)

	if err := http.ListenAndServe(":5000", cs); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
