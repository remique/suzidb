package bitcask

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigDirWorks(t *testing.T) {
	err := os.MkdirAll("tmp", 0755)
	assert.NoError(t, err)

	b, err := NewBitcask(WithDir("tmp"))
	assert.NoError(t, err)

	assert.Equal(t, b.Options.dir, "tmp")
}
