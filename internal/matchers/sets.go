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
	return gcustom.MakeMatcher(
		func(value any) (bool, error) {
			type sized interface {
				Len() int
			}

			s, ok := value.(sized)
			if !ok {
				return false, fmt.Errorf("HaveLenOf matcher expected actual with Len method with <int> return type.  Got:\n%s", format.Object(value, 1))
			}

			actualLen := s.Len()

			return actualLen == len, nil
		}).
		WithTemplate("Expected\n{{.FormattedActual}}\n{{.To}} have length\n{{format .Data 1}}").
		WithTemplateData(len)
}

func HaveLenOfZero() types.GomegaMatcher {
	return HaveLenOf(0)
}

func Contain[T comparable](elem T) types.GomegaMatcher {
	return gcustom.MakeMatcher(
		func(s set.Set[T]) (bool, error) {
			return s.Contains(elem), nil
		}).
		WithTemplate("Expected\n{{.FormattedActual}}\n{{.To}} contain\n{{format .Data 1}}").
		WithTemplateData(elem)
}

func HaveEmptyToSlice[T comparable]() types.GomegaMatcher {
	return gcustom.MakeMatcher(
		func(s set.Set[T]) (bool, error) {
			actual := s.ToSlice()
			return len(actual) == 0, nil
		}).
		WithTemplate("Expected ToSlice() of\n{{.FormattedActual}}\n{{.To}} return an empty slice")
}

func HaveToSliceThatConsistsOf[T comparable](first T, others ...T) types.GomegaMatcher {
	all := allOf(first, others)

	return gcustom.MakeMatcher(
		func(s set.Set[T]) (bool, error) {
			actual := s.ToSlice()
			return ConsistOf(all).Match(actual)
		}).
		WithTemplate("Expected ToSlice() of\n{{.FormattedActual}}\n{{.To}} consist of\n{{format .Data 1}}").
		WithTemplateData(all)
}

func HaveForEachThatEmitsNothing[T comparable]() types.GomegaMatcher {
	return gcustom.MakeMatcher(
		func(s set.Set[T]) (bool, error) {
			actual := forEachToSlice(s)
			return len(actual) == 0, nil
		}).
		WithTemplate("Expected ForEach() of\n{{.FormattedActual}}\n{{.To}} emit nothing")
}

func HaveForEachThatConsistsOf[T comparable](first any, others ...any) types.GomegaMatcher {
	all := allOf(first, others)

	return gcustom.MakeMatcher(
		func(s set.Set[T]) (bool, error) {
			actual := forEachToSlice(s)
			return ConsistOf(all).Match(actual)
		}).
		WithTemplate("Expected ForEach() of\n{{.FormattedActual}}\n{{.To}} to emit elements consisting of\n{{format .Data 1}}").
		WithTemplateData(all)
}

func HaveForEachThatConsistsOfElementsIn[T comparable](set set.Set[T]) types.GomegaMatcher {
	// TODO: Use gcustom.MakeMatcher to improve error message
	return WithTransform(forEachToSlice[T], ConsistOf(forEachToSlice(set)))
}

func BeNonMutableSet[T comparable]() types.GomegaMatcher {
	// TODO: Use gcustom.MakeMatcher to improve error message
	return WithTransform(
		func(s set.Set[T]) bool {
			_, mutable := s.(set.MutableSet[T])
			return mutable
		},
		BeFalse())
}

func forEachToSlice[T comparable](s set.Set[T]) []T {
	var result []T

	s.ForEach(func(elem T) {
		result = append(result, elem)
	})

	return result
}

func allOf[T any](first T, others []T) []T {
	all := []T{first}
	all = append(all, others...)
	return all
}
