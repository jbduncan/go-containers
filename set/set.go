package set

import "iter"

// Set is a generic, unordered collection of unique elements.
//
// An instance of Set can be made with set.Of. If a mutable set is needed, use set.NewMutable
// instead.
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
// An instance of MutableSet can be made with set.NewMutable.
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
	_ Set[int]        = (*mutableMapSet[int])(nil)
	_ MutableSet[int] = (*mutableMapSet[int])(nil)
)

// Of returns a new non-nil, empty, immutable Set. Its implementation is based on a Go map, with
// similar performance characteristics.
func Of[T comparable](elems ...T) Set[T] {
	delegate := make(map[T]struct{}, len(elems))
	for _, elem := range elems {
		delegate[elem] = struct{}{}
	}
	return &mapSet[T]{
		delegate: delegate,
	}
}

// NewMutable returns a new non-nil, empty MutableSet. Its implementation is based on a Go map,
// with similar performance characteristics.
func NewMutable[T comparable](elems ...T) MutableSet[T] {
	delegate := map[T]struct{}{}
	for _, elem := range elems {
		delegate[elem] = struct{}{}
	}
	return &mutableMapSet[T]{
		delegate: delegate,
	}
}

type mapSet[T comparable] struct {
	delegate map[T]struct{}
}

func (m *mapSet[T]) Contains(elem T) bool {
	_, ok := m.delegate[elem]
	return ok
}

func (m *mapSet[T]) Len() int {
	return len(m.delegate)
}

func (m *mapSet[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for elem := range m.delegate {
			if !yield(elem) {
				return
			}
		}
	}
}

func (m *mapSet[T]) String() string {
	return StringImpl[T](m)
}

type mutableMapSet[T comparable] struct {
	delegate map[T]struct{}
}

func (m *mutableMapSet[T]) Add(elem T, others ...T) bool {
	result := m.add(elem)
	for _, other := range others {
		added := m.add(other)
		result = result || added
	}
	return result
}

func (m *mutableMapSet[T]) add(elem T) bool {
	_, ok := m.delegate[elem]
	m.delegate[elem] = struct{}{}
	return !ok
}

func (m *mutableMapSet[T]) Remove(elem T, others ...T) bool {
	result := m.remove(elem)
	for _, other := range others {
		removed := m.remove(other)
		result = result || removed
	}
	return result
}

func (m *mutableMapSet[T]) remove(elem T) bool {
	_, ok := m.delegate[elem]
	delete(m.delegate, elem)
	return ok
}

func (m *mutableMapSet[T]) Contains(elem T) bool {
	_, ok := m.delegate[elem]
	return ok
}

func (m *mutableMapSet[T]) Len() int {
	return len(m.delegate)
}

func (m *mutableMapSet[T]) ForEach(fn func(elem T)) {
	for elem := range m.delegate {
		fn(elem)
	}
}

func (m *mutableMapSet[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for elem := range m.delegate {
			if !yield(elem) {
				return
			}
		}
	}
}

func (m *mutableMapSet[T]) String() string {
	return StringImpl[T](m)
}
