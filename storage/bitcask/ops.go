package bitcask

// func (b *Bitcask) set(key, value string) error {
// 	// Add to activeFile
// 	fmt.Println("Adding to file")
// 	serialized, err := json.Marshal(KeyVal{key: value})
// 	if err != nil {
// 		return err
// 	}

// 	b.activeFile.fd.Write(serialized)

// 	// myHeader := &Header{
// 	// 	Crc:       1,
// 	// 	Timestamp: 1,
// 	// 	KeySize:   1,
// 	// 	ValueSize: 1,
// 	// }
// 	// serialized2, err := json.Marshal(myHeader)
// 	// if err != nil {
// 	// 	return err
// 	// }

// 	// b.activeFile.fd.Write(serialized2)

// 	// Add to KeyDir
// 	fmt.Println("Adding to keydir")
// 	b.KeyDir[key] = KeyDirRecord{
// 		FileId:    1,
// 		ValueSize: 1,
// 		ValuePos:  1,
// 		Timestamp: 1,
// 	}

// 	return nil
// }
