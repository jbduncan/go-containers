package container

type Set[T comparable] interface {
	Contains(value T) bool
	Len() int
	ForEach(fn func(elem T))
	// Iter() Iterator
	String() string
}

// TODO: Implement in another file
// type Iterator[T any] interface {
// 	Value() T
// 	Next() bool
// }

type MutableSet[T comparable] interface {
	Set[T]

	Add(value T)
	Remove(value T)
}

func NewSet[T comparable]() MutableSet[T] {
	return new(set[T])
}

type set[T comparable] struct {
	elems []T
} // map[T]struct{}

var _ MutableSet[int] = (*set[int])(nil)

func (s *set[T]) Add(value T) {
	// s[v] = struct{}{}
	if !s.Contains(value) {
		s.elems = append(s.elems, value)
	}
}

func (s set[T]) Remove(value T) {
	// delete(s, v)
}

func (s set[T]) Contains(value T) bool {
	// _, ok := s[v]
	// return ok

	//return (s.elem != nil && *s.elem == value) || (s.elem2 != nil && *s.elem2 == value)
	for _, elem := range s.elems {
		if elem == value {
			return true
		}
	}
	return false
}

func (s set[T]) Len() int {
	// return len(s)
	return len(s.elems)
}

func (s set[T]) ForEach(fn func(elem T)) {
	// for v := range s {
	// 	f(v)
	// }
}

func (s set[T]) String() string {
	// var b strings.Builder
	// b.WriteRune('[')
	// i := 0
	// s.ForEach(func(v T) {
	// 	if i > 0 {
	// 		b.WriteString(", ")
	// 	}
	// 	b.WriteString(fmt.Sprintf("%v", v))
	// 	i++
	// })
	// b.WriteRune(']')
	// return b.String()
	return ""
}

type unmodifiableSet[T comparable] struct {
	set MutableSet[T]
}

var _ Set[int] = (*unmodifiableSet[int])(nil)

func (r unmodifiableSet[T]) Contains(value T) bool {
	// return r.set.Contains(v)
	return false
}

func (r unmodifiableSet[T]) Len() int {
	// return r.set.Len()
	return 0
}

func (r unmodifiableSet[T]) ForEach(fn func(elem T)) {
	// r.set.ForEach(f)
}

func (r unmodifiableSet[T]) String() string {
	// return r.set.String()
	return ""
}

func UnmodifiableSet[T comparable](set Set[T]) Set[T] {
	return nil
}
