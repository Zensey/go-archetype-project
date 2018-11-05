package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"bitbucket.org/Zensey/go-archetype-project/pkg/logger"
	"github.com/gbrlsnchs/jwt"
)

type Handler struct {
	logger.Logger
	mux *http.ServeMux
}

func NewHandler(log logger.Logger) *Handler {
	s := &Handler{mux: http.NewServeMux()}
	s.Logger = log
	s.mux.HandleFunc(spinsApiUri, s.spins)
	return s
}

func (s *Handler) spins(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != spinsApiUri {
		return
	}
	tok, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.Info("error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	req := TokenDto{JWT: &jwt.JWT{}}
	err = req.unpack(tok)
	if err != nil {
		s.Info("error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	st := newTState(req.Uid, req.Bet, req.Chips)
	err = st.play()
	if err != nil {
		s.Info("error", err)
	}
	resp := newResponseDto(st)
	w.Header().Set("Content-Type", contentType)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		s.Info("error", err)
	}
}

func (s *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
