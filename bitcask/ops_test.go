package bitcask

import (
	"bytes"
	"encoding/json"
	"hash/crc32"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetSingle(t *testing.T) {
	tmpDir := t.TempDir()
	b, err := NewBitcask(WithDir(tmpDir))
	assert.NoError(t, err)

	// Set key and value
	err = b.Set("a", "b")
	assert.NoError(t, err)

	// Assert that the value is in the active file.
	stat, err := b.ActiveFile.Fd.Stat()
	assert.NoError(t, err)

	// Build buffer with size of the file
	buf := make([]byte, stat.Size())

	// Read into buffer
	n, err := b.ActiveFile.Fd.ReadAt(buf, 0)
	assert.NoError(t, err)
	assert.Greater(t, n, 0)

	// Serialize the data
	var rec DiskRecord
	err = json.Unmarshal(buf, &rec)
	assert.NoError(t, err)

	expectedCrc := crc32.ChecksumIEEE(bytes.NewBufferString("a").Bytes())

	expected := DiskRecord{
		Header: Header{
			Crc:       expectedCrc,
			KeySize:   1,
			ValueSize: 1,
		},
		Key:   "a",
		Value: []byte("b"),
	}

	expectedKeydir := KeyDirRecord{
		FileId:    b.ActiveFile.Id,
		ValueSize: len("a"),
		ValuePos:  0,
		Timestamp: int(time.Now().Unix()),
	}

	// Test KeyDir as well

	assert.Equal(t, expectedKeydir.FileId, b.KeyDir["a"].FileId)
	assert.Equal(t, expectedKeydir.ValueSize, b.KeyDir["a"].ValueSize)
	assert.Equal(t, expectedKeydir.ValuePos, b.KeyDir["a"].ValuePos)
	assert.InDelta(t, expectedKeydir.Timestamp, b.KeyDir["a"].Timestamp, 2,
		"timestamp should be within 2 seconds")

	assert.Equal(t, expected.Key, rec.Key)
	assert.Equal(t, expected.Value, rec.Value)
	assert.Equal(t, string(expected.Value), "b")
	assert.Equal(t, expected.Header.Crc, rec.Header.Crc)
	assert.Equal(t, expected.Header.KeySize, rec.Header.KeySize)
	assert.Equal(t, expected.Header.ValueSize, rec.Header.ValueSize)
}
