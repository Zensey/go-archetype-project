package types

type TSpin struct {
	Stops []int // position of a reel
	Total int
	Row   TSymRow

	SpinType SpinType
}

func NewTBaseSpin(nReels int) TSpin {
	return TSpin{
		Row:   make([]Symbol, nReels),
		Stops: make([]int, nReels),
	}
}
