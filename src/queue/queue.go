package queue

import "fmt"

// Queue is a basic FIFO queue based on a circular list that resizes as needed.
type Queue struct {
	nodes []int
	size  int
	head  int
	tail  int
	count int
}

// New returns a new queue with the given initial size.
func New(size int) *Queue {
	return &Queue{
		nodes: make([]int, size),
		size:  size,
	}
}

// Push adds a node to the queue.
func (q *Queue) Push(n int) {
	if q.head == q.tail && q.count > 0 {
		nodes := make([]int, len(q.nodes)+q.size)
		copy(nodes, q.nodes[q.head:])
		copy(nodes[len(q.nodes)-q.head:], q.nodes[:q.head])
		q.head = 0
		q.tail = len(q.nodes)
		q.nodes = nodes
	}
	q.nodes[q.tail] = n
	q.tail = (q.tail + 1) % len(q.nodes)
	q.count++
}

// Pop removes and returns a node from the queue in first to last order.
func (q *Queue) Pop() (int, error) {
	if q.count == 0 {
		return 0, fmt.Errorf("Queue empty")
	}
	node := q.nodes[q.head]
	q.head = (q.head + 1) % len(q.nodes)
	q.count--
	return node, nil
}

// Empty returns true if there are no elements in the queue, fale otherwise
func (q *Queue) Empty() bool {
	if q.count == 0 {
		return true
	}
	return false
}
