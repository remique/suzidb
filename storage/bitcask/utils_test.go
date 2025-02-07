package bitcask

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLastFileId(t *testing.T) {
	res, err := getLastFileId([]string{"1.db", "3.db", "2.db"})
	assert.NoError(t, err)
	assert.Equal(t, 3, res)
}

func TestGetLastFileIdMissingOne(t *testing.T) {
	res, err := getLastFileId([]string{"1.db", "3.db"})
	assert.NoError(t, err)
	assert.Equal(t, 3, res)
}

func TestGetLastFileIdInvalidFileName(t *testing.T) {
	res, err := getLastFileId([]string{"x.db", "3.db"})
	assert.Error(t, err)
	assert.Equal(t, -1, res)
}
