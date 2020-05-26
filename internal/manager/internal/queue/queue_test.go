package queue

import (
	"testing"
	"time"
)

func TestQueue_Dequeue(t *testing.T) {

	q := New()

	go func() {
		_ = q.Dequeue()
		t.Fatalf("dequeue must not return anything")
	}()

	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	i1 := q.Dequeue()
	i2 := q.Dequeue()

	if val, ok := i1.(int); ok {
		if val != 1 {
			t.Fatalf("the value returned is not value")
		}
	}

	if val, ok := i2.(int); ok {
		if val != 2 {
			t.Fatalf("the value returned is not value")
		}
	}

	time.Sleep(time.Second)

}
