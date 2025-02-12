package bitcask

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDataFile(t *testing.T) {
	tmpDir := t.TempDir()
	res, err := NewDataFile(tmpDir, 1)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 1, res.Id)

	defer os.Remove(filepath.Join(tmpDir, "1.db"))
}
