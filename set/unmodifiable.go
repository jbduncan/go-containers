package set

type unmodifiableSet[T comparable] struct {
	set MutableSet[T]
}

var _ Set[int] = (*unmodifiableSet[int])(nil)

// Unmodifiable wraps a given MutableSet as a read-only Set view.
//
// This set cannot be cast back into a MutableSet.
//
// If the original MutableSet is ever mutated, then the returned Set will reflect those mutations.
//
// This function allows for two use cases:
//   - Immutable sets.
//   - Sets that can be mutated by your own code but not your clients' code.
func Unmodifiable[T comparable](set MutableSet[T]) Set[T] {
	return unmodifiableSet[T]{
		set: set,
	}
}

func (u unmodifiableSet[T]) Contains(elem T) bool {
	return u.set.Contains(elem)
}

func (u unmodifiableSet[T]) Len() int {
	return u.set.Len()
}

func (u unmodifiableSet[T]) ForEach(fn func(elem T)) {
	u.set.ForEach(fn)
}

func (u unmodifiableSet[T]) String() string {
	return u.set.String()
}

func (u unmodifiableSet[T]) Equal(other Set[T]) bool {
	// TODO
	panic("not yet implemented")
}
