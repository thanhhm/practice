package main

import (
	"fmt"
	"math/rand"
)

const (
	SIZE = 10
)

func sort(a []int) {
	size := len(a)
	for i := size - 1; i > 0; i-- {
		for j := 0; j < i; j++ {
			if a[j] > a[j+1] { // swap ASC
				swap(&a[j], &a[j+1])
			}
		}
	}
}

func swap(m, n *int) {
	temp := *m
	*m = *n
	*n = temp
}

func main() {
	a := []int{}
	for i := 0; i < SIZE; i++ {
		a = append(a, rand.Intn(20))
	}

	fmt.Printf("before sort: %v", a)
	sort(a)
	fmt.Printf("\nafter sort: %v", a)
}
