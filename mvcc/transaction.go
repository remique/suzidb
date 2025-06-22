package mvcc

import (
	"fmt"
	"strconv"

	"example.com/suzidb/storage"
)

type Transaction struct {
	Store storage.Storage
	id    uint64
}

func BeginTransaction(store storage.Storage) (*Transaction, error) {
	// TODO: This should also lock Storage in order to ensure that it is thread-safe.

	// TODO: This key should be encoded better.
	tIdString := store.Get("mvcc:next_tid")
	if tIdString == "" {
		// If we could not get mvcc:next_tid from storage, then we should set it to 2 and current
		// transaction id to 1.
		tIdString = "1"

		err := store.Set("mvcc:next_tid", strconv.FormatUint(2, 10))
		if err != nil {
			return nil, fmt.Errorf("Could not set next_tid(%d): %s", 2, err.Error())
		}

	}

	tId, err := strconv.ParseUint(tIdString, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Could not convert next_tid(%s) to uint64: %s", tIdString, err.Error())
	}

	if tId == 0 {
		tId = 1
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

func (t *Transaction) writeWithVersion(key, value string) error {
	vkv := VersionedKeyValue{
		key:     key,
		value:   value,
		version: t.id,
	}
	k, v := vkv.encode()

	err := t.Store.Set(k, v)
	if err != nil {
		return err
	}

	return nil
}

// func (t *Transaction) getLatestValue(key string) (VersionedKeyValue, error) {
// 	//
// }
