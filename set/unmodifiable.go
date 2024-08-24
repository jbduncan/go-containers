package set

import "iter"

var _ Set[int] = (*UnmodifiableSet[int])(nil)

// Unmodifiable wraps a given MutableSet as a read-only Set view. This allows for sets that can be mutated by your own
// code but not your clients' code.
//
// If the original MutableSet is ever mutated, then the returned Set will reflect those mutations.
//
// This set cannot be cast back into a MutableSet.
func Unmodifiable[T comparable](set MutableSet[T]) UnmodifiableSet[T] {
	return UnmodifiableSet[T]{
		set: set,
	}
}

type UnmodifiableSet[T comparable] struct {
	set MutableSet[T]
}

func (u UnmodifiableSet[T]) Contains(elem T) bool {
	return u.set.Contains(elem)
}

func (u UnmodifiableSet[T]) Len() int {
	return u.set.Len()
}

func (u UnmodifiableSet[T]) All() iter.Seq[T] {
	return u.set.All()
}

func (u UnmodifiableSet[T]) String() string {
	return u.set.String()
}
