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

	recs := []struct {
		diskRecord DiskRecord
		dataFile   *DataFile
	}{
		{
			diskRecord: *NewDiskRecord("first", "val1"),
			dataFile:   df1,
		},
		{
			diskRecord: *NewDiskRecord("second", "val2"),
			dataFile:   df1,
		},
		{
			diskRecord: *NewDiskRecord("third", "val3"),
			dataFile:   df2,
		},
	}

	// Write Records into appropriate files
	for _, record := range recs {
		encoded, err := record.diskRecord.encode()
		assert.NoError(t, err)

		written, err := record.dataFile.Fd.WriteString(string(encoded))
		assert.Greater(t, written, 0)
	}

	b, err := NewBitcask(WithDir(tmpDir))
	assert.NoError(t, err)

	// Assert that every Key was saved into KeyDir
	assert.Equal(t, len(b.KeyDir), len(recs))

	// Assert that KeyDirRecords match
	for _, record := range recs {
		assert.Equal(t, b.KeyDir[record.diskRecord.Key],
			KeyDirRecord{
				FileId:    record.dataFile.Id,
				ValueSize: len(record.diskRecord.Value),
				ValuePos:  0,
				Timestamp: int(record.diskRecord.Header.Timestamp),
			})
	}
}
