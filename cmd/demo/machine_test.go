package main

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Machine_CalcLineWin(t *testing.T) {
	s := newTState("uid", 20, 10000)

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
}

func Test_Machine_SingleSpinWithoutScatters(t *testing.T) {
	rand.Seed(0)
	s := newTState("uid", 20, 0)

	/* case w/o scatters */
	s.stops = TStops{27, 14, 3, 31, 27}
	win := s.calcSingleSpinWining(mainSpin)
	assert.Equal(t, 15, win)
}

func Test_Machine_PlayWithScatters(t *testing.T) {
	rand.Seed(0)

	/* case with scatters and free plays */
	s := newTState("uid", 1000, 10000)
	s.stops = TStops{26, 11, 21, 5, 2}
	s.isInitialized = true
	err := s.play()
	assert.Nil(t, err)
	assert.Equal(t, 16150, s.chips)
}
