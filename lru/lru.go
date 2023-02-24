package lru

type Node[k comparable, v any] struct {
	pre, nxt *Node[k, v]
	value    v
	key      any
}

type List[k comparable, v any] struct {
	head *Node[k, v]
	tail *Node[k, v]
	size int
}

type IList[k comparable, v any] interface {
	Prepend(key k, value v) *Node[k, v]
	Append(key k, value v) *Node[k, v]
	MoveToFront(*Node[k, v])
	Len() int
}

func NewList[k comparable, v any]() *List[k, v] {
	node_head := &Node[k, v]{
		pre: nil,
		nxt: nil,
	}
	node_tail := &Node[k, v]{
		pre: node_head,
		nxt: nil,
	}
	node_head.nxt = node_tail
	return &List[k, v]{
		head: node_head,
		tail: node_tail,
	}
}

func (list *List[k, v]) Prepend(key k, value v) *Node[k, v] {
	node := &Node[k, v]{
		pre:   list.head,
		nxt:   list.head.nxt,
		key:   key,
		value: value,
	}
	list.head.nxt = node
	list.size += 1
	return list.head.nxt
}

func (list *List[k, v]) Append(key k, value v) *Node[k, v] {
	node := &Node[k, v]{
		pre:   list.tail.pre,
		nxt:   list.tail,
		key:   key,
		value: value,
	}
	list.tail.pre.nxt = node
	list.tail.pre = node
	list.size += 1
	return node
}

func (list *List[k, v]) MoveToFront(node *Node[k, v]) {
	if list.head.nxt == node {
		return
	}
	node.pre.nxt = node.nxt
	node.nxt.pre = node.pre

	node.pre = list.head
	node.nxt = list.head.nxt
	list.head.nxt.pre = node
	list.head.nxt = node
}

func (list *List[k, v]) Remove(node *Node[k, v]) {
	node.pre.nxt = node.nxt
	node.nxt.pre = node.pre
	list.size -= 1
}

func (list *List[k, v]) Len() int {
	return list.size
}

// Lru implements a non-thread-safe lib of lru cache
type Lru[k comparable, v any] struct {
	list *List[k, v]
	hash map[k]*Node[k, v]
}

type ILru[k comparable, v any] interface {
	Add(key k, value v) (overwrite bool)
	Get(key k) (value v, exist bool)
	Remove(key k) (exist bool)
	RemoveOldest() (key k, value v)
	Clear()
	Len() int
}

func NewLru[k comparable, v any]() *Lru[k, v] {
	return &Lru[k, v]{
		list: NewList[k, v](),
		hash: make(map[k]*Node[k, v]),
	}
}

func (lru *Lru[k, v]) Add(key k, value v) {
	if node, ok := lru.hash[key]; ok {
		lru.list.MoveToFront(node)
	} else {
		lru.hash[key] = lru.list.Prepend(key, value)
	}
}

func (lru *Lru[k, v]) Get(key k) (value v, exist bool) {
	if node, ok := lru.hash[key]; ok {
		lru.list.MoveToFront(node)
		return node.value, true
	}
	var temp v
	return temp, false
}

func (lru *Lru[k, v]) Remove(key k) (exist bool) {
	if node, ok := lru.hash[key]; ok {
		lru.list.Remove(node)
		delete(lru.hash, key)
		return true
	}
	return false
}

func (lru *Lru[k, v]) Clear() {
	// gc will recycle it
	lru.hash = make(map[k]*Node[k, v])
	lru.list = NewList[k, v]()
}

func (lru *Lru[k, v]) Len() int {
	return lru.list.Len()
}
