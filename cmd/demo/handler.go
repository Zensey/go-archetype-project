package main

import (
	"encoding/gob"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"bitbucket.org/Zensey/go-archetype-project/pkg/logger"
)

type Handler struct {
	logger.Logger
	mux      *http.ServeMux
	requests []int64
	sync.Mutex
}

func NewHandler(log logger.Logger) *Handler {
	s := &Handler{mux: http.NewServeMux()}
	s.Logger = log
	s.mux.HandleFunc("/", s.index)
	return s
}

// This func discards all timestamps older then now - ttl (]
func (s *Handler) discardOlder(now int64, ttl int64) {
	if len(s.requests) == 0 {
		return
	}
	iFirst := 0
	for i, v := range s.requests {
		if v >= now-ttl {
			break
		}
		iFirst = i + 1
	}
	if iFirst > 0 {
		s.requests = s.requests[iFirst:]
	}
}

func (s *Handler) index(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != "/" {
		return
	}

	nReq := 0

	s.Lock()
	now := time.Now().Unix()
	s.discardOlder(now, ttlWindowSec)
	nReq = len(s.requests)
	s.requests = append(s.requests, now)
	s.Unlock()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.Itoa(nReq)))
}

func (s *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Handler) OnShutdown() error {
	s.Info("Shutdown..")
	f, err := os.Create("/tmp/codingTask.dat")
	if err != nil {
		return err
	}
	s.discardOlder(time.Now().Unix(), ttlWindowSec)
	s.Info("Saving recs:", len(s.requests))
	return gob.NewEncoder(f).Encode(s.requests)
}

func (s *Handler) BeforeStart() error {
	s.Info("Handler > Before start")
	f, err := os.Open("/tmp/codingTask.dat")
	if err != nil {
		return err
	}
	err = gob.NewDecoder(f).Decode(&s.requests)
	if err != nil {
		return err
	}
	s.Info("Restored recs:", len(s.requests))
	s.discardOlder(time.Now().Unix(), ttlWindowSec)
	return f.Close()
}
