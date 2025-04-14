package set

import "iter"

// Unmodifiable wraps a given set as a read-only set view. This prevents users
// from casting the given set into a Set or another mutable set
// implementation. Also, it allows libraries to mutate the original set without
// the library's users mutating it, too.
//
// If the given set is ever mutated, then the returned set will reflect those
// mutations.
//
// Note: Go needs the generic type to be defined explicitly, like:
//
//	s := set.Of(1)
//	s := set.Unmodifiable[int](a)
//	                     ^^^^^
func Unmodifiable[T comparable](s interface {
	Contains(elem T) bool
	Len() int
	All() iter.Seq[T]
	String() string
},
) UnmodifiableSet[T] {
	return UnmodifiableSet[T]{
		s: s,
	}
}

type UnmodifiableSet[T comparable] struct {
	s interface {
		Contains(elem T) bool
		Len() int
		All() iter.Seq[T]
		String() string
	}
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
