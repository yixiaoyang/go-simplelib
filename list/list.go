package list

type Node struct {
	pre, nxt *Node
	value    any
}

type List struct {
	head *Node
	tail *Node
	size int
}
