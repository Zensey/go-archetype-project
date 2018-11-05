package main

import (
	"bytes"
	"encoding/json"
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

	makeReq := func() (*TokenDto, *ResponseDto) {
		tokenBytes, err := tok.pack()
		assert.Nil(t, err)

		rr := bytes.NewReader(tokenBytes)
		resp, err := client.Post(url, contentType, rr)
		assert.Nil(t, err)

		respDto := &ResponseDto{}
		err = json.NewDecoder(resp.Body).Decode(&respDto)
		assert.Nil(t, err)

		respTok := newToken("", 0, 0)
		err = respTok.unpack([]byte(respDto.JWT))
		assert.Nil(t, err)
		return respTok, respDto
	}

	for i := 1; i <= 5; i++ {
		respTok, dto := makeReq()
		//lg.Infof("#req %d  respTok %v ", i, dto)

		assert.Equal(t, tok.Uid, respTok.Uid)
		assert.Equal(t, tok.Bet, respTok.Bet)
		assert.Equal(t, respTok.Chips, tok.Chips-tok.Bet+dto.Total)

		tok = respTok
		lg.Infof("#req %d chips: %v total: %d", i, tok.Chips, dto.Total)
	}

	err = app.stop()
	assert.Nil(t, err)
}
