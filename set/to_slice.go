package set

// ToSlice returns a new slice with the contents of the argument set copied
// into it.
//
// This function is implemented in terms of Set.ForEach, so the order of the
// elements in the returned slice is undefined; it may even change from one
// call of ToSlice to the next.
func ToSlice[T comparable](s Set[T]) []T {
	result := make([]T, 0, s.Len())

	s.ForEach(func(elem T) {
		result = append(result, elem)
	})

	return result
}
