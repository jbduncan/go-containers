package set_test

import (
	"testing"

	"github.com/jbduncan/go-containers/set"
	"github.com/jbduncan/go-containers/set/settest"
)

func TestSetNew(t *testing.T) {
	settest.Set(t, func(elements []string) set.Set[string] {
		s := set.New[string]()
		for _, element := range elements {
			s.Add(element)
		}
		return s
	})
}
