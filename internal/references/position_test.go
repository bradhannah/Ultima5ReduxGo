package references

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WithinN_One(t *testing.T) {
	p1 := Position{X: 0, Y: 0}
	assert.True(t, p1.IsWithinN(&Position{X: 0, Y: 0}, 1))
	assert.True(t, p1.IsWithinN(&Position{X: 1, Y: 1}, 1))

	assert.False(t, p1.IsWithinN(&Position{X: 2, Y: 2}, 1))
	assert.False(t, p1.IsWithinN(&Position{X: 2, Y: 1}, 1))
	assert.False(t, p1.IsWithinN(&Position{X: 1, Y: 2}, 1))
}

func Test_WithinN_OneHundred(t *testing.T) {
	p1 := Position{X: 100, Y: 100}

	assert.True(t, p1.IsWithinN(&Position{X: 100, Y: 100}, 5))
	assert.True(t, p1.IsWithinN(&Position{X: 50, Y: 50}, 51))

	assert.False(t, p1.IsWithinN(&Position{X: 2, Y: 2}, 1))
	assert.False(t, p1.IsWithinN(&Position{X: 2, Y: 1}, 1))
	assert.False(t, p1.IsWithinN(&Position{X: 1, Y: 2}, 1))
}
