package types

import "reflect"

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

func getType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}
