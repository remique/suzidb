package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// This test checks whether prefix in Memory works. This imitates how MVCC would get
// the keys in order to, later on, get the latest version.
func TestScanWithPrefix(t *testing.T) {
	ms := NewMemStorage()

	err := ms.Set("mytable:1::000001", "firstValue")
	assert.NoError(t, err)

	err = ms.Set("mytable:1::000002", "secondValue")
	assert.NoError(t, err)

	err = ms.Set("mytable:1::000003", "thirdValue")
	assert.NoError(t, err)

	res := ms.ScanWithPrefix("mytable:1")

	expectedRes := map[string]string{
		"mytable:1::000001": "firstValue",
		"mytable:1::000002": "secondValue",
		"mytable:1::000003": "thirdValue",
	}

	assert.Equal(t, expectedRes, res)
}
