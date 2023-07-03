package main

import "fmt"

const Size = 5

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
	return Cache{Queue: NewQueue(), Hash: Hash{}}
}

func NewQueue() Queue {
	head := &Node{}
	tail := &Node{}

	head.right = tail
	tail.left = head

	return Queue{Head: head, Tail: tail}
}

func (c *Cache) Check(str string) {
	node := &Node{}

	if val, ok := c.Hash[str]; ok {
		node = c.Remove(val)
	} else {
		node = &Node{val: str}
	}
	c.Add(node)
	c.Hash[str] = node
}

func (c *Cache) Remove(n *Node) *Node {
	fmt.Println("Remove: ", n.val)
	left := n.left
	right := n.right

	right.left = left
	left.right = right
	c.Queue.Length -= 1
	delete(c.Hash, n.val)
	return n
}

func (c *Cache) Add(n *Node) {
	fmt.Println("Add: ", n.val)
	temp := c.Queue.Head.right
	c.Queue.Head.right = n
	n.left = c.Queue.Head
	n.right = temp
	temp.left = n

	c.Queue.Length++
	if c.Queue.Length > Size {
		c.Remove(c.Queue.Tail.left)
	}
}

func (c *Cache) Display() {
	c.Queue.Display()
}

func (q *Queue) Display() {
	node := q.Head.right
	fmt.Printf("%d - [", q.Length)
	for i := 0; i < q.Length; i++ {
		fmt.Printf("{%s}", node.val)
		if i<q.Length-1 {
			fmt.Print("<-->")
		}
		node = node.right
	}
	fmt.Println("]")
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
