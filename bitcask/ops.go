package bitcask

import (
	"io"
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
		ValueSize: len(value),

		// TODO: Use i64 instead of int
		ValuePos:  int(pos),
		Timestamp: int(time.Now().Unix()),
	}

	return nil
}

// func (b *Bitcask) Get()
