package main

import "fmt"

type Node struct {
	val   string
	left  *Node
	right *Node
}

type Queue struct {
	Head   *Node
	Tail   *Node
	Length int
}

type Cache struct {
	Queue Queue
	Hash  Hash
}

func NewCache() Cache {
	return Cache{Queue: NewQueue(), Hash: NewHash{}}
}

func NewQueue() Queue {
	head := &Node{}
	tail := &Node{}

	head.right = tail
	tail.left = head

	return Queue{Head: head, Tail: tail}
}

type Hash map[string]*Node

func main() {
	fmt.Println("Start cache...")
	cache := NewCache()
	for _, word := range []string{"pikachu", " bulbasaur", "charmander", "squirtle", "butterfree", "pigeot"} {
		cache.Check(word)
		cache.Display()
	}
}
