package main

import (
	"fmt"
)

//           Symbol Five Four Three Two One SymbolGen
// 1         Atkins 5000  500    50   5   0      Wild
// 2          Steak 1000  200    40   3   0      Sym1
// 3            Ham  500  150    30   2   0      Sym2
// 4  Buffalo Wings  300  100    25   2   0      Sym3
// 5        Sausage  200   75    20   0   0      Sym4
// 6           Eggs  200   75    20   0   0      Sym5
// 7         Butter  100   50    15   0   0      Sym6
// 8         Cheese  100   50    15   0   0      Sym7
// 9          Bacon   50   25    10   0   0      Sym8
// 10    Mayonnaise   50   25    10   0   0      Sym9

type Symbol int

const (
	Atkins = Symbol(iota)
	Steak
	Ham
	Drumstick
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
	Sym3 = Drumstick
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

type Pays [nReels]int
type Line [nReels]int
type Reel [nReelSyms]Symbol
type Stops [nReels]int

var PayTable = [nSymbols]Pays{
	Atkins:     {5000, 500, 50, 5, 0},
	Steak:      {1000, 200, 40, 3, 0},
	Ham:        {500, 150, 30, 2, 0},
	Drumstick:  {300, 100, 25, 2, 0},
	Sausage:    {200, 75, 20, 0, 0},
	Eggs:       {200, 75, 20, 0, 0},
	Butter:     {100, 50, 15, 0, 0},
	Cheese:     {100, 50, 15, 0, 0},
	Bacon:      {50, 25, 10, 0, 0},
	Mayonnaise: {50, 25, 10, 0, 0},
	Scatter:    {100, 25, 5, 0, 0},
}

var Reels = [5]Reel{
	{Scatter, Sym9, Sym2, Sym4, Sym8, Sym5, Sym7, Sym9, Sym4, Sym6, Sym3, Sym8, Sym5, Sym9, Sym1, Sym3, Sym6, Sym7, Sym5, Wild, Sym8, Sym9, Sym2, Sym7, Sym5, Scatter, Sym6, Sym8, Sym4, Sym3, Sym1, Sym6},
	{Sym9, Sym3, Sym1, Sym4, Sym7, Sym9, Sym2, Sym6, Sym8, Sym1, Sym4, Sym9, Sym2, Wild, Sym6, Sym5, Sym7, Sym8, Sym4, Sym3, Scatter, Sym9, Sym6, Sym7, Sym8, Sym5, Sym3, Sym9, Sym1, Sym2, Sym7, Sym8},
	{Sym2, Sym6, Sym5, Scatter, Sym7, Sym9, Sym6, Sym2, Sym4, Sym8, Sym1, Sym3, Sym6, Sym9, Sym7, Sym4, Sym5, Sym8, Sym9, Sym3, Sym2, Sym4, Sym8, Sym7, Sym5, Wild, Sym3, Sym8, Sym6, Sym7, Sym9, Sym1},
	{Sym2, Sym7, Wild, Scatter, Sym6, Sym8, Sym7, Sym4, Sym1, Sym5, Sym8, Sym9, Sym4, Sym7, Sym6, Sym2, Sym9, Sym8, Sym3, Sym4, Sym7, Sym5, Sym6, Sym3, Sym8, Sym9, Sym5, Sym2, Sym4, Sym1, Sym9, Sym8},
	{Sym8, Scatter, Sym1, Sym2, Sym7, Sym4, Sym6, Sym8, Sym3, Sym7, Sym4, Sym2, Sym6, Sym1, Sym9, Sym5, Sym4, Sym2, Wild, Sym6, Sym3, Sym9, Sym5, Sym2, Sym8, Sym6, Sym1, Sym9, Sym4, Sym5, Sym7, Sym3},
}

var PayLines = [nLines]Line{
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
	s := mid - 1 - 1
	if s < 0 {
		s = nReelSyms + s
	}
	for i := 0; i < nRows; i++ {
		ret[i] = reel[s]
		s = (s + 1) % nReelSyms
	}
	return
}

func checkRepeats(seq [nReels]Symbol) int {
	calcSum := func(firstSym Symbol) int {
		repeats := 1
		for i := 1; i < nReels; i++ {
			if seq[i] != firstSym && seq[i] != Wild {
				break
			}
			repeats++
		}
		return PayTable[firstSym][nReels-repeats]
	}

	firstSym := seq[0]
	sum := calcSum(firstSym)

	// alternative case: when Wilds are first, substitute them with first not Wild
	if firstSym == Wild {
		for i := 1; i < nReels; i++ {
			if seq[i] != Wild {
				sum2 := calcSum(seq[i])
				if sum2 > sum {
					//fmt.Println("sum2", sum2)
					sum = sum2
					break
				}
				break
			}
		}
	}
	return sum
}

func calcScatterSum(seq [nReels]Symbol) int {
	scatters := 0
	for i := 0; i < nReels; i++ {
		if seq[i] == Scatter {
			scatters++
		}
	}
	return PayTable[Scatter][nReels-scatters]
}

func calcSum(aRand Stops) int {
	var M [5][3]Symbol

	for r := 0; r < nReels; r++ {
		M[r] = getReelSymSeq(r, aRand[r])
	}
	fmt.Println("rand", aRand)
	//M[0] = getReelSymSeq(0, m)
	//fmt.Println("M", M)

	sum := 0
	for l := 0; l < nLines; l++ {
		var seq [nReels]Symbol
		payLine := PayLines[l]
		for col := 0; col < nReels; col++ {
			row := payLine[col] - 1
			seq[col] = M[col][row]
		}
		fmt.Println("> ", seq, checkRepeats(seq))
		sum += checkRepeats(seq)
	}
	return sum
}

func play() {
	//for m := 0; m < 31; m++ {
	//	calcSum(m)
	//}

	var aRand [5]int
	//for r := 0; r < nReels; r++ {
	//	aRand[r] = rand.Intn(32)
	//}
	aRand = [5]int{27, 14, 3, 31, 27}
	win := calcSum(aRand) + calcScatterSum(aRand)
	fmt.Println("> ", win)
}
