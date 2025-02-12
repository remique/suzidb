package bitcask

func (b *Bitcask) Set(key, value string) error {
	record := NewDiskRecord(key, value)

	// Marshal values
	serialized, err := record.encode()
	if err != nil {
		return err
	}

	// Add to file
	_, err = b.ActiveFile.Fd.Write(serialized)
	if err != nil {
		return err
	}

	// Sync file
	err = b.ActiveFile.Fd.Sync()
	if err != nil {
		return err
	}

	// Set Record in keydir
	// b.KeyDir[key] = KeyDirRecord{
	// 	FileId:    1,
	// 	ValueSize: 1,
	// 	ValuePos:  1,
	// 	Timestamp: 1,
	// }

	return nil
}
