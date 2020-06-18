package main

import "fmt"

type Node struct {
	v    int
	next *Node
}

func (n *Node) String() string {
	if n == nil {
		return "nil"
	}
	return fmt.Sprintf("%d ", n.v)
}

type List struct {
	head *Node
}

func (l *List) push(v int) {
	a := &Node{v: v}
	a.next = l.head
	l.head = a
}

func (l *Node) count() (count int) {
	for l != nil {
		l = l.next
		count++
	}
	return
}

func (l *Node) reverse() *Node {
	c := l
	prev := (*Node)(nil)

	for {
		next := c.next
		c.next = prev
		prev = c
		if next == nil {
			break
		}
		c = next
	}
	return prev
}

func (l *Node) reverseN(n int) *Node {
	c := l
	prev := (*Node)(nil)
	end := c

	for i := n; i > 0; i-- {
		next := c.next
		c.next = prev

		prev = c
		c = next
		if next == nil {
			break
		}
	}

	if c != nil {
		if c.count() >= n {
			end.next = c.reverseN(2)
		} else {
			end.next = c
		}
	}
	return prev
}

func (l *Node) print() {
	for l != nil {
		fmt.Print(l.v, " -> ")
		l = l.next
	}
	fmt.Println()
}
