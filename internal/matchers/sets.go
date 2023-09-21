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
				// This error message is more descriptive than making `value` of type `sized` directly.
				return false, fmt.Errorf(
					"HaveLenOf matcher expected actual with Len method with <int> return type.  Got:\n%s",
					format.Object(value, 1))
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

func ContainAtLeast[T comparable](first T, others ...T) types.GomegaMatcher {
	elements := allOf(first, others)

	return gcustom.MakeMatcher(
		func(s set.Set[T]) (bool, error) {
			for _, element := range elements {
				if !s.Contains(element) {
					return false, nil
				}
			}
			return true, nil
		}).
		WithTemplate("Expected\n{{.FormattedActual}}\n{{.To}} contain at least\n{{format .Data 1}}").
		WithTemplateData(elements)
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
	elements := allOf(first, others)

	return gcustom.MakeMatcher(
		func(s set.Set[T]) (bool, error) {
			actual := forEachToSlice(s)
			return ConsistOf(elements).Match(actual)
		}).
		WithTemplate("Expected ForEach() of\n{{.FormattedActual}}\n{{.To}} emit elements consisting of\n{{format .Data 1}}").
		WithTemplateData(elements)
}

func HaveForEachThatConsistsOfElementsInSlice[T comparable](elements []T) types.GomegaMatcher {
	return gcustom.MakeMatcher(
		func(s set.Set[T]) (bool, error) {
			actual := forEachToSlice(s)
			return ConsistOf(elements).Match(actual)
		}).
		WithTemplate("Expected ForEach() of\n{{.FormattedActual}}\n{{.To}} emit elements consisting of\n{{format .Data 1}}").
		WithTemplateData(elements)
}

func HaveForEachThatConsistsOfElementsInSet[T comparable](s set.Set[T]) types.GomegaMatcher {
	elements := forEachToSlice(s)

	return gcustom.MakeMatcher(
		func(s2 set.Set[T]) (bool, error) {
			actual := forEachToSlice(s2)
			return ConsistOf(elements).Match(actual)
		}).
		WithTemplate("Expected ForEach() of\n{{.FormattedActual}}\n{{.To}} to emit elements consisting of\n{{format .Data 1}}").
		WithTemplateData(elements)
}

func BeNonMutableSet[T comparable]() types.GomegaMatcher {
	return gcustom.MakeMatcher(
		func(s set.Set[T]) (bool, error) {
			_, mutable := s.(set.MutableSet[T])
			return !mutable, nil
		}).
		WithTemplate("Expected\n{{.FormattedActual}}\n{{.To}} implement set.Set but not set.MutableSet")
}

func forEachToSlice[T comparable](s set.Set[T]) []T {
	var result []T

	s.ForEach(func(elem T) {
		result = append(result, elem)
	})

	return result
}

func allOf[T any](first T, others []T) []T {
	all := make([]T, 0, len(others)+1)
	all = append(all, first)
	all = append(all, others...)
	return all
}
