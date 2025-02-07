package bitcask

import (
	"os"
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

func TestGlob(t *testing.T) {
	_, err1 := NewDataFile(".", 1)
	_, err2 := NewDataFile(".", 2)
	_, err3 := NewDataFile(".", 3)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NoError(t, err3)

	res, err := glob(".")
	assert.NoError(t, err)
	assert.Equal(t, []string{"1.db", "2.db", "3.db"}, res)

	defer os.Remove("./1.db")
	defer os.Remove("./2.db")
	defer os.Remove("./3.db")
}

func TestNewActiveFileId(t *testing.T) {
	_, err1 := NewDataFile(".", 1)
	_, err2 := NewDataFile(".", 2)
	_, err3 := NewDataFile(".", 3)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NoError(t, err3)

	res, err := generateNewActiveFileId(".")
	assert.NoError(t, err)
	assert.Equal(t, 4, res)

	defer os.Remove("./1.db")
	defer os.Remove("./2.db")
	defer os.Remove("./3.db")
}
