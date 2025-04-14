package set

import "iter"

// Union returns the set union of sets a and b.
//
// The returned set is a read-only view that implements Set, so changes to a
// and b will be reflected in the returned set.
//
// Set.Len runs in O(b) time for the returned set.
//
// Note: Go needs the generic type to be defined explicitly, like:
//
//	a := set.Of(1)
//	b := set.Of(2)
//	u := set.Union[int](a, b)
//	              ^^^^^
func Union[T comparable](a, b interface {
	Contains(elem T) bool
	All() iter.Seq[T]
	Len() int
},
) UnionSet[T] {
	return UnionSet[T]{
		a: a,
		b: b,
	}
}

type UnionSet[T comparable] struct {
	a, b interface {
		Contains(elem T) bool
		All() iter.Seq[T]
		Len() int
	}
}

func (u UnionSet[T]) Contains(elem T) bool {
	return u.a.Contains(elem) || u.b.Contains(elem)
}

func (u UnionSet[T]) Len() int {
	bLen := 0
	for elem := range u.b.All() {
		if !u.a.Contains(elem) {
			bLen++
		}
	}
	return u.a.Len() + bLen
}

func (u UnionSet[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for elem := range u.a.All() {
			if !yield(elem) {
				return
			}
		}

		for elem := range u.b.All() {
			if u.a.Contains(elem) {
				continue
			}
			if !yield(elem) {
				return
			}
		}
	}
}

func (u UnionSet[T]) String() string {
	return StringImpl[T](u)
}
