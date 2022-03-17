package container

type Set[T comparable] interface {
	Has(v T) bool
	Len() int
	ForEach(f func(T))
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

	Add(v T)
	Remove(v T)
	AsReadOnly() Set[T]
}

func NewSet[T comparable]() MutableSet[T] {
	// return make(set[T])
	return nil
}

type set[T comparable] map[T]struct{}

var _ MutableSet[int] = (*set[int])(nil)

func (s set[T]) Add(v T) {
	// s[v] = struct{}{}
}

func (s set[T]) Remove(v T) {
	// delete(s, v)
}

func (s set[T]) Has(v T) bool {
	// _, ok := s[v]
	// return ok
	return false
}

func (s set[T]) Len() int {
	// return len(s)
	return 0
}

func (s set[T]) ForEach(f func(T)) {
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

func (s set[T]) AsReadOnly() Set[T] {
	// return readOnlySet[T]{set: s}
	return nil
}

type readOnlySet[T comparable] struct {
	set MutableSet[T]
}

var _ Set[int] = (*readOnlySet[int])(nil)

func (r readOnlySet[T]) Has(v T) bool {
	// return r.set.Contains(v)
	return false
}

func (r readOnlySet[T]) Len() int {
	// return r.set.Len()
	return 0
}

func (r readOnlySet[T]) ForEach(f func(T)) {
	// r.set.ForEach(f)
}

func (r readOnlySet[T]) String() string {
	// return r.set.String()
	return ""
}
