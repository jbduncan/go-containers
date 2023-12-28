package set_test

import (
	"testing"

	"github.com/jbduncan/go-containers/set"
	"github.com/jbduncan/go-containers/set/settest"
)

func TestSetOf(t *testing.T) {
	settest.Set(t, func(elements []string) set.Set[string] {
		return set.Of(elements...)
	})
}

func TestSetNewMutable(t *testing.T) {
	settest.Set(t, func(elements []string) set.Set[string] {
		s := set.NewMutable[string]()
		for _, element := range elements {
			s.Add(element)
		}
		return s
	})
}
