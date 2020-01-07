package domain

func NewLogicError(text string) error {
	return &LogicError{text}
}

// errorString is a trivial implementation of error.
type LogicError struct {
	s string
}

func (e *LogicError) Error() string {
	return e.s
}
