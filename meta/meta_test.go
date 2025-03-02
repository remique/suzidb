package meta

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeRows(t *testing.T) {
	left := Row{
		"carid":   "1",
		"brandid": "1",
	}

	right := Row{
		"brandid":   "1",
		"brandname": "suzuki",
	}

	expected := Row{
		"cars.carid":       "1",
		"cars.brandid":     "1",
		"brands.brandid":   "1",
		"brands.brandname": "suzuki",
	}

	res := MergeRows(left, right, "cars", "brands")

	assert.Equal(t, expected, res)
}
