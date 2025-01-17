package slicesx

import "slices"

func AllOf[T any](first T, rest []T) []T {
	return slices.Concat([]T{first}, rest)
}
