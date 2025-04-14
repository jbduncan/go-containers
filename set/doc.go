// Package set provides a set data structure, which is a generic, unordered container of elements where no two elements
// can be equal according to Go's == operator.
//
// A mutable set can be created with Of, which can be passed into Unmodifiable to turn it into a read-only set view.
//
// Set has a String method that satisfies fmt.Stringer.
//
// The union of two sets can be created with Union.
//
// Two sets can be compared for "equality" with Equal, returning true if they have the same elements in any order,
// otherwise false.
//
// Third-party set implementations can be tested with settest.Set.
package set
