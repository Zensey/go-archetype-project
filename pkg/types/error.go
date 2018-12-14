package types

func NewErrorLogic(text string) error {
	return &ErrorLogic{text}
}

// errorString is a trivial implementation of error.
type ErrorLogic struct {
	s string
}

func (e *ErrorLogic) Error() string {
	return e.s
}
