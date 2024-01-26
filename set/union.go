package set

// Union returns the set union of sets a and b.
//
// The returned set is an unmodifiable view, so changes to a and b will be reflected in the returned set.
//
// Set.Len runs in O(b) time.
//
// Note: If passing in a MutableSet, Go will force the programmer to specify the generic type explicitly, like:
//
//	a := set.NewMutable(1)
//	b := set.NewMutable(2)
//	u := set.Union[int](a, b)
//	              ^^^^^
func Union[T comparable](a, b Set[T]) Set[T] {
	return union[T]{
		a: a,
		b: b,
	}
}

type union[T comparable] struct {
	a Set[T]
	b Set[T]
}

func (u union[T]) Contains(elem T) bool {
	return u.a.Contains(elem) || u.b.Contains(elem)
}

func (u union[T]) Len() int {
	bLen := 0
	u.b.ForEach(func(elem T) {
		if !u.a.Contains(elem) {
			bLen++
		}
	})
	return u.a.Len() + bLen
}

func (u union[T]) ForEach(fn func(elem T)) {
	u.a.ForEach(fn)
	u.b.ForEach(func(elem T) {
		if !u.a.Contains(elem) {
			fn(elem)
		}
	})
}

func (u union[T]) String() string {
	return StringImpl[T](u)
}
