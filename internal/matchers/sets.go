package matchers

import (
	"fmt"

	"github.com/jbduncan/go-containers/set"
	//lint:ignore ST1001 dot importing gomega matchers is best practice and
	// this package is used by test code only
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/gcustom"
	"github.com/onsi/gomega/types"
)

func HaveLenOf(len int) types.GomegaMatcher {
	return gcustom.MakeMatcher(func(value any) (bool, error) {
		type sized interface {
			Len() int
		}

		s, ok := value.(sized)
		if !ok {
			return false, fmt.Errorf("HaveLenOf matcher expected actual with Len method with <int> return type.  Got:\n%s", format.Object(value, 1))
		}

		actualLen := s.Len()

		return actualLen == len, nil
	}).WithTemplate("Expected\n{{.FormattedActual}}\n {{.To}} have length\n{{format .Data 1}}").
		WithTemplateData(len)
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
