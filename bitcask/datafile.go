package bitcask

import (
	"os"
	"strconv"
)

type DataFile struct {
	Id int
	Fd *os.File
}

// Creates a new DataFile. Please note that it simply opens a file
// and keeps a reference to the file descriptor. This can be used
// for both activeFile as well as staleFiles.
func NewDataFile(dirName string, id int) (*DataFile, error) {
	idStr := strconv.Itoa(id)
	path := dirName + "/" + idStr + ".db"
	fd, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE, 0644)

	if err != nil {
		return nil, err
	}

	return &DataFile{
		Id: id,
		Fd: fd,
	}, nil
}
