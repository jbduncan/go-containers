// Package set provides a set data structure, which is a generic, unordered container of elements where no two elements
// can be equal according to Go's == operator.
//
// A set can be created with the New function.
//
// An existing set can be made "unmodifiable", which turns it into a read-only Set view. Read the docs for the function
// Unmodifiable for more information.
//
// Two sets can be compared for "equality", returning true if they both have the same elements as each other in any
// order, otherwise false. Read the docs for the function Equal for more information.
//
// Third-party set implementations can be tested with settest.Set.
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

	// TODO: Set: make Iterator method that returns an Iterator.
	//       Note: this depends on any of:
	//         - The "range over func" proposal: https://github.com/golang/go/issues/61405
	//         - Making a custom map type that we can easily make an iterator for
	//         - Using reflect.MapIter
	// Iterator returns an iterator for the elements in this set.
	// Iterator() iterator.Iterator[T]
}

// MutableSet is a Set with additional methods for adding elements to the set and removing them.
type MutableSet[T comparable] interface {
	Set[T]

	// Add adds the given element to this set if it is not already present. Returns true if the element was not already
	// present in the set, otherwise false.
	Add(elem T) bool

	// Remove removes the given element from this set if it is present. Returns true if the element was already present
	// in the set, otherwise false.
	Remove(elem T) bool
}

// TODO: Make all set implementations incomparable with == by using the same trick as:
//       https://github.com/tailscale/tailscale/blob/e5e5ebda44e7d28df279e89d3cc3a8b904843304/types/structs/structs.go

// New returns a new non-nil, empty, mutable set.
func New[T comparable]() *MapSet[T] {
	return &MapSet[T]{
		delegate: map[T]struct{}{},
	}
}

// TODO: Introduce set.Of[T comparable](elems ...T), which returns an unmodifiable set
// TODO: Introduce set.OfMutable[T comparable](elems ...T) ...

// MapSet is a generic, unordered container of elements where no two elements can be equal according to Go's ==
// operator. It implements Set and MutableSet. Its implementation is based on a Go map, with similar performance
// characteristics.
type MapSet[T comparable] struct {
	delegate map[T]struct{}
}

var (
	_ Set[int]        = (*MapSet[int])(nil)
	_ MutableSet[int] = (*MapSet[int])(nil)
)

// Add adds the given element to this set if it is not already present. Returns true if the element was not already
// present in the set, otherwise false.
func (s *MapSet[T]) Add(elem T) bool {
	_, ok := s.delegate[elem]
	s.delegate[elem] = struct{}{}
	return !ok
}

// Remove removes the given element from this set if it is present. Returns true if the element was already present in
// the set, otherwise false.
func (s *MapSet[T]) Remove(elem T) bool {
	_, ok := s.delegate[elem]
	delete(s.delegate, elem)
	return ok
}

// Contains returns true if this set contains the given element, otherwise it returns false.
func (s *MapSet[T]) Contains(elem T) bool {
	_, ok := s.delegate[elem]
	return ok
}

// Len returns the number of elements in this set.
func (s *MapSet[T]) Len() int {
	return len(s.delegate)
}

// ForEach runs the given function on each element in this set.
//
// The iteration order of the elements is undefined; it may even change from one call to the next.
func (s *MapSet[T]) ForEach(fn func(elem T)) {
	for elem := range s.delegate {
		fn(elem)
	}
}

// String returns a string representation of all the elements in this set.
//
// The format of this string is a single "[" followed by a comma-separated list (", ") of this set's elements in the
// same order as ForEach (which is undefined and may change from one call to the next), followed by a single "]".
//
// This method satisfies fmt.Stringer.
func (s *MapSet[T]) String() string {
	return StringImpl[T](s)
}
