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
func (vkv *VersionedKeyValue) encode() (outputKey string, outputValue string) {
	formattedKey := fmt.Sprintf("%s::%020d", vkv.key, vkv.version)

	return formattedKey, vkv.value
}

// Transaction can be in either of the three states denoted below.
type TransactionState uint8

const (
	InProgress TransactionState = iota
	Commited
	Aborted
)

// Denotes an isolation level for a transaction. Usually managed by Transaction Manager, but
// can be constructed manually via BeginTransaction.
type IsolationLevel uint8

const (
	// Most basic isolation level and not used by default. Reads get the most recent value,
	// regardless of the transaction state.
	ReadUncommited IsolationLevel = iota

	// This is an upgraded isolation level and is the default in SuziDB.
	// Every read gets the latest _commited_ value.
	ReadCommited
)
