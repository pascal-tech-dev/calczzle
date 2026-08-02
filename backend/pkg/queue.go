package pkg

// Queue is a FIFO container backed by a slice.
type Queue[T any] struct {
	items []T
}

// NewQueue creates an empty queue.
func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{items: make([]T, 0)}
}

// Enqueue adds an item to the back of the queue.
func (q *Queue[T]) Enqueue(item T) {
	q.items = append(q.items, item)
}

// Dequeue removes and returns the front item.
// The boolean is false when the queue is empty.
func (q *Queue[T]) Dequeue() (T, bool) {
	var zero T
	if len(q.items) == 0 {
		return zero, false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

// Peek returns the front item without removing it.
// The boolean is false when the queue is empty.
func (q *Queue[T]) Peek() (T, bool) {
	var zero T
	if len(q.items) == 0 {
		return zero, false
	}
	return q.items[0], true
}

// Len returns the number of items in the queue.
func (q *Queue[T]) Len() int {
	return len(q.items)
}

// IsEmpty reports whether the queue has no items.
func (q *Queue[T]) IsEmpty() bool {
	return len(q.items) == 0
}
