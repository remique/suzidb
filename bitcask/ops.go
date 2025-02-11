package bitcask

import (
	"bytes"
	"encoding/json"
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

	// Add to file
	_, err = b.ActiveFile.Fd.Write(serialized)

	if err != nil {
		return err
	}
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
