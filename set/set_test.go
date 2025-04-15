package set_test

import (
	"testing"

	"github.com/jbduncan/go-containers/set"
	"github.com/jbduncan/go-containers/set/settest"
)

func TestSetOf(t *testing.T) {
	t.Parallel()

	settest.TestSet(t, func(elements []int) settest.Set[int] {
		return set.Of(elements...)
	})
}

func TestSetInitializedWithAdd(t *testing.T) {
	t.Parallel()

	settest.TestSet(t, func(elements []int) settest.Set[int] {
		s := set.Of[int]()
		for _, element := range elements {
			s.Add(element)
		}
		return s
	})
}
