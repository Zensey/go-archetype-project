package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"testing"

	"bitbucket.org/Zensey/go-archetype-project/pkg/logger"
	"github.com/stretchr/testify/assert"
)

const (
	url = "http://localhost:8080" + spinsApiUri
)

func Test_Server(t *testing.T) {
	app, err := newServer()
	assert.Nil(t, err)
	err = app.start()
	assert.Nil(t, err)

	rand.Seed(1)
	lg, _ := logger.NewLogger(logger.LogLevelInfo, "test", logger.BackendConsole)
	client := http.Client{}

	tok := newToken("asfasf", 1000, 10000)
	tok.SetAlgorithm(hs256)

	makeReq := func() (respTok *TokenDto, respDto *ResponseDto, err error) {
		tokenBytes, err := tok.pack()
		if err != nil {
			return
		}

		rr := bytes.NewReader(tokenBytes)
		resp, err := client.Post(url, contentType, rr)
		if err != nil {
			return
		}
		if resp.StatusCode != 200 {
			err = errors.New(resp.Status)
			return
		}
		assert.Equal(t, 200, resp.StatusCode)

		respDto = &ResponseDto{}
		err = json.NewDecoder(resp.Body).Decode(respDto)
		if err != nil {
			return
		}

		respTok = newToken("", 0, 0) // empty token
		err = respTok.unpack([]byte(respDto.JWT))
		return
	}

	for i := 1; i <= 5; i++ {
		respTok, dto, err := makeReq()
		assert.Nil(t, err)
		if err != nil {
			break
		}

		assert.Equal(t, tok.Uid, respTok.Uid)
		assert.Equal(t, tok.Bet, respTok.Bet)
		assert.Equal(t, respTok.Chips, tok.Chips-tok.Bet+dto.Total)
		tok = respTok
		lg.Infof("#req %d chips: %v total: %d", i, tok.Chips, dto.Total)
	}

	err = app.stop()
	assert.Nil(t, err)
}
