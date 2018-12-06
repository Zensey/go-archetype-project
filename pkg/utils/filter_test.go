package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Filter(t *testing.T) {
	var badWords = []string{"fee", "nee", "cruul", "leent"}

	hasBad := DetectBadWords("I really love the product and will recommend!", badWords)
	assert.Equal(t, false, hasBad)
}
