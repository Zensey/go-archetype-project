package types

type IMachineState interface {
	Play() error
	GetBaseState() TBaseState
}
