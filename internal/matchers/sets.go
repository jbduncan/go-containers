package matchers

import (
	"fmt"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
	"go-containers/container/set"
)

func HaveLenOf(len int) types.GomegaMatcher {
	// TODO: Use gcustom.MakeMatcher to improve error message
	//  when value.Len() (the actual value) isn't equal to len
	//  (the expected value).
	return WithTransform(
		func(value any) (int, error) {
			errNoLenMethod := fmt.Errorf(format.Message(value, "to have a Len method with a single return value of type <int>"))

			type sized interface {
				Len() int
			}

			s, ok := value.(sized)
			if !ok {
				return 0, errNoLenMethod
			}

			return s.Len(), nil
		},
		Equal(len))
}

func HaveLenOfZero() types.GomegaMatcher {
	return HaveLenOf(0)
}

func Contain[T comparable](elem T) types.GomegaMatcher {
	// TODO: Use gcustom.MakeMatcher to improve error message
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

func HaveEmptyToSlice[T comparable]() types.GomegaMatcher {
	// TODO: Use gcustom.MakeMatcher to improve error message
	return WithTransform(
		func(s set.Set[T]) []T {
			return s.ToSlice()
		},
		BeEmpty())
}

func HaveToSliceThatConsistsOf[T comparable](first T, others ...T) types.GomegaMatcher {
	all := []T{first}
	all = append(all, others...)

	// TODO: Use gcustom.MakeMatcher to improve error message
	return WithTransform(
		func(s set.Set[T]) []T {
			return s.ToSlice()
		},
		ConsistOf(all))
}

// TODO: Rename to HaveForEachThatConsistsOf.

func BeSetWithForEachThatProduces(first string, others ...string) types.GomegaMatcher {
	all := []string{first}
	all = append(all, others...)

	// TODO: Use gcustom.MakeMatcher to improve error message
	return WithTransform(ForEachToSlice[string], ConsistOf(all))
}

// TODO: Rename to HaveForEachThatProducesNothing

func BeSetWithForEachThatProducesNothing() types.GomegaMatcher {
	// TODO: Use gcustom.MakeMatcher to improve error message
	return WithTransform(ForEachToSlice[string], BeEmpty())
}

// TODO: This matcher is a duplicate of BeSetWithForEachThatProduces.
//       Eliminate one or the other.

func BeSetThatConsistsOf[T comparable](first any, others ...any) types.GomegaMatcher {
	all := []any{first}
	all = append(all, others...)

	// TODO: Use gcustom.MakeMatcher to improve error message
	return WithTransform(ForEachToSlice[T], ConsistOf(all))
}
