package fs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	cart "github.com/amelendres/go-shopping-cart/pkg"
)

// CartRepository stores carts in filesystem
type CartRepository struct {
	database *json.Encoder
	carts    []cart.Cart
}

func NewCartStore(file *os.File) (cart.Repository, error) {

	err := initialiseCartDBFile(file)

	if err != nil {
		return nil, fmt.Errorf("problem initialising cart db file, %v", err)
	}

	carts, err := newCartsFromJSON(file)

	if err != nil {
		return nil, fmt.Errorf("problem loading carts from file %s, %v", file.Name(), err)
	}

	return &CartRepository{
		database: json.NewEncoder(&tape{file}),
		carts:    carts,
	}, nil
}

// CartStoreFromFile creates a CartStore from the contents of a JSON file found at path
func CartStoreFromFile(path string) (cart.Repository, func(), error) {
	db, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return nil, nil, fmt.Errorf("problem opening %s %v", path, err)
	}

	closeFunc := func() {
		db.Close()
	}

	store, err := NewCartStore(db)

	if err != nil {
		return nil, nil, fmt.Errorf("problem creating file system cart store, %v ", err)
	}

	return store, closeFunc, nil
}

func initialiseCartDBFile(file *os.File) error {
	file.Seek(0, 0)

	info, err := file.Stat()

	if err != nil {
		return fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err)
	}

	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, 0)
	}

	return nil
}

func (f *CartRepository) Get(id string) (*cart.Cart, error) {
	c, err := f.find(id)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (f *CartRepository) Save(cart cart.Cart) error {
	c, err := f.find(string(cart.ID))
	if err != nil {
		f.carts = append(f.carts, cart)
		return nil
	}

	*c = cart
	return f.database.Encode(f.carts)
}

func (f *CartRepository) find(id string) (*cart.Cart, error) {
	for i, c := range f.carts {
		if string(c.ID) == id {
			return &f.carts[i], nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Cart <%s> not found", id))
}

func newCartsFromJSON(rdr io.Reader) ([]cart.Cart, error) {
	var carts []cart.Cart
	err := json.NewDecoder(rdr).Decode(&carts)

	if err != nil {
		err = fmt.Errorf("problem parsing carts, %v", err)
	}

	return carts, err
}
