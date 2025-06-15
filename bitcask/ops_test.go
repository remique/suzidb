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
		ValueSize: len(buf),
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

func TestSetMultiple(t *testing.T) {
	tmpDir := t.TempDir()
	b, err := NewBitcask(WithDir(tmpDir))
	assert.NoError(t, err)

	// Set keys and values
	err = b.Set("key_a", "value_a")
	assert.NoError(t, err)

	err = b.Set("key_b", "value_b")
	assert.NoError(t, err)

	get1, err := b.Get("key_a")
	assert.NoError(t, err)

	get2, err := b.Get("key_b")
	assert.NoError(t, err)

	firstExpected := DiskRecord{
		Header: Header{
			Crc:       crc32.ChecksumIEEE(bytes.NewBufferString("key_a").Bytes()),
			KeySize:   uint32(len("key_a")),
			ValueSize: uint32(len("value_a")),
		},
		Key:   "key_a",
		Value: []byte("value_a"),
	}

	secondExpected := DiskRecord{
		Header: Header{
			Crc:       crc32.ChecksumIEEE(bytes.NewBufferString("key_b").Bytes()),
			KeySize:   uint32(len("key_b")),
			ValueSize: uint32(len("value_b")),
		},
		Key:   "key_b",
		Value: []byte("value_b"),
	}

	// Is there a better way to do this?
	assert.Equal(t, firstExpected.Value, get1.Value)
	assert.Equal(t, firstExpected.Key, get1.Key)
	assert.Equal(t, firstExpected.Header.Crc, get1.Header.Crc)
	assert.Equal(t, firstExpected.Header.KeySize, get1.Header.KeySize)
	assert.Equal(t, firstExpected.Header.ValueSize, get1.Header.ValueSize)
	assert.InDelta(t, time.Now().Unix(), get1.Header.Timestamp, 2,
		"timestamp should be within 2 seconds")

	assert.Equal(t, secondExpected.Value, get2.Value)
	assert.Equal(t, secondExpected.Key, get2.Key)
	assert.Equal(t, secondExpected.Header.Crc, get2.Header.Crc)
	assert.Equal(t, secondExpected.Header.KeySize, get2.Header.KeySize)
	assert.Equal(t, secondExpected.Header.ValueSize, get2.Header.ValueSize)
	assert.InDelta(t, time.Now().Unix(), get2.Header.Timestamp, 2,
		"timestamp should be within 2 seconds")
}
