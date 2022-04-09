package set

import (
	"fmt"
	"strings"
)

type Set[T comparable] interface {
	Contains(value T) bool
	Len() int
	ForEach(fn func(elem T))
	String() string
	// TODO: Set: make Iter method that returns an Iterator
	// TODO: Set: make Equals method and discourage == from being used (documenting that its use is undefined).
}

type MutableSet[T comparable] interface {
	Set[T]

	Add(value T)
	Remove(value T)
}

func New[T comparable]() MutableSet[T] {
	return &set[T]{
		delegate: map[T]struct{}{},
	}
}

type set[T comparable] struct {
	delegate map[T]struct{}
}

var _ Set[int] = (*set[int])(nil)
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

func Unmodifiable[T comparable](set MutableSet[T]) Set[T] {
	return unmodifiableSet[T]{
		set: set,
	}
}

func (u unmodifiableSet[T]) Contains(value T) bool {
	return u.set.Contains(value)
}

func (u unmodifiableSet[T]) Len() int {
	return u.set.Len()
}

func (u unmodifiableSet[T]) ForEach(fn func(elem T)) {
	u.set.ForEach(fn)
}

func (u unmodifiableSet[T]) String() string {
	return u.set.String()
}
