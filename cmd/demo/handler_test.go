package main

import (
	"reflect"
	"testing"

	"bitbucket.org/Zensey/go-archetype-project/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func Test_Discard(t *testing.T) {
	l, _ := logger.NewLogger(logger.LogLevelInfo, "test", logger.BackendConsole)

	s := NewHandler(l)
	s.requests = []int64{1, 2, 3, 4, 5, 6, 7, 8, 9}
	s.discardOlder(11, 3)
	assert.Equal(t, true, reflect.DeepEqual(s.requests, []int64{8, 9}))
	s.Info(s.requests)

	s.requests = []int64{1, 2, 3, 4, 5, 6, 7, 8, 9}
	s.discardOlder(1, 3)
	assert.Equal(t, true, reflect.DeepEqual(s.requests, []int64{1, 2, 3, 4, 5, 6, 7, 8, 9}))
	s.Info(s.requests)

	s.requests = []int64{1, 2, 3, 4, 5, 6, 7, 8, 9}
	s.discardOlder(-1, 3)
	assert.Equal(t, true, reflect.DeepEqual(s.requests, []int64{1, 2, 3, 4, 5, 6, 7, 8, 9}))
	s.Info(s.requests)

	s.requests = []int64{1, 2, 3, 4, 5, 6, 7, 8, 9}
	s.discardOlder(10, 0)
	assert.Equal(t, true, reflect.DeepEqual(s.requests, []int64{}))
	s.Info(s.requests)
}
