package main

import "fmt"

// Stack represents a stack data structure
type Stack[T any] struct {
	elements []T
}

// Push adds an element to the top of the stack
func (s *Stack[T]) Push(element T) {
	s.elements = append(s.elements, element)
}

// Pop removes and returns the top element of the stack
func (s *Stack[T]) Pop() (T, error) {
	if len(s.elements) == 0 {
    var zero T
		return zero, fmt.Errorf("stack is empty")
	}
	element := s.elements[len(s.elements)-1]
	s.elements = s.elements[:len(s.elements)-1]
	return element, nil
}

// Peek returns the top element of the stack without removing it
func (s *Stack[T]) Peek() (T, error) {
	if len(s.elements) == 0 {
    var zero T
		return zero, fmt.Errorf("stack is empty")
	}
	return s.elements[len(s.elements)-1], nil
}

func (s *Stack[T]) Get(index int) (T, error) {
  if len(s.elements) - 1 > index {
    var zero T
    return zero, fmt.Errorf("Index greater than stack length")
  }
  return s.elements[index], nil
}


func (s *Stack[T]) Size() int {
	return len(s.elements)
}

// IsEmpty checks if the stack is empty
func (s *Stack[T]) IsEmpty() bool {
	return len(s.elements) == 0
}
