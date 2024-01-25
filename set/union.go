package set

func Union[T comparable](a Set[T], b Set[T]) Set[T] {
	return &union[T]{
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
