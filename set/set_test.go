package set_test

import (
	"testing"

	"github.com/jbduncan/go-containers/set"
	"github.com/jbduncan/go-containers/set/settest"
)

func TestSetOf(t *testing.T) {
	settest.Set(t, func(elems []string) set.Set[string] {
		return set.Of(elems...)
	})
}

func TestSetNewMutableInitializedWithAdd(t *testing.T) {
	settest.Set(t, func(elems []string) set.Set[string] {
		s := set.NewMutable[string]()
		for _, elem := range elems {
			s.Add(elem)
		}
		return s
	})
}

func TestSetNewMutableInitializedWithInit(t *testing.T) {
	settest.Set(t, func(elems []string) set.Set[string] {
		return set.NewMutable[string](elems...)
	})
}
