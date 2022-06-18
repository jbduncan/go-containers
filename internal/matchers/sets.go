package matchers

import (
	"github.com/onsi/gomega/types"
	"go-containers/container/set"

	. "github.com/onsi/gomega"
)

func HaveLenOf[T comparable](len int) types.GomegaMatcher {
	return WithTransform(
		func(set set.Set[T]) int {
			return set.Len()
		},
		Equal(len))
}

func Contain[T comparable](elem T) types.GomegaMatcher {
	return WithTransform(
		func(set set.Set[T]) bool {
			return set.Contains(elem)
		},
		BeTrue())
}

func ForEachToSlice[T comparable](s set.Set[T]) []T {
	var result []T

	s.ForEach(func(elem T) {
		result = append(result, elem)
	})

	return result
}
