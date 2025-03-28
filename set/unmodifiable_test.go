package set_test

import (
	"testing"

	internalsettest "github.com/jbduncan/go-containers/internal/settest"
	"github.com/jbduncan/go-containers/set"
	"github.com/jbduncan/go-containers/set/settest"
)

func TestSetUnmodifiable(t *testing.T) {
	settest.Set(t, func(elements []string) set.Set[string] {
		s := set.Of[string]()
		for _, element := range elements {
			s.Add(element)
		}
		return set.Unmodifiable[string](s)
	})

	t.Run(
		"empty unmodifiable set: add to underlying set",
		func(t *testing.T) {
			s := set.Of[string]()
			unmodSet := set.Unmodifiable[string](s)

			s.Add("link")

			internalsettest.SetLen(t, "set.Unmodifiable", unmodSet, 1)
			internalsettest.SetAll(
				t,
				"set.Unmodifiable",
				unmodSet,
				[]string{"link"},
			)
			internalsettest.SetContains(
				t,
				"set.Unmodifiable",
				unmodSet,
				[]string{"link"},
				nil,
			)
			internalsettest.SetString(
				t,
				"set.Unmodifiable",
				unmodSet,
				[]string{"link"},
			)
		})

	t.Run(
		"empty unmodifiable set: "+
			"add x2 to underlying set: "+
			"has two-element string representation",
		func(t *testing.T) {
			s := set.Of[string]()
			unmodSet := set.Unmodifiable[string](s)

			s.Add("link")
			s.Add("zelda")

			internalsettest.SetString(
				t,
				"set.Unmodifiable",
				unmodSet,
				[]string{"link", "zelda"},
			)
		})
}
