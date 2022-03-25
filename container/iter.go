package container

type Iterator[T any] interface {
	Value() T
	Next() bool
}
