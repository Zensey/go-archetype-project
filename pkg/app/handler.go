package app

import (
	"encoding/gob"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Zensey/slog"
)

const (
	// settings
	throttlePeriodSec = 10
	maxRequests       = 4
	ttl               = 10
	toNanoSec         = int64(time.Second)
)

type clientState struct {
	//address  string
	requests []int64
	sync.Mutex
}

type Handler struct {
	slog.Logger

	mux     *http.ServeMux
	clients map[string]*clientState
	sync.Mutex
}

func NewHandler(log slog.Logger) *Handler {
	s := &Handler{mux: http.NewServeMux()}
	s.Logger = log
	s.mux.HandleFunc("/", s.index)
	s.clients = make(map[string]*clientState, 0)

	return s
}

// This func discards all timestamps older then now - ttl
// unit of both arguments is nanosecond

func (cs *clientState) discardOld(now int64, ttl int64) int64 {

	if len(cs.requests) == 0 {
		return 0
	}
	// n=2 t=3
	// 1,1,1,2,2,2,5,5,5
	// t-(now-5)

	for i, v := range cs.requests {
		fmt.Println(">>>>", i, v, now-v, ttl)
		if now-v <= ttl {
			cs.requests = cs.requests[i:]
			break
		}
	}

	if len(cs.requests) > 0 {
		return cs.requests[0]
	}
	return 0
}

func nowNano() int64 {
	return time.Now().UnixNano()
}

func (s *Handler) throttle(ip string) clientState {
	now := nowNano()

	// global lock
	s.Lock()
	c, ok := s.clients[ip]
	if !ok {
		c = &clientState{
			//address: ip,
		}
		s.clients[ip] = c
	}
	s.Unlock()

	// lock on a specific entry
	c.Lock()
	defer c.Unlock()
	c.requests = append(c.requests, now)
	theBeginning := c.discardOld(now, throttlePeriodSec*toNanoSec)

	nReq := len(c.requests)
	if nReq > maxRequests {
		restNano := throttlePeriodSec*toNanoSec - (now - theBeginning) // calc the rest of period [nanosec]
		fmt.Println("restNano >", restNano)
		time.Sleep(time.Duration(restNano))
	}
	return *c
}

func (s *Handler) index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		return
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	fmt.Println(ip)

	cs := s.throttle(ip)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.Itoa(len(cs.requests))))
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
	//s.discardOld(nowNano(), ttlNanoSec)

	s.Info("Saving state..")
	return gob.NewEncoder(f).Encode(s.clients)
}

func (s *Handler) LoadState() error {
	s.Info("Handler > Load state")

	f, err := os.Open(dataFile)
	if err != nil {
		return err
	}
	defer f.Close()

	err = gob.NewDecoder(f).Decode(&s.clients)
	if err != nil {
		return err
	}
	s.Info("Restored recs:", len(s.clients))
	//s.discardOld(nowNano(), ttlNanoSec)

	return nil
}
