package mvcc

import (
	"fmt"
	"strconv"

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

type Transaction struct {
	Store storage.Storage
	id    uint64
}

func BeginTransaction(store storage.Storage) (*Transaction, error) {
	// TODO: This should also lock Storage in order to ensure that it is thread-safe.

	// TODO: This key should be encoded better.
	tIdString := store.Get("mvcc:next_tid")
	if tIdString == "" {
		return nil, fmt.Errorf("Could not get next_tid from storage")
	}

	tId, err := strconv.ParseUint(tIdString, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Could not convert next_tid(%s) to uint64: %s", tIdString, err.Error())
	}

	// Now increase next_tid
	nextTid := strconv.FormatUint(tId+1, 10)
	err = store.Set("mvcc:next_tid", nextTid)
	if err != nil {
		return nil, fmt.Errorf("Could not set next_tid(%s): %s", nextTid, err.Error())
	}

	return &Transaction{
		Store: store,
		id:    tId,
	}, nil
}

// Formats VersionedKeyValue into writeable key and value pair into the storage.
func (vkv *VersionedKeyValue) encode() (key string, value string) {
	formattedKey := fmt.Sprintf("%s::%s", key, vkv.version)

	return formattedKey, value
}

func (t *Transaction) Set(key, value string) error {
	// TODO: Write a versioned key value

	return nil
}
