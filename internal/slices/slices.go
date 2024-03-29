package slices

import "golang.org/x/exp/slices"

func CartesianProduct[T any](values [][]T) [][]T {
	result := [][]T{{}}
	for _, innerValues := range values {
		var newResult [][]T
		for _, rest := range result {
			for _, tail := range innerValues {
				newResult = append(newResult, CopyToNonNilSlice(append(rest, tail)))
			}
		}
		result = newResult
	}
	return result
}

func CopyToNonNilSlice[T any](values []T) []T {
	if values == nil {
		return []T{}
	}

	return slices.Clone(values)
}

func Repeat[T any](value T, times int) []T {
	result := make([]T, times)
	for i := range result {
		result[i] = value
	}
	return result
}
