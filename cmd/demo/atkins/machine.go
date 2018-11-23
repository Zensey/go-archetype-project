package atkins

import (
	"errors"
	"math/rand"

	. "github.com/Zensey/go-archetype-project/cmd/demo/types"
)

func getPay(sym Symbol, repeats int) int {
	if repeats > 0 {
		return PayTable[sym][nReels-repeats]
	}
	return 0
}

func getReelSymSeq(r int, mid int) (ret [nRows]Symbol) {
	reel := Reels[r]
	ind := (nReelSyms + mid - 1 - 1) % nReelSyms

	for i := 0; i < nRows; i++ {
		ret[i] = reel[ind]
		ind = (ind + 1) % nReelSyms
	}
	return
}

/////////////////////////////////////////////////////////////
type TAtkinsState struct {
	TBaseState
	TSpin // current spin

	scatters      int
	freeRuns      int
	isInitialized bool // for debug purposes
}

func NewAtkins(uid string, bet, chips int) *TAtkinsState {
	return &TAtkinsState{
		TBaseState: TBaseState{
			Uid:   uid,
			Bet:   bet,
			Chips: chips,
		},
		TSpin: NewTBaseSpin(nReels),
	}
}

func (s *TAtkinsState) calcLineWin() int {
	calcSum := func(firstSym Symbol) int {
		repeats := 1
		for i := 1; i < nReels; i++ {
			if s.Row[i] == firstSym || (s.Row[i] == Wild && firstSym != Scatter) { // wild cannot subst scatter
				repeats++
				continue
			}
			break
		}
		return getPay(firstSym, repeats)
	}

	firstSym := s.Row[0]
	sum := calcSum(firstSym)

	// alternative case: wilds are first; treat them as first-non-wild symbol
	if firstSym == Wild {
		for i := 1; i < nReels; i++ {
			if s.Row[i] != Wild {
				sumAlt := calcSum(s.Row[i])
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

func (s *TAtkinsState) calcSingleSpinWining(spinType SpinType) int {
	s.SpinType = spinType
	if spinType == MainSpin {
		s.WithdrawBet()
	}

	coins := int(s.Bet / nLines)
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
			s.Row[col] = T[col][row]
		}
		sum += s.calcLineWin()
	}
	sum += getPay(Scatter, scatters) // scatter pay
	if spinType == FreeSpin {
		sum = sum * 3 // in free games all wins are tripled
	}

	s.Total = sum * coins
	s.HandleSpin(s.TSpin)

	return s.Total
}

func (s *TAtkinsState) stopRandom() {
	for r := 0; r < nReels; r++ {
		s.Stops[r] = rand.Intn(nReelSyms)
	}
}

func (s *TAtkinsState) Play() error {
	if s.Chips < s.Bet {
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

func (s *TAtkinsState) GetBaseState() TBaseState {
	return s.TBaseState
}
