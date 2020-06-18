package main

import "fmt"

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func findFirstPositive(a []int) (i int, v int) {
	for i := 0; i < len(a); i++ {
		if a[i] > 0 {
			return i, a[i]
		}
	}
	return -1, 0
}

func main() {
	a := []int{2, 3, -7, 6, 8, 1, -10, 15}
	fmt.Println(a)

	_, firstPos := findFirstPositive(a)
	for i, _ := range a {
		if a[i] <= 0 {
			a[i] = firstPos
		}
	}
	for i, _ := range a {
		j := abs(a[i]) - 1
		if j < len(a) && a[j] > 0 {
			a[j] = -a[j]
		}
	}

	firstPosI, _ := findFirstPositive(a)
	fmt.Println(a, firstPosI+1)
	return
}
