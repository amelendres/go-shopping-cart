package shop

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// FileSystemCartStore stores carts filesystem
type FileSystemCartStore struct {
	database *json.Encoder
	products Products
}

func NewFileSystemCartStore(file *os.File) (*FileSystemCartStore, error) {

	err := initialiseCartDBFile(file)

	if err != nil {
		return nil, fmt.Errorf("problem initialising cart db file, %v", err)
	}

	products, err := NewProducts(file)

	if err != nil {
		return nil, fmt.Errorf("problem loading cart store from file %s, %v", file.Name(), err)
	}

	return &FileSystemCartStore{
		database: json.NewEncoder(&tape{file}),
		products: products,
	}, nil
}

// FileSystemCartStoreFromFile creates a CartStore from the contents of a JSON file found at path
func FileSystemCartStoreFromFile(path string) (*FileSystemCartStore, func(), error) {
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

func (f *FileSystemCartStore) GetProducts() Products {
	sort.Slice(f.products, func(i, j int) bool {
		return f.products[i].Qty > f.products[j].Qty
	})
	return f.products
}

func (f *FileSystemCartStore) AddProduct(product Product) {
	prod := f.products.Find(product.Name)

	if prod != nil {
		prod.Qty++
	} else {
		f.products = append(f.products, product)
	}

	f.database.Encode(f.products)
}
