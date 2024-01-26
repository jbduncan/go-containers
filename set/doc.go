// Package set provides a set data structure, which is a generic, unordered container of elements where no two elements
// can be equal according to Go's == operator.
//
// An immutable set can be created with Of. Likewise, a mutable set can be created with NewMutable.
//
// A mutable set can be passed into Unmodifiable, which turns it into a read-only Set view.
//
// Set has a String method that satisfies fmt.Stringer.
//
// The union of two sets can be created with Union.
//
// Two sets can be compared for "equality" with Equal, returning true if they have the same elements in any order,
// otherwise false.
//
// The contents of a set can be copied into a slice with ToSlice.
//
// Third-party set implementations can be tested with settest.Set.
package set
