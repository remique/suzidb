package bitcask

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	b, err := NewBitcask()
	assert.NoError(t, err)

	// Set key and value
	b.Set("a", "b")

	_ = b.ActiveFile.Fd.Sync()

	// Assert that the value is in the active file.
	stat, err := b.ActiveFile.Fd.Stat()
	assert.NoError(t, err)

	// Build buffer with size of the file
	buf := make([]byte, stat.Size())
	fmt.Println(stat.Size())

	// Read
	n, err := b.ActiveFile.Fd.ReadAt(buf, 0)
	// assert.NoError(t, err)
	assert.Greater(t, n, 0)

	// Serialize the data
	var rec DiskRecord
	err = json.Unmarshal(buf, &rec)
	assert.NoError(t, err)

	expectedCrc := crc32.ChecksumIEEE(bytes.NewBufferString("a").Bytes())

	expected := DiskRecord{
		Header: Header{
			Crc: expectedCrc,
		},
		Key:   "a",
		Value: bytes.NewBufferString("b").Bytes(),
	}

	assert.Equal(t, expected.Key, rec.Key)
	assert.Equal(t, expected.Value, rec.Value)
	assert.Equal(t, expected.Header.Crc, rec.Header.Crc)
	// assert.Equal(t, expected.Header.KeySize, rec.Header.KeySize)
	// assert.Equal(t, expected.Header.ValueSize, rec.Header.ValueSize)
}
