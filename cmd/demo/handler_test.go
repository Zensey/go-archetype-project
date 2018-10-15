package main

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"testing"
	"time"

	"bitbucket.org/Zensey/go-archetype-project/pkg/logger"
	"github.com/stretchr/testify/assert"
)

const (
	nWorker   = 4
	nRequests = 1000
	url       = "http://localhost:8080/"
)

func Test_DiscardOld(t *testing.T) {
	l, _ := logger.NewLogger(logger.LogLevelInfo, "test", logger.BackendConsole)

	s := NewHandler(l)
	s.requests = []int64{1, 2, 3, 4, 5, 6, 7, 8, 9}
	s.discardOld(11, 3)
	assert.Equal(t, true, reflect.DeepEqual(s.requests, []int64{8, 9}))
	s.Info(s.requests)

	s.requests = []int64{1, 2, 3, 4, 5, 6, 7, 8, 9}
	s.discardOld(1, 3)
	assert.Equal(t, true, reflect.DeepEqual(s.requests, []int64{1, 2, 3, 4, 5, 6, 7, 8, 9}))
	s.Info(s.requests)

	s.requests = []int64{1, 2, 3, 4, 5, 6, 7, 8, 9}
	s.discardOld(-1, 3)
	assert.Equal(t, true, reflect.DeepEqual(s.requests, []int64{1, 2, 3, 4, 5, 6, 7, 8, 9}))
	s.Info(s.requests)

	s.requests = []int64{1, 2, 3, 4, 5, 6, 7, 8, 9}
	s.discardOld(10, 0)
	assert.Equal(t, true, reflect.DeepEqual(s.requests, []int64{}))
	s.Info(s.requests)
}

type reqResult struct {
	val int64
	err error
}

func request(res chan reqResult, client *http.Client, l logger.Logger, requests int) {
	for i := 0; i < requests; i++ {
		resp, err := client.Get(url)
		if err != nil {
			res <- reqResult{err: err}
			return
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			res <- reqResult{err: err}
			return
		}
		val, err := strconv.ParseInt(string(b), 10, 64)
		res <- reqResult{val, err}
	}
}

func runRequestsTest(ch chan reqResult, client *http.Client, l logger.Logger, t *testing.T) int64 {
	for n := 0; n < nWorker; n++ {
		go request(ch, client, l, nRequests)
	}
	maxI := int64(0)
	for i := 0; i < nRequests*nWorker; i++ {
		r := <-ch
		assert.Nil(t, r.err)
		if maxI < r.val {
			maxI = r.val
		}
	}
	return maxI
}

func Test_Server(t *testing.T) {
	lg, _ := logger.NewLogger(logger.LogLevelInfo, "test", logger.BackendConsole)
	client := &http.Client{}
	ch := make(chan reqResult, nWorker)

	app, err := newApp()
	assert.Nil(t, err)

	app.Info("Clearing state store")
	err = app.saveState() // clear state
	assert.Nil(t, err)

	// run requests 3 times with 5 second pause
	for i := 1; i <= 3; i++ {
		err = app.start()
		assert.Nil(t, err)
		maxI := runRequestsTest(ch, client, lg, t)
		assert.Equal(t, int64(nRequests*nWorker*i)-1, maxI)

		err = app.stop()
		assert.Nil(t, err)

		pause := 5 * time.Second
		app.Info("Waiting", pause)
		time.Sleep(pause)
	}
}
