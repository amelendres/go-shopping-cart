package fs

import (
	cart "github.com/amelendres/go-shopping-cart/pkg"
	"io/ioutil"
	"os"
	"testing"
)

func CreateTempFile(t *testing.T, initialData string) (*os.File, func()){
	t.Helper()

	tmpfile, err := ioutil.TempFile("", "db")

	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}

	tmpfile.Write([]byte(initialData))

	removeFile := func() {
		os.Remove(tmpfile.Name())
	}

	return tmpfile, removeFile
}

func TestFileSystemStore(t *testing.T) {

	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := CreateTempFile(t, "")
		defer cleanDatabase()

		_, err := NewCartStore(database)

		cart.AssertNoError(t, err)
	})
}
