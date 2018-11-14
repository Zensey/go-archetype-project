package types

import (
	"encoding/json"
	"errors"
)

type (
	Symbol   int
	SpinType int
	TSymRow  []Symbol // length: nReels
)

const (
	MainSpin = SpinType(iota)
	FreeSpin
)

func (u SpinType) MarshalJSON() ([]byte, error) {
	stype := "main"
	if u != MainSpin {
		stype = "free"
	}
	return json.Marshal(stype)
}

func (u *SpinType) UnmarshalJSON(data []byte) error {
	v := ""
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	switch v {
	case "main":
		*u = MainSpin
	case "free":
		*u = FreeSpin
	default:
		return errors.New("unknown spin type " + v)
	}
	return nil
}

///////////////////////////////////////////////////////////
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

///////////////////////////////////////////////////////////
type TBaseState struct {
	Uid   string
	Bet   int
	Chips int
	Win   int // Total win

	spins []TSpin // played spins
}

func (s *TBaseState) HandleSpin(spin TSpin) {
	s.Win += spin.Total
	s.Chips += spin.Total

	s.spins = append(s.spins, spin)
}

func (s *TBaseState) WithdrawBet() {
	s.Chips = s.Chips - s.Bet
}

func (s *TBaseState) GetSpins() []TSpin {
	return s.spins
}
