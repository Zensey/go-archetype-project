package main

func main() {
	l := List{}
	l.push(5)
	l.push(4)
	l.push(3)
	l.push(2)
	l.push(1)

	l.head.print()
	new := l.head.reverseN(2)
	new.print()

	return
}
