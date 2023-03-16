package iterator

// Iterator generates a series of elements, one by one. Usually it gets its elements from an underlying collection.
// TODO: Document with an example use.
type Iterator[T any] interface {
	// Value returns the next element in the iteration. If the iteration has no more elements, it panics.
	Value() T

	// Next returns true if the iteration has more elements, otherwise false.
	Next() bool
}
