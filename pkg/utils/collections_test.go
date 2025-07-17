package utils_test

import (
	"testing"

	"github.com/GoCodingX/repartners/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestDedupe(t *testing.T) {
	t.Run("ints with duplicates", func(t *testing.T) {
		input := []int{1, 2, 2, 3, 1, 4}
		expected := []int{1, 2, 3, 4}

		assert.Equal(t, expected, utils.Dedupe(input))
	})

	t.Run("strings with duplicates", func(t *testing.T) {
		input := []string{"a", "b", "a", "c", "b"}
		expected := []string{"a", "b", "c"}

		assert.Equal(t, expected, utils.Dedupe(input))
	})

	t.Run("empty slice", func(t *testing.T) {
		var input []int

		var expected []int

		assert.Equal(t, expected, utils.Dedupe(input))
	})

	t.Run("all duplicates", func(t *testing.T) {
		input := []int{5, 5, 5, 5}
		expected := []int{5}

		assert.Equal(t, expected, utils.Dedupe(input))
	})

	t.Run("no duplicates", func(t *testing.T) {
		input := []int{1, 2, 3, 4}
		expected := []int{1, 2, 3, 4}

		assert.Equal(t, expected, utils.Dedupe(input))
	})
}
