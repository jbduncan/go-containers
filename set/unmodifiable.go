package set

import "iter"

var _ Set[int] = (*UnmodifiableSet[int])(nil)

// Unmodifiable wraps a given set as a read-only Set view. This prevents users from casting the given set into a
// MutableSet and allows for sets that can be mutated by your own code but not your users' code.
//
// If the given set is ever mutated, then the returned set will reflect those mutations.
func Unmodifiable[T comparable](s Set[T]) UnmodifiableSet[T] {
	return UnmodifiableSet[T]{
		s: s,
	}
}

type UnmodifiableSet[T comparable] struct {
	s Set[T]
}

func (u UnmodifiableSet[T]) Contains(elem T) bool {
	return u.s.Contains(elem)
}

func (u UnmodifiableSet[T]) Len() int {
	return u.s.Len()
}

func (u UnmodifiableSet[T]) All() iter.Seq[T] {
	return u.s.All()
}

func (u UnmodifiableSet[T]) String() string {
	return u.s.String()
}
