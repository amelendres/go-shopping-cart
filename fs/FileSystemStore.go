package fs

import (
	"encoding/json"
	"fmt"
	"github.com/amelendres/go-shopping-cart"
	"io"
	"os"
)

// FileSystemCartRepository stores carts filesystem
type FileSystemCartRepository struct {
	database *json.Encoder
	carts    []shopping.Cart
}

func NewFileSystemCartStore(file *os.File) (*FileSystemCartRepository, error) {

	err := initialiseCartDBFile(file)

	if err != nil {
		return nil, fmt.Errorf("problem initialising cart db file, %v", err)
	}

	carts, err := newCartsFromJSON(file)

	if err != nil {
		return nil, fmt.Errorf("problem loading carts from file %s, %v", file.Name(), err)
	}

	return &FileSystemCartRepository{
		database: json.NewEncoder(&tape{file}),
		carts:    carts,
	}, nil
}

// FileSystemCartStoreFromFile creates a CartStore from the contents of a JSON file found at path
func FileSystemCartStoreFromFile(path string) (*FileSystemCartRepository, func(), error) {
	db, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return nil, nil, fmt.Errorf("problem opening %s %v", path, err)
	}

	closeFunc := func() {
		db.Close()
	}

	store, err := NewFileSystemCartStore(db)

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

func (f *FileSystemCartRepository) Get(id string) (*shopping.Cart, error) {
	_, cart := f.find(id)
	if cart == nil {
		return nil, shopping.ErrUnknownCart
	}

	return cart, nil
}

func (f *FileSystemCartRepository) Save(cart *shopping.Cart) error {
	_, crt := f.find(string(cart.ID))

	if crt != nil {
		crt = cart
	} else {
		f.carts = append(f.carts, *cart)
	}

	err := f.database.Encode(f.carts)
	return err
}

func (f *FileSystemCartRepository) find(id string) (int, *shopping.Cart) {
	for i, cart := range f.carts {
		if string(cart.ID) == id {
			return i, &f.carts[i]
		}
	}
	return 0, nil
}

func newCartsFromJSON(rdr io.Reader) ([]shopping.Cart, error) {
	var carts []shopping.Cart
	err := json.NewDecoder(rdr).Decode(&carts)

	if err != nil {
		err = fmt.Errorf("problem parsing carts, %v", err)
	}

	return carts, err
}
