package compare

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMin(t *testing.T) {
	tests := []struct {
		name   string
		vals   []int
		expect int
	}{
		{"empty", []int{}, 0},
		{"single", []int{5}, 5},
		{"normal", []int{3, 1, 4, 1, 5, 9, 2, 6}, 1},
		{"negative", []int{-3, -1, -4, -1, -5}, -5},
		{"mixed", []int{-1, 0, 1}, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Min(tt.vals...)
			assert.Equal(t, tt.expect, result)
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		name   string
		vals   []int
		expect int
	}{
		{"empty", []int{}, 0},
		{"single", []int{5}, 5},
		{"normal", []int{3, 1, 4, 1, 5, 9, 2, 6}, 9},
		{"negative", []int{-3, -1, -4, -1, -5}, -1},
		{"mixed", []int{-1, 0, 1}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Max(tt.vals...)
			assert.Equal(t, tt.expect, result)
		})
	}
}

func TestMinMaxFloat(t *testing.T) {
	t.Run("float64", func(t *testing.T) {
		assert.Equal(t, 1.5, Min(3.14, 2.71, 1.5, 4.2))
		assert.Equal(t, 4.2, Max(3.14, 2.71, 1.5, 4.2))
	})

	t.Run("string", func(t *testing.T) {
		assert.Equal(t, "apple", Min("banana", "apple", "cherry"))
		assert.Equal(t, "cherry", Max("banana", "apple", "cherry"))
	})
}

func TestContains(t *testing.T) {
	tests := []struct {
		name   string
		vals   []int
		target int
		expect bool
	}{
		{"empty", []int{}, 1, false},
		{"found", []int{1, 2, 3}, 2, true},
		{"not found", []int{1, 2, 3}, 4, false},
		{"first", []int{1, 2, 3}, 1, true},
		{"last", []int{1, 2, 3}, 3, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Contains(tt.vals, tt.target)
			assert.Equal(t, tt.expect, result)
		})
	}
}

func TestContainsString(t *testing.T) {
	assert.True(t, Contains([]string{"a", "b", "c"}, "b"))
	assert.False(t, Contains([]string{"a", "b", "c"}, "d"))
}

func TestUnique(t *testing.T) {
	tests := []struct {
		name   string
		vals   []int
		expect []int
	}{
		{"empty", []int{}, []int{}},
		{"single", []int{1}, []int{1}},
		{"no duplicates", []int{1, 2, 3}, []int{1, 2, 3}},
		{"with duplicates", []int{1, 2, 2, 3, 3, 3}, []int{1, 2, 3}},
		{"preserves order", []int{3, 1, 2, 1, 3}, []int{3, 1, 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Unique(tt.vals)
			assert.Equal(t, tt.expect, result)
		})
	}
}

func TestUniqueString(t *testing.T) {
	result := Unique([]string{"a", "b", "a", "c", "b"})
	assert.Equal(t, []string{"a", "b", "c"}, result)
}

func TestChunk(t *testing.T) {
	tests := []struct {
		name   string
		vals   []int
		size   int
		expect [][]int
	}{
		{"empty", []int{}, 2, nil},
		{"zero size", []int{1, 2, 3}, 0, nil},
		{"negative size", []int{1, 2, 3}, -1, nil},
		{"size larger than len", []int{1, 2}, 5, [][]int{{1, 2}}},
		{"exact size", []int{1, 2, 3}, 3, [][]int{{1, 2, 3}}},
		{"normal", []int{1, 2, 3, 4, 5}, 2, [][]int{{1, 2}, {3, 4}, {5}}},
		{"normal larger", []int{1, 2, 3, 4, 5, 6, 7}, 3, [][]int{{1, 2, 3}, {4, 5, 6}, {7}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Chunk(tt.vals, tt.size)
			assert.Equal(t, tt.expect, result)
		})
	}
}

func TestChunkString(t *testing.T) {
	result := Chunk([]string{"a", "b", "c", "d", "e"}, 2)
	assert.Equal(t, [][]string{{"a", "b"}, {"c", "d"}, {"e"}}, result)
}
