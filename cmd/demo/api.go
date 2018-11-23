package main

import (
	. "github.com/Zensey/go-archetype-project/cmd/demo/types"
	"github.com/gbrlsnchs/jwt"
)

const (
	spinsApiUri = "/api/machines/atkins-diet/spins"
	secret      = "key"
	contentType = "application/json"
)

// TODO: consider useing RSA for more secure communication
var hs256 = jwt.NewHS256(secret)

type TokenDto struct {
	*jwt.JWT

	Uid   string `json: "uid"`
	Chips int    `json: "chips"`
	Bet   int    `json: "bet"`
}

func newToken(uid string, bet, chips int) *TokenDto {
	return &TokenDto{JWT: &jwt.JWT{},
		Uid:   uid,
		Bet:   bet,
		Chips: chips,
	}
}

func (req *TokenDto) unpack(token []byte) error {
	payload, sig, err := jwt.ParseBytes(token)
	if err != nil {
		return err
	}
	if err = hs256.Verify(payload, sig); err != nil {
		return err
	}
	if err = jwt.Unmarshal(payload, &req); err != nil {
		return err
	}
	return nil
}

func (req *TokenDto) pack() ([]byte, error) {
	payload, err := jwt.Marshal(req)
	if err != nil {
		return nil, err
	}
	token, err := hs256.Sign(payload)
	if err != nil {
		return nil, err
	}
	return token, nil
}

/////////////////////////////////////////////////////
type TSpinDto struct {
	Type  SpinType `json:"type"`
	Total int      `json:"total"`
	Stops []int    `json:"stops"`
}

func newTSpinDto(s TSpin) TSpinDto {
	return TSpinDto{
		Type:  s.SpinType,
		Stops: s.Stops,
		Total: s.Total,
	}
}

type ResponseDto struct {
	Total int        `json:"total"`
	Spins []TSpinDto `json:"spins"`
	JWT   string     `json:"jwt"`
}

func newResponseDto(s IMachineState) ResponseDto {
	r := ResponseDto{}
	bs := s.GetBaseState()
	r.Total = bs.Win

	for _, sp := range bs.GetSpins() {
		spinDto := newTSpinDto(sp)
		r.Spins = append(r.Spins, spinDto)
	}

	t := newToken(bs.Uid, bs.Bet, bs.Chips)
	sgn, err := t.pack()
	if err == nil {
		r.JWT = string(sgn)
	}
	return r
}
