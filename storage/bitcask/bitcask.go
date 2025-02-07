package bitcask

type Bitcask struct {
	KeyDir KeyDir
	// TODO: Set it to DataFile
	activeFile *DataFile
	// TODO: Set it to []DataFile
	staleFiles []*DataFile
}

func NewBitcask() (*Bitcask, error) {
	af, err := NewDataFile(1)
	if err != nil {
		return nil, err
	}

	b := &Bitcask{
		KeyDir:     KeyDir{},
		activeFile: af,
	}

	// err = b.buildKeydir()
	// if err != nil {
	// 	return nil, err
	// }

	return b, nil
}

// func (b *Bitcask) buildKeydir() error {
// 	// Get the size of the file
// 	stat, err := b.activeFile.fd.Stat()
// 	if err != nil {
// 		return err
// 	}

// 	// Build buffer with size of the file
// 	buf := make([]byte, stat.Size())

// 	// Read
// 	_, err = b.activeFile.fd.Read(buf)
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

// 		fmt.Println(data)
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
