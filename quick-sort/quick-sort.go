package main

import (
	"fmt"
	"math/rand"
)

const (
	SIZE = 10
)

func swap(m, n *int) {
	temp := *m
	*m = *n
	*n = temp
}

func partition(a []int, left, right, pivot int) int {
	for left < right {
		for a[left] < a[pivot] { // Move right
			left++
		}
		for a[right] > a[pivot] { // Move left
			right--
		}
		if left < right {
			swap(&a[left], &a[right])
		}
	}
	// Meet new pivot
	swap(&a[left], &a[pivot])
	return left
}

func quickSort(a []int, left, right int) {
	if left >= right {
		return
	}

	pivot := partition(a, left, right-1, right)
	quickSort(a, left, pivot-1)  // QuickSort left patition
	quickSort(a, pivot+1, right) // QuickSort right patition
}

func main() {
	a := []int{}
	for i := 0; i < SIZE; i++ {
		a = append(a, rand.Intn(20))
	}

	fmt.Printf("Before sort: %v", a)

	left := 0
	right := len(a) - 1
	quickSort(a, left, right)
	fmt.Printf("\nAfter sort: %v", a)
}
