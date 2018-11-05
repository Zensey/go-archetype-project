package main

import (
	"encoding/json"
	"errors"
)

type TStops []int
type Symbol int
type SpinType int

const (
	mainSpin = SpinType(iota)
	freeSpin
)

func (u SpinType) MarshalJSON() ([]byte, error) {
	stype := "main"
	if u != mainSpin {
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
		*u = mainSpin
	case "free":
		*u = freeSpin
	default:
		return errors.New("unknown spin type " + v)
	}
	return nil
}
