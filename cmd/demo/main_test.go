package main

import (
	"bufio"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

func ReadFile(filename string, consumeInt func(x int)) {

	f, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file ", err)
		os.Exit(1)
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		path, err := r.ReadString(10) // 0x0A separator = newline
		if err == io.EOF {
			// do something here
			break
		} else if err != nil {
			return
		}

		path = strings.TrimSuffix(path, "\n")
		path = strings.TrimSuffix(path, "\r\r\r")

		i, err := strconv.Atoi(path)
		if err != nil {
			fmt.Println(err)
			return
		}
		consumeInt(i)
	}
}

func Test1Odd(t *testing.T) {
	w := newWindow(3)
	testData := []int{100, 102, 101, 110, 120, 115}
	testAnswers := []int{-1, 101, 101, 102, 110, 115}

	m := -100
	for i, v := range testData {
		w.addDelay(v)
		m = w.getMedian()
		assert.Equal(t, testAnswers[i], m, "Wrong median value")
	}
	assert.Equal(t, 115, m, "Wrong median value")
}

func Test1Even(t *testing.T) {
	w := newWindow(2)
	testData := []int{100, 102, 101, 110, 120, 115}

	for _, v := range testData {
		w.addDelay(v)
	}
	m := w.getMedian()
	assert.Equal(t, 117, m, "Wrong median value")
}

func TestFromFile0(t *testing.T) {
	w := newWindow(6)
	m := -100
	ReadFile("../../data/test0.csv", func(x int) {
		w.addDelay(x)
		m = w.getMedian()
	})
	assert.Equal(t, 106, m, "Wrong median value")
}

func TestFromFile2(t *testing.T) {
	w := newWindow(1000)
	m := -100
	ReadFile("../../data/test2.csv", func(x int) {
		w.addDelay(x)
		m = w.getMedian()
	})
	assert.Equal(t, 289, m, "Wrong median value")
}

func TestFromFile3(t *testing.T) {
	w := newWindow(10000)
	m := -100
	ReadFile("../../data/test3.csv", func(x int) {
		w.addDelay(x)
		m = w.getMedian()
	})
	assert.Equal(t, 298, m, "Wrong median value")
}

func TestFromFile4(t *testing.T) {
	w := newWindow(100000)
	m := -100
	ReadFile("../../data/test4.csv", func(x int) {
		w.addDelay(x)
		m = w.getMedian()
	})
	assert.Equal(t, 301, m, "Wrong median value")
}
