package matchers

import (
	"fmt"
	"slices"

	"github.com/jbduncan/go-containers/internal/slicesx"
	"github.com/jbduncan/go-containers/set"
	// dot importing gomega matchers is best practice and this package is used by test code only
	. "github.com/onsi/gomega" //nolint:stylecheck
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/gcustom"
	"github.com/onsi/gomega/types"
)

func HaveLenOf(length int) types.GomegaMatcher {
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

			return actualLen == length, nil
		}).
		WithTemplate("Expected\n{{.FormattedActual}}\n{{.To}} have length\n{{format .Data 1}}").
		WithTemplateData(length)
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

func HaveAllThatEmitsNothing[T comparable]() types.GomegaMatcher {
	return gcustom.MakeMatcher(
		func(s set.Set[T]) (bool, error) {
			actual := slices.Collect(s.All())
			return len(actual) == 0, nil
		}).
		WithTemplate("Expected All() of\n{{.FormattedActual}}\n{{.To}} emit nothing")
}

func HaveAllThatConsistsOf[T comparable](first any, rest ...any) types.GomegaMatcher {
	elements := slicesx.AllOf(first, rest)

	return gcustom.MakeMatcher(
		func(s set.Set[T]) (bool, error) {
			actual := slices.Collect(s.All())
			return ConsistOf(elements).Match(actual)
		}).
		WithTemplate("Expected All() of\n{{.FormattedActual}}\n{{.To}} emit elements consisting of\n{{format .Data 1}}").
		WithTemplateData(elements)
}

func HaveAllThatConsistsOfElementsInSlice[T comparable](elements []T) types.GomegaMatcher {
	return gcustom.MakeMatcher(
		func(s set.Set[T]) (bool, error) {
			actual := slices.Collect(s.All())
			return ConsistOf(elements).Match(actual)
		}).
		WithTemplate("Expected All() of\n{{.FormattedActual}}\n{{.To}} emit elements consisting of\n{{format .Data 1}}").
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
