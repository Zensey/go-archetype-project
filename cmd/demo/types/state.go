package types

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
