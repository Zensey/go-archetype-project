package types

import (
	"encoding/json"
	"errors"
)

type (
	Symbol   int
	SpinType int
	TSymRow  []Symbol // length: nReels
)

const (
	MainSpin = SpinType(iota)
	FreeSpin
)

func (u SpinType) MarshalJSON() ([]byte, error) {
	stype := "main"
	if u != MainSpin {
		stype = "free"
	}
	return json.Marshal(stype)
}

func (u *SpinType) UnmarshalJSON(data []byte) error {
	v := ""
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	switch v {
	case "main":
		*u = MainSpin
	case "free":
		*u = FreeSpin
	default:
		return errors.New("unknown spin type " + v)
	}
	return nil
}
