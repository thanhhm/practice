package main

import (
	"fmt"
)

type ListNode struct {
	Value interface{}
	Next  *ListNode
}

func (l *ListNode) pushBack(v interface{}) {
	if l != nil {
		// Go to the last node
		for l.Next != nil {
			l = l.Next
		}

		// Push new node with value v
		l.Next = &ListNode{
			Value: v,
		}
	} else {
		l = &ListNode{
			Value: v,
		}
	}

}

func (l *ListNode) pushFront(v interface{}) *ListNode {
	return &ListNode{
		Value: v,
		Next:  l,
	}
}

func (l *ListNode) printListNode() {
	fmt.Printf("[")
	for l != nil {
		fmt.Printf("% v ", l.Value)
		l = l.Next
	}
	fmt.Printf("]")
}

func main() {
	l := &ListNode{
		Value: 1,
	}

	l.pushBack(2)
	l.pushBack(3)
	l.printListNode()
	fmt.Printf("\n%v", l)

}
