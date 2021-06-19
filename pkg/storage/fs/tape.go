package fs

import (
	"os"
)

type tape struct {
	file *os.File
}

func NewTape(file *os.File) *tape{
	return &tape{file}
}

func (t *tape) Write(p []byte) (n int, err error) {
	t.file.Truncate(0)
	t.file.Seek(0, 0)
	return t.file.Write(p)
}
