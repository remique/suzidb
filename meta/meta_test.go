package meta

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeRows2WithPrefixes(t *testing.T) {
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

	res := MergeRows2(
		WithMergeRow(left, "cars"),
		WithMergeRow(right, "brands"),
	)

	assert.Equal(t, expected, res)
}

func TestMergeRows2NoPrefixes(t *testing.T) {
	left := Row{
		"carid":   "1",
		"brandid": "1",
	}

	right := Row{
		"brandid":   "1",
		"brandname": "suzuki",
	}

	expected := Row{
		"carid":     "1",
		"brandid":   "1",
		"brandname": "suzuki",
	}

	res := MergeRows2(
		WithMergeRow(left, ""),
		WithMergeRow(right, ""),
	)

	assert.Equal(t, expected, res)
}
