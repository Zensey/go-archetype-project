package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Main(t *testing.T) {
	s := newTState(20, 10000)

	s.seq = TSymRow{0, 0, 1, 2, 2}
	assert.Equal(t, 40, s.calcLineWin())

	s.seq = TSymRow{0, 0, 1, 1, 1}
	assert.Equal(t, 1000, s.calcLineWin())

	s.seq = TSymRow{0, 0, 1, 0, 1}
	assert.Equal(t, 1000, s.calcLineWin())

	s.seq = TSymRow{0, 0, 0, 0, 0}
	assert.Equal(t, 5000, s.calcLineWin())

	// wild may substitute for any symbol, except the scale
	s.seq = TSymRow{0, 0, 10, 10, 10}
	assert.Equal(t, 5, s.calcLineWin())

	// wild may substitute for any symbol, except the scale
	s.seq = TSymRow{10, 10, 10, 0, 1}
	assert.Equal(t, 5, s.calcLineWin())

	/* case w/o scatters */
	s.stops = TStops{27, 14, 3, 31, 27}
	win := s.calcSingleSpinWining(1)
	assert.Equal(t, 15, win)

	/* case with scatters and free plays */
	s2 := newTState(1000, 10000)
	s2.stops = TStops{26, 11, 21, 5, 2}
	s2.play()
	assert.Equal(t, 49750, s2.chips)

}
