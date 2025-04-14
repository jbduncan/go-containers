package set_test

import (
	"testing"

	"github.com/jbduncan/go-containers/set"
	"github.com/jbduncan/go-containers/set/settest"
)

func TestSetOf(t *testing.T) {
	t.Parallel()

	settest.TestSet(t, func(elems []string) settest.Set[string] {
		return set.Of(elems...)
	})
}

func TestSetInitializedWithAdd(t *testing.T) {
	t.Parallel()

	settest.TestSet(t, func(elems []string) settest.Set[string] {
		s := set.Of[string]()
		for _, elem := range elems {
			s.Add(elem)
		}
		return s
	})
}
