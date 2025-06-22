package mvcc

import (
	"strconv"
	"testing"

	"example.com/suzidb/storage"

	"github.com/stretchr/testify/assert"
)

func TestEncodingVkv(t *testing.T) {
	vkv := VersionedKeyValue{
		key:     "somekey",
		value:   "someval",
		version: 1,
	}

	resKey, resValue := vkv.encode()

	assert.Equal(t, "somekey::00000000000000000001", resKey)
	assert.Equal(t, "someval", resValue)
}

func TestNewTransaction(t *testing.T) {
	ms := storage.NewMemStorage()

	trans1, err := BeginTransaction(ms)
	assert.NoError(t, err)
	assert.Equal(t, trans1.id, uint64(1))

	trans2, err := BeginTransaction(ms)
	assert.NoError(t, err)
	assert.Equal(t, trans2.id, uint64(2))
}

func TestTransactionLimit(t *testing.T) {
	ms := storage.NewMemStorage()

	maxuint64 := uint64(1<<64 - 1)
	maxMinusOne := strconv.FormatUint(maxuint64, 10)
	err := ms.Set("mvcc:next_tid", maxMinusOne)
	assert.NoError(t, err)

	trans1, err := BeginTransaction(ms)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1<<64-1), trans1.id)

	// After reaching the limit we receive id: 1. Which means, that
	// currently we support up to 18446744073709551615 transactions.
	trans2, err := BeginTransaction(ms)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), trans2.id)
}

func TestWriteWithVersion(t *testing.T) {
	ms := storage.NewMemStorage()

	trans1, err := BeginTransaction(ms)
	assert.NoError(t, err)
	assert.Equal(t, trans1.id, uint64(1))

	trans2, err := BeginTransaction(ms)
	assert.NoError(t, err)
	assert.Equal(t, trans1.id, uint64(1))

	err = trans1.writeWithVersion("key1", "value1")
	assert.NoError(t, err)

	err = trans2.writeWithVersion("key1", "value2")
	assert.NoError(t, err)

	res := ms.ScanWithPrefix("key1")

	expectedRed := map[string]string{
		"key1::00000000000000000001": "value1",
		"key1::00000000000000000002": "value2",
	}

	assert.Equal(t, expectedRed, res)
}
