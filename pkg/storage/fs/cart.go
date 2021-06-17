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

func NewCartStore(f *os.File) (cart.Repository, error) {

	err := initDBFile(f)

	if err != nil {
		return nil, fmt.Errorf("problem initialising cart db file, %v", err)
	}

	carts, err := newCartsFromJSON(f)

	if err != nil {
		return nil, fmt.Errorf("problem loading carts from file %s, %v", f.Name(), err)
	}

	return &CartRepository{
		database: json.NewEncoder(&tape{f}),
		carts:    carts,
	}, nil
}

func (f *CartRepository) Get(ID string) (*cart.Cart, error) {
	c, err := f.find(ID)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (f *CartRepository) Save(c cart.Cart) error {
	curr, err := f.find(c.ID().String())
	if err != nil {
		f.carts = append(f.carts, c)
	}else{
		*curr = c
	}

	return f.database.Encode(f.carts)
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

func initDBFile(f *os.File) error {
	f.Seek(0, 0)

	info, err := f.Stat()

	if err != nil {
		return fmt.Errorf("problem getting file info from file %s, %v", f.Name(), err)
	}

	if info.Size() == 0 {
		f.Write([]byte("[]"))
		f.Seek(0, 0)
	}

	return nil
}

func (f *CartRepository) find(ID string) (*cart.Cart, error) {
	for i, c := range f.carts {
		if c.ID().String() == ID {
			return &f.carts[i], nil
		}
	}
	return nil, errors.New(fmt.Sprintf("Cart <%s> not found", ID))
}

func (f *CartRepository) key(ID string) (int, error) {
	for i, c := range f.carts {
		if c.ID().String() == ID {
			return i, nil
		}
	}
	return -1, errors.New(fmt.Sprintf("Cart <%s> not found", ID))
}

func newCartsFromJSON(rdr io.Reader) ([]cart.Cart, error) {
	var carts []cart.Cart
	err := json.NewDecoder(rdr).Decode(&carts)

	if err != nil {
		err = fmt.Errorf("problem parsing carts, %v", err)
	}

	return carts, err
}
