package queue

import (
	"sync"
)

// Queue is a blocking queue implementation
// once there are no more items dequeue will block
type Queue interface {
	Enqueue(i interface{})
	Dequeue() interface{}
}

func New() Queue {
	q := &queue{
		store: make([]interface{}, 0),
		l:     sync.Mutex{},
		l2:    sync.Mutex{},
	}

	q.l2.Lock()

	return q
}

type queue struct {
	store []interface{}
	l     sync.Mutex
	l2    sync.Mutex
}

func (q *queue) Enqueue(i interface{}) {
	q.l.Lock()
	defer q.l.Unlock()

	// check if l2 is locked
	qLength := len(q.store)

	q.store = append(q.store, i)

	if qLength == 0 {
		q.l2.Unlock()
	}
}

func (q *queue) Dequeue() interface{} {

	q.l2.Lock()

	q.l.Lock()
	defer q.l.Unlock()

	i := q.store[0]
	q.store = q.store[1:]

	l := len(q.store)

	if l > 0 {
		q.l2.Unlock()
	}

	return i
}
