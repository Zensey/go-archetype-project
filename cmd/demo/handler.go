package main

import (
	"net/http"
)

type HandlerCtx struct {
	mux    *http.ServeMux
}

func NewHandler() *HandlerCtx {
	s := &HandlerCtx{mux: http.NewServeMux()}
	s.mux.HandleFunc("/", s.index)
	return s
}

func (s *HandlerCtx) index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}

func (s *HandlerCtx) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "example Go server")
	s.mux.ServeHTTP(w, r)
}