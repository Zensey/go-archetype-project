package main

import (
	. "bitbucket.org/Zensey/go-archetype-project/cmd/demo/types"
	"github.com/gbrlsnchs/jwt"
)

const spinsApiUri = "/api/machines/atkins-diet/spins"
const secret = "key"
const contentType = "application/json"

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
	Stops TStops   `json:"stops"`
}

func newTSpinDto(s TBaseSpin) TSpinDto {
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
	r.Total = s.GetWin()

	for _, sp := range s.GetSpins() {
		spinDto := newTSpinDto(sp)
		r.Spins = append(r.Spins, spinDto)
	}

	t := newToken(s.GetUid(), s.GetBet(), s.GetChips())
	sgn, err := t.pack()
	if err == nil {
		r.JWT = string(sgn)
	}
	return r
}
