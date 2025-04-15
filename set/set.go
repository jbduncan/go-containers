package set

import (
	"iter"
	"maps"
)

// Of returns a new non-nil, empty Set, which is a generic, unordered
// collection of unique elements. Its implementation is based on a Go map, with
// similar performance characteristics.
func Of[T comparable](elements ...T) Set[T] {
	delegate := make(map[T]struct{}, len(elements))
	for _, elem := range elements {
		delegate[elem] = struct{}{}
	}
	return Set[T]{
		delegate: delegate,
	}
}

// Set is a generic, unordered collection of unique elements. Its
// implementation is based on a Go map, with similar performance
// characteristics.
type Set[T comparable] struct {
	delegate map[T]struct{}
}

// Contains returns true if this set contains the given element, otherwise it
// returns false.
func (m Set[T]) Contains(elem T) bool {
	_, ok := m.delegate[elem]
	return ok
}

// Len returns the number of elements in this set.
func (m Set[T]) Len() int {
	return len(m.delegate)
}

// All returns an iter.Seq that returns each and every element in this set.
//
// The iteration order is undefined; it may even change from one call to the
// next.
func (m Set[T]) All() iter.Seq[T] {
	return maps.Keys(m.delegate)
}

// String returns a string representation of all the elements in this set.
//
// The format of this string is a single "[" followed by a comma-separated list
// (", ") of this set's elements in the same order as All (which is undefined
// and may change from one call to the next), followed by a single "]".
//
// This method satisfies fmt.Stringer.
func (m Set[T]) String() string {
	return StringImpl[T](m)
}

// Add adds the given element(s) to this set. If any of the elements are
// already present, the set will not add those elements again. Returns true if
// this set changed as a result of this call, otherwise false.
func (m Set[T]) Add(elem T, others ...T) bool {
	result := m.addInternal(elem)
	for _, other := range others {
		added := m.addInternal(other)
		result = result || added
	}
	return result
}

func (m Set[T]) addInternal(elem T) bool {
	_, ok := m.delegate[elem]
	m.delegate[elem] = struct{}{}
	return !ok
}

// Remove removes the given element(s) from this set. If any of the elements
// are already absent, the set will not attempt to remove those elements.
// Returns true if this set changed as a result of this call, otherwise false.
func (m Set[T]) Remove(elem T, others ...T) bool {
	result := m.removeInternal(elem)
	for _, other := range others {
		removed := m.removeInternal(other)
		result = result || removed
	}
	return result
}

func (m Set[T]) removeInternal(elem T) bool {
	_, ok := m.delegate[elem]
	delete(m.delegate, elem)
	return ok
}
