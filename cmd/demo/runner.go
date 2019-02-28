package main

import (
	"github.com/ancientlore/go-avltree"
)

type MedianRunner struct {
	tree   *avltree.Tree
	n      int
	window chan int
}

func compareInt(a interface{}, b interface{}) int {
	if a.(int) < b.(int) {
		return -1
	} else if a.(int) > b.(int) {
		return 1
	}
	return 0
}

func newWindow(n int) *MedianRunner {
	runner := &MedianRunner{
		tree:   avltree.New(compareInt, avltree.AllowDuplicates),
		n:      0,
		window: make(chan int, n),
	}
	return runner
}

func (w *MedianRunner) addDelay(x int) {
	if w.n == cap(w.window) {
		key := <-w.window
		w.tree.Remove(key)
	} else {
		w.n++
	}
	w.window <- x
	w.tree.Add(x)
}

func (w *MedianRunner) getMedian() int {
	len := w.tree.Len()
	if len <= 1 {
		return -1
	}
	if len%2 == 1 {
		return w.tree.At(len / 2).(int)
	} else {
		return (w.tree.At(len/2-1).(int) + w.tree.At(len/2).(int)) / 2
	}
}
