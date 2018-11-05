package types

import (
	"encoding/json"
	"errors"
)

type TStops []int
type Symbol int
type SpinType int

type TBaseSpin struct {
	SpinType SpinType
	Stops    TStops
	Total    int
}

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
