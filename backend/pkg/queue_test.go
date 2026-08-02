package pkg

import "testing"

func TestQueueEnqueueDequeue(t *testing.T) {
	q := NewQueue[int]()

	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	if q.Len() != 3 {
		t.Fatalf("Len() = %d, want 3", q.Len())
	}

	got, ok := q.Dequeue()
	if !ok || got != 1 {
		t.Fatalf("Dequeue() = (%d, %v), want (1, true)", got, ok)
	}

	got, ok = q.Dequeue()
	if !ok || got != 2 {
		t.Fatalf("Dequeue() = (%d, %v), want (2, true)", got, ok)
	}

	got, ok = q.Dequeue()
	if !ok || got != 3 {
		t.Fatalf("Dequeue() = (%d, %v), want (3, true)", got, ok)
	}

	if !q.IsEmpty() {
		t.Fatal("expected empty queue after dequeuing all items")
	}
}

func TestQueueDequeueEmpty(t *testing.T) {
	q := NewQueue[string]()

	got, ok := q.Dequeue()
	if ok || got != "" {
		t.Fatalf("Dequeue() on empty = (%q, %v), want (\"\", false)", got, ok)
	}
}

func TestQueuePeek(t *testing.T) {
	q := NewQueue[int]()

	if _, ok := q.Peek(); ok {
		t.Fatal("Peek() on empty should return false")
	}

	q.Enqueue(10)
	q.Enqueue(20)

	got, ok := q.Peek()
	if !ok || got != 10 {
		t.Fatalf("Peek() = (%d, %v), want (10, true)", got, ok)
	}
	if q.Len() != 2 {
		t.Fatalf("Peek should not remove items; Len() = %d, want 2", q.Len())
	}
}

func TestQueueIsEmpty(t *testing.T) {
	q := NewQueue[int]()
	if !q.IsEmpty() {
		t.Fatal("new queue should be empty")
	}

	q.Enqueue(1)
	if q.IsEmpty() {
		t.Fatal("queue with items should not be empty")
	}
}
