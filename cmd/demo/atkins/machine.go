package atkins

import (
	"errors"
	"math/rand"

	. "github.com/Zensey/go-archetype-project/cmd/demo/types"
)

const (
	Atkins = Symbol(iota)
	Steak
	Ham
	BuffaloWings
	Sausage
	Eggs
	Butter
	Cheese
	Bacon
	Mayonnaise
	Scatter

	// Aliases
	Wild = Atkins
	Sym1 = Steak
	Sym2 = Ham
	Sym3 = BuffaloWings
	Sym4 = Sausage
	Sym5 = Eggs
	Sym6 = Butter
	Sym7 = Cheese
	Sym8 = Bacon
	Sym9 = Mayonnaise
)
const (
	nSymbols  = Scatter + 1
	nReelSyms = 32
	nReels    = 5
	nRows     = 3 // n of visible rows (horizontal)
	nLines    = 20
)

type (
	TReel    [nReelSyms]Symbol
	TPays    [nReels]int
	TWinLine [nReels]int
	TSymRow  [nReels]Symbol

	TSpin struct {
		TBaseSpin
		seq      TSymRow
		scatters int
	}

	TState struct {
		TSpin // current spin

		freeRuns int
		spins    []TSpin // played spins

		uid   string
		bet   int
		chips int
		win   int // Total win

		isInitialized bool // for debug purposes
	}
)

var PayTable = [nSymbols]TPays{
	/*             5     4    3   2  1   // n repeats in a row in a given line    */
	Atkins:       {5000, 500, 50, 5, 0},
	Steak:        {1000, 200, 40, 3, 0},
	Ham:          {500, 150, 30, 2, 0},
	BuffaloWings: {300, 100, 25, 2, 0},
	Sausage:      {200, 75, 20, 0, 0},
	Eggs:         {200, 75, 20, 0, 0},
	Butter:       {100, 50, 15, 0, 0},
	Cheese:       {100, 50, 15, 0, 0},
	Bacon:        {50, 25, 10, 0, 0},
	Mayonnaise:   {50, 25, 10, 0, 0},
	Scatter:      {100, 25, 5, 0, 0},
}

func getPay(sym Symbol, repeats int) int {
	if repeats > 0 {
		return PayTable[sym][nReels-repeats]
	}
	return 0
}

var Reels = [nReels]TReel{
	{Scatter, Sym9, Sym2, Sym4, Sym8, Sym5, Sym7, Sym9, Sym4, Sym6, Sym3, Sym8, Sym5, Sym9, Sym1, Sym3, Sym6, Sym7, Sym5, Wild, Sym8, Sym9, Sym2, Sym7, Sym5, Scatter, Sym6, Sym8, Sym4, Sym3, Sym1, Sym6},
	{Sym9, Sym3, Sym1, Sym4, Sym7, Sym9, Sym2, Sym6, Sym8, Sym1, Sym4, Sym9, Sym2, Wild, Sym6, Sym5, Sym7, Sym8, Sym4, Sym3, Scatter, Sym9, Sym6, Sym7, Sym8, Sym5, Sym3, Sym9, Sym1, Sym2, Sym7, Sym8},
	{Sym2, Sym6, Sym5, Scatter, Sym7, Sym9, Sym6, Sym2, Sym4, Sym8, Sym1, Sym3, Sym6, Sym9, Sym7, Sym4, Sym5, Sym8, Sym9, Sym3, Sym2, Sym4, Sym8, Sym7, Sym5, Wild, Sym3, Sym8, Sym6, Sym7, Sym9, Sym1},
	{Sym2, Sym7, Wild, Scatter, Sym6, Sym8, Sym7, Sym4, Sym1, Sym5, Sym8, Sym9, Sym4, Sym7, Sym6, Sym2, Sym9, Sym8, Sym3, Sym4, Sym7, Sym5, Sym6, Sym3, Sym8, Sym9, Sym5, Sym2, Sym4, Sym1, Sym9, Sym8},
	{Sym8, Scatter, Sym1, Sym2, Sym7, Sym4, Sym6, Sym8, Sym3, Sym7, Sym4, Sym2, Sym6, Sym1, Sym9, Sym5, Sym4, Sym2, Wild, Sym6, Sym3, Sym9, Sym5, Sym2, Sym8, Sym6, Sym1, Sym9, Sym4, Sym5, Sym7, Sym3},
}

var WinLines = [nLines]TWinLine{
	{2, 2, 2, 2, 2},
	{1, 1, 1, 1, 1},
	{3, 3, 3, 3, 3},
	{1, 2, 3, 2, 1},
	{3, 2, 1, 2, 3},
	{2, 1, 1, 1, 2},
	{2, 3, 3, 3, 2},
	{1, 1, 2, 3, 3},
	{3, 3, 2, 1, 1},
	{2, 1, 2, 3, 2},
	{2, 3, 2, 1, 2},
	{1, 2, 2, 2, 1},
	{3, 2, 2, 2, 3},
	{1, 2, 1, 2, 1},
	{3, 2, 3, 2, 3},
	{2, 2, 1, 2, 2},
	{2, 2, 3, 2, 2},
	{1, 1, 3, 1, 1},
	{3, 3, 1, 3, 3},
	{1, 3, 3, 3, 1},
}

func getReelSymSeq(r int, mid int) (ret [3]Symbol) {
	reel := Reels[r]
	ind := (nReelSyms + mid - 1 - 1) % nReelSyms

	for i := 0; i < nRows; i++ {
		ret[i] = reel[ind]
		ind = (ind + 1) % nReelSyms
	}
	return
}

/////////////////////////////////////////////////////////////

func NewAtkins(uid string, bet, chips int) *TState {
	return &TState{
		uid:   uid,
		bet:   bet,
		chips: chips,
	}
}

func (s TState) GetSpins() (res []TBaseSpin) {
	for _, v := range s.spins {
		res = append(res, v.TBaseSpin)
	}
	return
}

func (s TState) GetUid() string {
	return s.uid
}

func (s TState) GetBet() int {
	return s.bet
}

func (s TState) GetChips() int {
	return s.chips
}

func (s TState) GetWin() int {
	return s.win
}

func (s *TState) calcLineWin() int {
	calcSum := func(firstSym Symbol) int {
		repeats := 1
		for i := 1; i < nReels; i++ {
			if s.seq[i] == firstSym || (s.seq[i] == Wild && firstSym != Scatter) { // wild cannot subst scatter
				repeats++
				continue
			}
			break
		}
		return getPay(firstSym, repeats)
	}

	firstSym := s.seq[0]
	sum := calcSum(firstSym)

	// alternative case: wilds are first; treat them as first-non-wild symbol
	if firstSym == Wild {
		for i := 1; i < nReels; i++ {
			if s.seq[i] != Wild {
				sumAlt := calcSum(s.seq[i])
				if sumAlt > sum { // highest win pays
					sum = sumAlt
					break
				}
				break
			}
		}
	}
	return sum
}

func (s *TState) calcSingleSpinWining(spinType SpinType) int {
	if spinType == MainSpin {
		s.chips = s.chips - s.bet
	}
	coins := int(s.bet / nLines)
	scatters := 0

	var T [nReels][nRows]Symbol // visible grid

	for col := 0; col < nReels; col++ {
		seq := getReelSymSeq(col, s.Stops[col])
		T[col] = seq

		for r := 0; r < nRows; r++ {
			if seq[r] == Scatter {
				scatters++
			}
		}
	}
	s.scatters = scatters
	if scatters >= 3 {
		s.freeRuns += 10
	}

	sum := 0
	for l := 0; l < nLines; l++ {
		payLine := WinLines[l]
		for col := 0; col < nReels; col++ {
			row := payLine[col] - 1
			s.seq[col] = T[col][row]
		}
		sum += s.calcLineWin()
	}
	sum += getPay(Scatter, scatters) // scatter pay

	if spinType == FreeSpin {
		sum = sum * 3 // in free games all wins are tripled
	}
	s.Total = sum * coins

	s.win += s.Total
	s.chips += s.Total
	s.SpinType = spinType
	s.spins = append(s.spins, s.TSpin)

	return s.Total
}

func (s *TState) stopRandom() {
	s.Stops = make(TStops, nReels)
	for r := 0; r < nReels; r++ {
		s.Stops[r] = rand.Intn(32)
	}
}

func (s *TState) Play() error {
	if s.chips < s.bet {
		return errors.New("insufficient chips")
	}
	if !s.isInitialized {
		s.stopRandom()
	}
	s.calcSingleSpinWining(MainSpin)
	for s.freeRuns > 0 {
		s.stopRandom()
		s.calcSingleSpinWining(FreeSpin)
		s.freeRuns--
	}
	return nil
}
