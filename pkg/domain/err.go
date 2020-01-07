package domain

func NewLogicError(text string) error {
	return &LogicError{text}
}

type LogicError struct {
	s string
}

func (e *LogicError) Error() string {
	return e.s
}
