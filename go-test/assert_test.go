package go_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestEqual 使用Equal
func TestEqual(t *testing.T) {
	expected := 100
	var b = 100
	assert.Equal(t, expected, b, "")
}

// TestEqual 使用NotEqual
func TestNotEqual(t *testing.T) {
	expected := 100
	var b = 200
	var c = 300
	assert.NotEqual(t, expected, b, "")
	assert.NotEqual(t, expected, c, "")
}
func TestFalse(t *testing.T) {
	assert.False(t, 1+1 == 3, "1+1 == 3 should be false")
}
