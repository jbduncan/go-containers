package set

import (
	"iter"
	"maps"
)

// Set is a generic, unordered collection of unique elements.
//
// An instance of Set can be made with set.Of.
type Set[T comparable] interface {
	// Contains returns true if this set contains the given element, otherwise it returns false.
	Contains(elem T) bool

	// Len returns the number of elements in this set.
	Len() int

	// All returns an iter.Seq that returns each and every element in this set.
	//
	// The iteration order is undefined; it may even change from one call to the next.
	All() iter.Seq[T]

	// String returns a string representation of all the elements in this set.
	//
	// The format of this string is a single "[" followed by a comma-separated list (", ") of this set's elements in
	// the same order as All (which is undefined and may change from one call to the next), followed by a single
	// "]".
	//
	// This method satisfies fmt.Stringer.
	String() string
}

// MutableSet is a Set with additional methods for adding and removing elements.
//
// An instance of MutableSet can be made with set.Of.
type MutableSet[T comparable] interface {
	Set[T]

	// Add adds the given element(s) to this set. If any of the elements are already present, the set will not add
	// those elements again. Returns true if this set changed as a result of this call, otherwise false.
	Add(elem T, others ...T) bool

	// Remove removes the given element(s) from this set. If any of the elements are already absent, the set will not
	// attempt to remove those elements. Returns true if this set changed as a result of this call, otherwise false.
	Remove(elem T, others ...T) bool
}

var (
	_ Set[int]        = (*MapSet[int])(nil)
	_ MutableSet[int] = (*MapSet[int])(nil)
)

// Of returns a new non-nil, empty MapSet, which implements MutableSet. Its implementation is based on a Go map, with
// similar performance characteristics.
func Of[T comparable](elems ...T) *MapSet[T] {
	delegate := make(map[T]struct{}, len(elems))
	for _, elem := range elems {
		delegate[elem] = struct{}{}
	}
	return &MapSet[T]{
		delegate: delegate,
	}
}

type MapSet[T comparable] struct {
	delegate map[T]struct{}
}

func (m *MapSet[T]) Add(elem T, others ...T) bool {
	result := m.add(elem)
	for _, other := range others {
		added := m.add(other)
		result = result || added
	}
	return result
}

func (m *MapSet[T]) add(elem T) bool {
	_, ok := m.delegate[elem]
	m.delegate[elem] = struct{}{}
	return !ok
}

func (m *MapSet[T]) Remove(elem T, others ...T) bool {
	result := m.remove(elem)
	for _, other := range others {
		removed := m.remove(other)
		result = result || removed
	}
	return result
}

func (m *MapSet[T]) remove(elem T) bool {
	_, ok := m.delegate[elem]
	delete(m.delegate, elem)
	return ok
}

func (m *MapSet[T]) Contains(elem T) bool {
	_, ok := m.delegate[elem]
	return ok
}

func (m *MapSet[T]) Len() int {
	return len(m.delegate)
}

func (m *MapSet[T]) All() iter.Seq[T] {
	return maps.Keys(m.delegate)
}

func (m *MapSet[T]) String() string {
	return StringImpl[T](m)
}
