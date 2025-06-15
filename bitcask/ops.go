package bitcask

import (
	"fmt"
	"io"
	"slices"
	"time"
)

func (b *Bitcask) Set(key, value string) error {
	record := NewDiskRecord(key, value)

	// Marshal values
	serialized, err := record.encode()
	if err != nil {
		return err
	}

	// Here seek the latest position at which we write it at.
	// And save it to set it in KeyDir later.
	pos, err := b.ActiveFile.Fd.Seek(0, io.SeekCurrent)
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
	b.KeyDir[key] = KeyDirRecord{
		FileId:    b.ActiveFile.Id,
		ValueSize: len(serialized),

		// TODO: Use i64 instead of int
		ValuePos:  int(pos),
		Timestamp: int(time.Now().Unix()),
	}

	return nil
}

func (b *Bitcask) Get(key string) (*DiskRecord, error) {
	fromKeydir, ok := b.KeyDir[key]
	if !ok {
		return nil, fmt.Errorf("No value found")
	}

	fileToRead := b.ActiveFile
	buffer := make([]byte, fromKeydir.ValueSize)

	if fromKeydir.FileId != b.ActiveFile.Id {
		// We need to find b.StaleFiles ID
		idx := slices.IndexFunc(b.StaleFiles, func(df *DataFile) bool { return df.Id == fromKeydir.FileId })
		if idx < 0 {
			return nil, fmt.Errorf("No fileId: %d", fromKeydir.FileId)
		}

		fileToRead = b.StaleFiles[idx]
	}

	_, err := fileToRead.Fd.ReadAt(buffer, int64(fromKeydir.ValuePos))
	if err != nil {
		return nil, err
	}

	fmt.Println("buffer", string(buffer))

	// Decode bytes -> DiskRecord
	dr, err := decode(buffer)
	if err != nil {
		return nil, err
	}

	return dr, nil
}
