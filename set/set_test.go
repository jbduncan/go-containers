package set_test

import (
	"testing"

	"github.com/jbduncan/go-containers/set"
	"github.com/jbduncan/go-containers/set/settest"
)

func TestSetOf(t *testing.T) {
	t.Parallel()

	settest.TestMutable(t, func(elements []int) settest.MutableSet[int] {
		return set.Of(elements...)
	})
}

func TestSetInitializedWithAdd(t *testing.T) {
	t.Parallel()

	settest.TestMutable(t, func(elements []int) settest.MutableSet[int] {
		s := set.Of[int]()
		for _, element := range elements {
			s.Add(element)
		}
		return s
	})
}
