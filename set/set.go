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

import (
	"fmt"
	"strings"
)

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
	// The order in which the elements are returned is undefined; it may even
	// change from one call to the next.
	ForEach(fn func(elem T))

	// TODO: Can we move all of our Set.String() implementations into a helper
	//       function?

	// String returns a string representation of all the elements in this set.
	//
	// The format of this string is undefined. The order of the elements in
	// this string is also undefined; it may even change from one call to the
	// next.
	String() string

	// TODO: Set: make Iterator method that returns an Iterator.
	//       Note: this depends on making a custom map type that we can easily make an
	//       iterator for.
	// Iterator returns an iterator for the elements in this set.
	// Iterator() iterator.Iterator[T]

	// TODO: Set: make Equal method and discourage == from being used (documenting that its use is undefined).
}

// MutableSet is a Set with additional methods for adding and removing
// elements to and from the set.
type MutableSet[T comparable] interface {
	Set[T]

	// Add adds the given element into this set, if it is not already present.
	Add(elem T) // TODO: Return true if set did not already have the element, otherwise false

	// Remove removes the given element from this set.
	Remove(elem T) // TODO: Return true if set had the element, otherwise false
}

// TODO: Consider returning a public version of the concrete type, rather
//       than the MutableSet interface, to allow new methods to be added
//       without breaking backwards compatibility:
//       - https://github.com/golang/go/wiki/CodeReviewComments#interfaces

// New returns a new empty MutableSet.
func New[T comparable]() MutableSet[T] {
	return &set[T]{
		delegate: map[T]struct{}{},
	}
}

type set[T comparable] struct {
	delegate map[T]struct{}
}

// TODO: If the Set and MutableSet interfaces are ever eliminated, move them and these
//       compile-time type assertions to a test package.

var (
	_ Set[int]        = (*set[int])(nil)
	_ MutableSet[int] = (*set[int])(nil)
)

func (s *set[T]) Add(elem T) {
	s.delegate[elem] = struct{}{}
}

func (s *set[T]) Remove(elem T) {
	delete(s.delegate, elem)
}

func (s *set[T]) Contains(elem T) bool {
	_, ok := s.delegate[elem]
	return ok
}

func (s *set[T]) Len() int {
	return len(s.delegate)
}

func (s *set[T]) ForEach(fn func(elem T)) {
	for elem := range s.delegate {
		fn(elem)
	}
}

func (s *set[T]) String() string {
	var builder strings.Builder

	builder.WriteRune('[')
	index := 0
	for elem := range s.delegate {
		if index > 0 {
			builder.WriteString(", ")
		}

		builder.WriteString(fmt.Sprintf("%v", elem))
		index++
	}

	builder.WriteRune(']')
	return builder.String()
}
