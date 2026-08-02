package pkg

// Stack is a LIFO container backed by a slice.
type Stack[T any] struct {
	items []T
}

// NewStack creates an empty stack.
func NewStack[T any]() *Stack[T] {
	return &Stack[T]{items: make([]T, 0)}
}

// Push adds an item to the top of the stack.
func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

// Pop removes and returns the top item.
// The boolean is false when the stack is empty.
func (s *Stack[T]) Pop() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	i := len(s.items) - 1
	item := s.items[i]
	s.items = s.items[:i]
	return item, true
}

// Peek returns the top item without removing it.
// The boolean is false when the stack is empty.
func (s *Stack[T]) Peek() (T, bool) {
	var zero T
	if len(s.items) == 0 {
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

// Len returns the number of items on the stack.
func (s *Stack[T]) Len() int {
	return len(s.items)
}

// IsEmpty reports whether the stack has no items.
func (s *Stack[T]) IsEmpty() bool {
	return len(s.items) == 0
}
