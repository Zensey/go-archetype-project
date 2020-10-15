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

// Find the smallest positive number missing from an unsorted array

func task1() {
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

/* task2:
	You are given a list of n-1 integers and these integers are in the range of 1 to n. There are no duplicates in list. One of the integers is missing in the list.

	I/P    [1, 2, 4, ,6, 3, 7, 8]
	O/P    5

       S=(1+n)*n/2
       S_act
       missing = S_calc - S_act
*/

func main() {
	task1()
	return
}
