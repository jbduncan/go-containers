package container

import (
	"fmt"
	"strings"
)

type Set[T comparable] interface {
	Contains(value T) bool
	Len() int
	ForEach(fn func(elem T))
	// Iter() Iterator
	String() string
}

type MutableSet[T comparable] interface {
	Set[T]

	Add(value T)
	Remove(value T)
}

func NewSet[T comparable]() MutableSet[T] {
	return &set[T]{
		delegate: map[T]struct{}{},
	}
}

type set[T comparable] struct {
	delegate map[T]struct{}
}

var _ MutableSet[int] = (*set[int])(nil)

func (s *set[T]) Add(value T) {
	s.delegate[value] = struct{}{}
}

func (s *set[T]) Remove(value T) {
	delete(s.delegate, value)
}

func (s set[T]) Contains(value T) bool {
	_, ok := s.delegate[value]
	return ok
}

func (s set[T]) Len() int {
	return len(s.delegate)
}

func (s set[T]) ForEach(fn func(elem T)) {
	for elem := range s.delegate {
		fn(elem)
	}
}

func (s set[T]) String() string {
	var b strings.Builder
	b.WriteRune('[')
	i := 0
	for elem := range s.delegate {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(fmt.Sprintf("%v", elem))
		i++
	}
	b.WriteRune(']')
	return b.String()
}

type unmodifiableSet[T comparable] struct {
	set MutableSet[T]
}

var _ Set[int] = (*unmodifiableSet[int])(nil)

func (r unmodifiableSet[T]) Contains(value T) bool {
	// return r.set.Contains(v)
	return false
}

func (r unmodifiableSet[T]) Len() int {
	// return r.set.Len()
	return 0
}

func (r unmodifiableSet[T]) ForEach(fn func(elem T)) {
	// r.set.ForEach(f)
}

func (r unmodifiableSet[T]) String() string {
	// return r.set.String()
	return ""
}

func UnmodifiableSet[T comparable](set Set[T]) Set[T] {
	return nil
}
