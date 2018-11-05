package types

type IMachineState interface {
	Play() error

	GetSpins() []TBaseSpin
	GetUid() string
	GetBet() int
	GetChips() int
	GetWin() int
}
