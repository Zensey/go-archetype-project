package atkins

import . "github.com/Zensey/go-archetype-project/cmd/demo/types"

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

var PayTable = [nSymbols][nReels]int{
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

var Reels = [nReels][nReelSyms]Symbol{
	{Scatter, Sym9, Sym2, Sym4, Sym8, Sym5, Sym7, Sym9, Sym4, Sym6, Sym3, Sym8, Sym5, Sym9, Sym1, Sym3, Sym6, Sym7, Sym5, Wild, Sym8, Sym9, Sym2, Sym7, Sym5, Scatter, Sym6, Sym8, Sym4, Sym3, Sym1, Sym6},
	{Sym9, Sym3, Sym1, Sym4, Sym7, Sym9, Sym2, Sym6, Sym8, Sym1, Sym4, Sym9, Sym2, Wild, Sym6, Sym5, Sym7, Sym8, Sym4, Sym3, Scatter, Sym9, Sym6, Sym7, Sym8, Sym5, Sym3, Sym9, Sym1, Sym2, Sym7, Sym8},
	{Sym2, Sym6, Sym5, Scatter, Sym7, Sym9, Sym6, Sym2, Sym4, Sym8, Sym1, Sym3, Sym6, Sym9, Sym7, Sym4, Sym5, Sym8, Sym9, Sym3, Sym2, Sym4, Sym8, Sym7, Sym5, Wild, Sym3, Sym8, Sym6, Sym7, Sym9, Sym1},
	{Sym2, Sym7, Wild, Scatter, Sym6, Sym8, Sym7, Sym4, Sym1, Sym5, Sym8, Sym9, Sym4, Sym7, Sym6, Sym2, Sym9, Sym8, Sym3, Sym4, Sym7, Sym5, Sym6, Sym3, Sym8, Sym9, Sym5, Sym2, Sym4, Sym1, Sym9, Sym8},
	{Sym8, Scatter, Sym1, Sym2, Sym7, Sym4, Sym6, Sym8, Sym3, Sym7, Sym4, Sym2, Sym6, Sym1, Sym9, Sym5, Sym4, Sym2, Wild, Sym6, Sym3, Sym9, Sym5, Sym2, Sym8, Sym6, Sym1, Sym9, Sym4, Sym5, Sym7, Sym3},
}

var WinLines = [nLines][nReels]int{
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
