package main

type IMachineState interface {
	play() error

	getSpins() []TSpin
	getUid() string
	getBet() int
	getChips() int
	getWin() int
}
