// Package set provides a set data structure, which is a generic, unordered container of elements where no two elements
// can be equal according to Go's == operator.
//
// A set can be created with the New function.
//
// An existing set can be made "unmodifiable" with Unmodifiable, which turns it into a read-only Set view.
//
// Two sets can be compared for "equality" with Equal, returning true if they both have the same elements as each other
// in any order, otherwise false.
//
// The contents of a set can be copied into a slice with ToSlice.
//
// Third-party set implementations can be tested with settest.Set.
package set
