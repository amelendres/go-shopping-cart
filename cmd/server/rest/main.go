package main

import (
	"log"
	"net/http"
	"os"

	"github.com/amelendres/go-shopping-cart/pkg/adding"
	"github.com/amelendres/go-shopping-cart/pkg/creating"
	"github.com/amelendres/go-shopping-cart/pkg/listing"
	server "github.com/amelendres/go-shopping-cart/pkg/server/rest"
	"github.com/amelendres/go-shopping-cart/pkg/storage/pgsql"
	"github.com/sirupsen/logrus"
)

func main() {

	dbURI := os.Getenv("DATABASE_URI")
	conn, err := pgsql.NewConn(dbURI)
	if err != nil {
		log.Fatal(err)
	}
	store := pgsql.NewCartRepository(conn)

	logger := logrus.New()

	cs := server.NewCartServer(
		logger,
		creating.NewCartCreator(store),
		adding.NewProductAdder(store),
		listing.NewProductLister(store),
	)

	if err := http.ListenAndServe(":5000", cs); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
