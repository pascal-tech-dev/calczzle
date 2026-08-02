package pkg

import "testing"

func TestStackPushPop(t *testing.T) {
	s := NewStack[int]()

	s.Push(1)
	s.Push(2)
	s.Push(3)

	if s.Len() != 3 {
		t.Fatalf("Len() = %d, want 3", s.Len())
	}

	got, ok := s.Pop()
	if !ok || got != 3 {
		t.Fatalf("Pop() = (%d, %v), want (3, true)", got, ok)
	}

	got, ok = s.Pop()
	if !ok || got != 2 {
		t.Fatalf("Pop() = (%d, %v), want (2, true)", got, ok)
	}

	got, ok = s.Pop()
	if !ok || got != 1 {
		t.Fatalf("Pop() = (%d, %v), want (1, true)", got, ok)
	}

	if !s.IsEmpty() {
		t.Fatal("expected empty stack after popping all items")
	}
}

func TestStackPopEmpty(t *testing.T) {
	s := NewStack[string]()

	got, ok := s.Pop()
	if ok || got != "" {
		t.Fatalf("Pop() on empty = (%q, %v), want (\"\", false)", got, ok)
	}
}

func TestStackPeek(t *testing.T) {
	s := NewStack[int]()

	if _, ok := s.Peek(); ok {
		t.Fatal("Peek() on empty should return false")
	}

	s.Push(10)
	s.Push(20)

	got, ok := s.Peek()
	if !ok || got != 20 {
		t.Fatalf("Peek() = (%d, %v), want (20, true)", got, ok)
	}
	if s.Len() != 2 {
		t.Fatalf("Peek should not remove items; Len() = %d, want 2", s.Len())
	}
}

func TestStackIsEmpty(t *testing.T) {
	s := NewStack[int]()
	if !s.IsEmpty() {
		t.Fatal("new stack should be empty")
	}

	s.Push(1)
	if s.IsEmpty() {
		t.Fatal("stack with items should not be empty")
	}
}
