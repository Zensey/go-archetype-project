package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Main(t *testing.T) {
	var seq [nReels]Symbol

	seq = [nReels]Symbol{0, 0, 1, 2, 2}
	assert.Equal(t, 40, checkRepeats(seq))

	seq = [nReels]Symbol{0, 0, 1, 1, 1}
	assert.Equal(t, 1000, checkRepeats(seq))

	seq = [nReels]Symbol{0, 0, 1, 0, 1}
	assert.Equal(t, 1000, checkRepeats(seq))

	seq = [nReels]Symbol{0, 0, 0, 0, 0}
	assert.Equal(t, 5000, checkRepeats(seq))

	seq = [nReels]Symbol{0, Scatter, 0, Scatter, Scatter}
	assert.Equal(t, 5, calcScatterSum(seq))
}
