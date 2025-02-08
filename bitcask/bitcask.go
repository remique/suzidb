package bitcask

import (
	// "bytes"
	"fmt"
	// "encoding/json"
	"strconv"
	"strings"
)

type Bitcask struct {
	KeyDir     KeyDir
	ActiveFile *DataFile
	staleFiles []*DataFile
}

func NewBitcask() (*Bitcask, error) {
	newActiveId, err := generateNewActiveFileId(".")
	if err != nil {
		return nil, err
	}

	rest, err := glob(".")
	if err != nil {
		return nil, err
	}

	af, err := NewDataFile(".", newActiveId)
	fmt.Println("newActive", newActiveId)
	if err != nil {
		return nil, err
	}

	b := &Bitcask{
		KeyDir:     KeyDir{},
		ActiveFile: af,
	}

	// Load stalefiles
	// Move to separate function
	for _, item := range rest {
		asInt, err := strconv.Atoi(strings.Trim(item, ".db"))
		if err != nil {
			return nil, err
		}

		sf, err := NewDataFile(".", asInt)
		b.staleFiles = append(b.staleFiles, sf)
	}

	// err = b.buildKeydir()
	// if err != nil {
	// 	return nil, err
	// }

	return b, nil
}

// func (b *Bitcask) buildKeydir() error {
// 	// Get the size of the file
// 	stat, err := b.activeFile.Fd.Stat()
// 	if err != nil {
// 		return err
// 	}

// 	// Build buffer with size of the file
// 	buf := make([]byte, stat.Size())

// 	// Read
// 	_, err = b.activeFile.Fd.Read(buf)
// 	if err != nil {
// 		return err
// 	}

// 	dec := json.NewDecoder(bytes.NewReader(buf))

// 	for dec.More() {
// 		var data map[string]string
// 		err := dec.Decode(&data)
// 		if err != nil {
// 			return err
// 		}

// 		// Assign to KeyDirRecord here
// 		// for key, value := range data {
// 		// 	keydir[key] = value
// 		// }
// 	}

// 	// FileId    int
// 	// ValueSize int
// 	// ValuePos  int
// 	// Timestamp int

// 	// Poza tym chcielibysmy zapisywac tez header i odczytywac

// 	// fmt.Println(b.KeyDir)

// 	return nil
// }
