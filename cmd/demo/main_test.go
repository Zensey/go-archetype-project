package main

import (
	"bytes"
	"testing"

	"github.com/Zensey/go-archetype-project/pkg/domain"
)

func Test_Main(t *testing.T) {

	p := domain.Producer{}
	b := bytes.Buffer{}

	for i := 0; i < 1000000; i++ {
		b.Reset()
		p.GetNewMsgID(&b)
	}
}
