package main

import (
	"fmt"
	"math/rand"
)

const (
	SIZE = 10
)

func mergeSort(a []int) []int {
	length := len(a)
	if length == 1 {
		return a
	}

	a1 := mergeSort(a[0 : length/2])
	a2 := mergeSort(a[length/2:])

	return merge(a1, a2)
}

func merge(a []int, b []int) []int {
	var c []int
	c = make([]int, 0)

	for len(a) > 0 && len(b) > 0 {
		if a[0] > b[0] { // Add lower number to c[]
			c = append(c, b[0])
			b = append(b[:0], b[1:]...) // Remove b[0]
		} else {
			c = append(c, a[0])
			a = append(a[:0], a[1:]...) // Remove a[0]
		}
	}

	if len(a) > 0 {
		c = append(c, a[0:]...)
		a = a[:0]
	}
	if len(b) > 0 {
		c = append(c, b[0:]...)
		b = b[:0]
	}

	return c
}

func main() {
	a := []int{}
	for i := 0; i < SIZE; i++ {
		a = append(a, rand.Intn(20))
	}

	fmt.Printf("before sort: %v", a)
	fmt.Printf("\nafter sort: %v", mergeSort(a))
}
