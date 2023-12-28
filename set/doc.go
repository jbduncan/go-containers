// Package set provides a set data structure, which is a generic, unordered container of elements where no two elements
// can be equal according to Go's == operator.
//
// An immutable set can be created with Of. A mutable set can be created with NewMutable.
//
// An existing mutable set can be passed into Unmodifiable, which turns it into a read-only Set view.
//
// Two sets can be compared for "equality" with Equal, returning true if they both have the same elements in any order,
// otherwise false.
//
// The contents of a set can be copied into a slice with ToSlice.
//
// Third-party set implementations can be tested with settest.Set.
package set
