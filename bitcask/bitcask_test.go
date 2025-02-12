package bitcask

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigDirWorks(t *testing.T) {
	tmpDir := t.TempDir()
	err := os.MkdirAll(tmpDir, 0755)
	assert.NoError(t, err)

	b, err := NewBitcask(WithDir(tmpDir))
	assert.NoError(t, err)

	assert.Equal(t, b.Options.dir, tmpDir)

}

func TestKeydir(t *testing.T) {
	tmpDir := t.TempDir()

	df1, err := NewDataFile(tmpDir, 1)
	assert.NoError(t, err)
	df2, err := NewDataFile(tmpDir, 2)
	assert.NoError(t, err)

	firstWrite, err := df1.Fd.WriteString(`{"Header":{"Crc":3904355907,"Timestamp":1739364247,"KeySize":1,"ValueSize":1},"Key":"b","Value":"Yg=="}{"Header":{"Crc":3904355907,"Timestamp":1739364247,"KeySize":1,"ValueSize":1},"Key":"c","Value":"Yg=="}`)
	assert.NoError(t, err)
	assert.Greater(t, firstWrite, 0)
	secondWrite, err := df2.Fd.WriteString(`{"Header":{"Crc":3904355907,"Timestamp":1739364247,"KeySize":1,"ValueSize":1},"Key":"d","Value":"Yg=="}`)
	assert.NoError(t, err)
	assert.Greater(t, secondWrite, 0)

	b, err := NewBitcask(WithDir(tmpDir))
	assert.NoError(t, err)

	assert.Equal(t, len(b.KeyDir), 3)

	t.Fatal(b.KeyDir)

}
