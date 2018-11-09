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

const (
	throttlePeriodSec = 5
	maxRequests       = 20
	nanoSecMultiplier = 1e9
)

type clientState struct {
	address      string
	theBeginning int64
	nRequests    int

	sync.Mutex
}

type Handler struct {
	logger.Logger

	mux      *http.ServeMux
	requests []int64
	clients  map[string]*clientState

	sync.Mutex
}

func NewHandler(log logger.Logger) *Handler {
	s := &Handler{mux: http.NewServeMux()}
	s.Logger = log
	s.mux.HandleFunc("/", s.index)
	s.clients = make(map[string]*clientState, 0)

	return s
}

// This func discards all timestamps older then now - ttl
// unit of both arguments is nanosecond

func (s *Handler) discardOld(now int64, ttl int64) {
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

func nowNano() int64 {
	return time.Now().UnixNano()
}

func (s *Handler) throttle(ip string) {
	now := nowNano()

	// global lock
	s.Lock()
	c := s.clients[ip]
	if c == nil {
		c = &clientState{
			address:      ip,
			theBeginning: now,
			nRequests:    0,
		}
		s.clients[ip] = c
	}
	s.Unlock()

	// lock on a specific entry
	c.Lock()
	c.nRequests++
	restNano := c.theBeginning + throttlePeriodSec*nanoSecMultiplier - now // calc the rest of period [nanosec]
	//s.Info("throttle", ip, c.nRequests, c.theBeginning, now, restNano,)

	if restNano > 0 && c.nRequests > maxRequests {
		time.Sleep(time.Duration(restNano) * time.Nanosecond)
		c.nRequests = 1
		c.theBeginning = nowNano()
		s.Info("throttle reset for", ip)
	}
	c.Unlock()
}

func (s *Handler) index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		return
	}
	now := nowNano()

	//ip, _, err := net.SplitHostPort(r.RemoteAddr)
	//if err != nil {
	//	s.Error(err)
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}
	ip := r.URL.Query()["ip"][0]
	s.throttle(ip)

	s.Lock()
	s.discardOld(now, ttlNanoSec)
	nReq := len(s.requests)
	s.requests = append(s.requests, now)
	s.Unlock()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.Itoa(nReq)))
}

func (s *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Handler) SaveState() error {
	s.Info("Shutdown..")

	f, err := os.Create(dataFile)
	if err != nil {
		return err
	}
	s.discardOld(nowNano(), ttlNanoSec)
	s.Info("Saving recs:", len(s.requests))
	return gob.NewEncoder(f).Encode(s.requests)
}

//LoadState
func (s *Handler) LoadState() error {
	s.Info("Handler > Load state")

	f, err := os.Open(dataFile)
	if err != nil {
		return err
	}
	defer f.Close()

	err = gob.NewDecoder(f).Decode(&s.requests)
	if err != nil {
		return err
	}
	s.Info("Restored recs:", len(s.requests))
	s.discardOld(nowNano(), ttlNanoSec)
	return nil
}
