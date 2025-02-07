package bitcask

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDataFile(t *testing.T) {
	res, err := NewDataFile(".", 1)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 1, res.Id)

	defer os.Remove("./1.db")
}
