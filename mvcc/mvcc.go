package mvcc

import (
	"fmt"

	"example.com/suzidb/storage"
)

type MvccManager struct {
	Store storage.Storage
}

type VersionedKeyValue struct {
	key     string
	value   string
	version uint64
}

// Formats VersionedKeyValue into writeable key and value pair into the storage.
func (vkv *VersionedKeyValue) encode() (key string, value string) {
	formattedKey := fmt.Sprintf("%s::%020d", key, vkv.version)

	return formattedKey, value
}

func (t *Transaction) Set(key, value string) error {
	// TODO: Write a versioned key value

	return nil
}
