package main

type Window struct {
	arr    []int
	wsz    int
	pos    int
	isFull bool
}

func newWindow(n int) Window {
	return Window{
		arr: make([]int, n),
		wsz: n,
		pos: 0,
	}
}

func (w *Window) addDelay(x int) {
	w.arr[w.pos] = x
	w.pos++
	if w.pos == w.wsz {
		w.isFull = true
		w.pos = 0
	}
}

func (w *Window) getMedian() int {
	len := len(w.arr)
	if !w.isFull {
		len = w.pos
	}
	slice := make([]int, len)
	copy(slice, w.arr)

	if len <= 1 {
		return -1
	}
	if len%2 == 1 {
		return quickselect_(slice, len/2)
	} else {
		return (quickselect_(slice, len/2-1) + quickselect_(slice, len/2)) / 2
	}
}
