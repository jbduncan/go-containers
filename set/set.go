// Package set provides a set data structure, which is a generic, unordered container of elements where no two elements
// can be equal according to Go's '==' operator.
//
// A set can be created with the New function.
//
// Sets satisfy two interfaces: Set and MutableSet. Read those interfaces' docs and their methods' docs for more
// information.
//
// An existing set can be made "unmodifiable", which turns it into a read-only Set view. Read the docs for the function
// Unmodifiable for more information.
package set

// Set is a generic, unordered collection of unique elements. This interface has methods for reading the set; for
// writing to the set, use the MutableSet interface.
type Set[T comparable] interface {
	// Contains returns true if this set contains the given element, otherwise
	// it returns false.
	Contains(elem T) bool

	// Len returns the number of elements in this set.
	Len() int

	// ForEach runs the given function on each element in this set.
	//
	// The iteration order of the elements is undefined; it may even change
	// from one call to the next.
	ForEach(fn func(elem T))

	// String returns a string representation of all the elements in this set.
	//
	// The format of this string is a single "[" followed by a comma-separated
	// list (", ") of this set's elements in the same order as ForEach (which
	// is undefined and may change from one call to the next), followed by a
	// single "]".
	String() string

	// TODO: Set: make Iterator method that returns an Iterator.
	//       Note: this depends on any of:
	//         - The "range over func" proposal: https://github.com/golang/go/issues/61405
	//         - Making a custom map type that we can easily make an iterator for
	//         - Using reflect.MapIter
	// Iterator returns an iterator for the elements in this set.
	// Iterator() iterator.Iterator[T]

	// Equal returns true if this set has the same elements as the other set
	// in any order. Otherwise, it returns false.
	//
	// This method should be used over ==, the behaviour of which is
	// undefined.
	//
	// Equal implementations follow these rules:
	//   - Reflexive: for any potentially-nil set x, x.Equal(x) should be true.
	//   - Symmetric: for any potentially-nil sets x and y, x.Equal(y) and y.Equal(x) should have the same results.
	//   - Transitive: for any potentially-nil sets x, y and z, if x.Equal(y) and y.Equal(z), then x.Equal(z) should be
	//     true.
	//   - Consistent: for any potentially-nil sets x and y, multiple calls to x.Equal(y) should consistently return
	//     true or consistently return false, as long as no information used by the Equal calls is changed.
	Equal(other Set[T]) bool
}

// MutableSet is a Set with additional methods for adding elements to the set
// and removing them.
type MutableSet[T comparable] interface {
	Set[T]

	// Add adds the given element to this set if it is not already present.
	// Returns true if the element was not already present in the set,
	// otherwise false.
	Add(elem T) bool

	// Remove removes the given element from this set if it is present. Returns
	// true if the element was already present in the set, otherwise false.
	Remove(elem T) bool
}

// TODO: Consider returning a public version of the concrete type, rather
//       than the MutableSet interface, to allow new methods to be added
//       without breaking backwards compatibility:
//       - https://github.com/golang/go/wiki/CodeReviewComments#interfaces
//       In doing so, move Set.Equal to its own helper method (as gonum does
//       for its graph type) and make all implementations incomparable with ==
//       by using the same trick as:
//       https://github.com/tailscale/tailscale/blob/e5e5ebda44e7d28df279e89d3cc3a8b904843304/types/structs/structs.go

// New returns a new empty MutableSet.
func New[T comparable]() MutableSet[T] {
	return set[T]{}
}

type set[T comparable] map[T]struct{}

// TODO: If the Set and MutableSet interfaces are ever eliminated, move them and these
//       compile-time type assertions to a test package.

var (
	_ Set[int]        = (*set[int])(nil)
	_ MutableSet[int] = (*set[int])(nil)
)

func (s set[T]) Add(elem T) bool {
	_, ok := s[elem]
	s[elem] = struct{}{}
	return !ok
}

func (s set[T]) Remove(elem T) bool {
	_, ok := s[elem]
	delete(s, elem)
	return ok
}

func (s set[T]) Contains(elem T) bool {
	_, ok := s[elem]
	return ok
}

func (s set[T]) Len() int {
	return len(s)
}

func (s set[T]) ForEach(fn func(elem T)) {
	for elem := range s {
		fn(elem)
	}
}

func (s set[T]) String() string {
	return StringImpl[T](s)
}

func (s set[T]) Equal(other Set[T]) bool {
	// TODO
	panic("not yet implemented")
}
