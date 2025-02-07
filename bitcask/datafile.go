package bitcask

import (
	"os"
)

type DataFile struct {
	id int
	fd *os.File
}

func NewDataFile(id int) (*DataFile, error) {
	// TODO: Use id + ".db" instead of 1db
	fd, err := os.Open("1.db")
	if err != nil {
		return nil, err
	}

	return &DataFile{
		id: id,
		fd: fd,
	}, nil
}
