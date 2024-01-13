package set

// Set is a generic, unordered collection of unique elements. This interface has methods for reading the set; for
// writing to the set, use the MutableSet interface.
type Set[T comparable] interface {
	// Contains returns true if this set contains the given element, otherwise it returns false.
	Contains(elem T) bool

	// Len returns the number of elements in this set.
	Len() int

	// ForEach runs the given function on each element in this set.
	//
	// The iteration order of the elements is undefined; it may even change from one call to the next.
	ForEach(fn func(elem T))

	// String returns a string representation of all the elements in this set.
	//
	// The format of this string is a single "[" followed by a comma-separated list (", ") of this set's elements in
	// the same order as ForEach (which is undefined and may change from one call to the next), followed by a single
	// "]".
	//
	// This method satisfies fmt.Stringer.
	String() string
}

// MutableSet is a Set with additional methods for adding elements to the set and removing them.
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
	_ Set[int]        = (*MutableMapSet[int])(nil)
	_ MutableSet[int] = (*MutableMapSet[int])(nil)
)

// Of returns a new non-nil, empty, immutable set. See MapSet for more details.
func Of[T comparable](elems ...T) *MapSet[T] {
	delegate := map[T]struct{}{}
	for _, elem := range elems {
		delegate[elem] = struct{}{}
	}
	return &MapSet[T]{
		delegate: delegate,
	}
}

// NewMutable returns a new non-nil, empty, mutable set. See MutableMapSet for more details.
func NewMutable[T comparable](elems ...T) *MutableMapSet[T] {
	delegate := map[T]struct{}{}
	for _, elem := range elems {
		delegate[elem] = struct{}{}
	}
	return &MutableMapSet[T]{
		delegate: delegate,
	}
}

// MapSet is an immutable, generic, unordered container of elements where no two elements can be equal according to
// Go's == operator. It implements Set and MutableSet. Its implementation is based on a Go map, with similar
// performance characteristics.
type MapSet[T comparable] struct {
	delegate map[T]struct{}
}

// Contains returns true if this set contains the given element, otherwise it returns false.
func (m *MapSet[T]) Contains(elem T) bool {
	_, ok := m.delegate[elem]
	return ok
}

// Len returns the number of elements in this set.
func (m *MapSet[T]) Len() int {
	return len(m.delegate)
}

// ForEach runs the given function on each element in this set.
//
// The iteration order of the elements is undefined; it may even change from one call to the next.
func (m *MapSet[T]) ForEach(fn func(elem T)) {
	for elem := range m.delegate {
		fn(elem)
	}
}

// String returns a string representation of all the elements in this set.
//
// The format of this string is a single "[" followed by a comma-separated list (", ") of this set's elements in the
// same order as ForEach (which is undefined and may change from one call to the next), followed by a single "]".
//
// This method satisfies fmt.Stringer.
func (m *MapSet[T]) String() string {
	return StringImpl[T](m)
}

// MutableMapSet is a mutable, generic, unordered container of elements where no two elements can be equal according to
// Go's == operator. It implements Set and MutableSet. Its implementation is based on a Go map, with similar performance
// characteristics.
type MutableMapSet[T comparable] struct {
	delegate map[T]struct{}
}

// Add adds the given element(s) to this set. If any of the elements are already present, the set will not add those
// elements again. Returns true if this set changed as a result of this call, otherwise false.
func (m *MutableMapSet[T]) Add(elem T, others ...T) bool {
	result := m.add(elem)
	for _, elem := range others {
		added := m.add(elem)
		result = result || added
	}
	return result
}

func (m *MutableMapSet[T]) add(elem T) bool {
	_, ok := m.delegate[elem]
	m.delegate[elem] = struct{}{}
	return !ok
}

// Remove removes the given element(s) from this set. If any of the elements are already absent, the set will not
// attempt to remove those elements. Returns true if this set changed as a result of this call, otherwise false.
func (m *MutableMapSet[T]) Remove(elem T, others ...T) bool {
	result := m.remove(elem)
	for _, other := range others {
		removed := m.remove(other)
		result = result || removed
	}
	return result
}

func (m *MutableMapSet[T]) remove(elem T) bool {
	_, ok := m.delegate[elem]
	delete(m.delegate, elem)
	return ok
}

// Contains returns true if this set contains the given element, otherwise it returns false.
func (m *MutableMapSet[T]) Contains(elem T) bool {
	_, ok := m.delegate[elem]
	return ok
}

// Len returns the number of elements in this set.
func (m *MutableMapSet[T]) Len() int {
	return len(m.delegate)
}

// ForEach runs the given function on each element in this set.
//
// The iteration order of the elements is undefined; it may even change from one call to the next.
func (m *MutableMapSet[T]) ForEach(fn func(elem T)) {
	for elem := range m.delegate {
		fn(elem)
	}
}

// String returns a string representation of all the elements in this set.
//
// The format of this string is a single "[" followed by a comma-separated list (", ") of this set's elements in the
// same order as ForEach (which is undefined and may change from one call to the next), followed by a single "]".
//
// This method satisfies fmt.Stringer.
func (m *MutableMapSet[T]) String() string {
	return StringImpl[T](m)
}
