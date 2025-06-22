package mvcc

import (
	"testing"

	"example.com/suzidb/storage"

	"github.com/stretchr/testify/assert"
)

func TestNewTransaction(t *testing.T) {
	ms := storage.NewMemStorage()

	trans1, err := BeginTransaction(ms)
	assert.NoError(t, err)
	assert.Equal(t, trans1.id, uint64(1))

	trans2, err := BeginTransaction(ms)
	assert.NoError(t, err)
	assert.Equal(t, trans2.id, uint64(2))
}
