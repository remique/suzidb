package bitcask

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func (b *Bitcask) Set(key, value string) error {
	valueBytes := bytes.NewBufferString(value).Bytes()
	// Generate new header
	h := NewHeader(key, valueBytes)

	record := DiskRecord{
		Header: *h,
		Key:    key,
		Value:  valueBytes,
	}

	// Marshal values
	serialized, err := json.Marshal(record)
	if err != nil {
		return err
	}

	fmt.Println(serialized)
	fmt.Println(record)

	// Add to file
	n, err := b.ActiveFile.Fd.Write(serialized)
	fmt.Println("Wrote ", n, "bytes")

	if err != nil {
		return err
	}
	err = b.ActiveFile.Fd.Sync()
	if err != nil {
		return err
	}

	// b.KeyDir[key] = KeyDirRecord{
	// 	FileId:    1,
	// 	ValueSize: 1,
	// 	ValuePos:  1,
	// 	Timestamp: 1,
	// }

	return nil
}
