package set_test

import (
	"testing"

	internalsettest "github.com/jbduncan/go-containers/internal/settest"
	"github.com/jbduncan/go-containers/set"
	"github.com/jbduncan/go-containers/set/settest"
)

func TestSetUnmodifiable(t *testing.T) {
	t.Parallel()

	settest.TestSet(t, func(elements []int) settest.Set[int] {
		s := set.Of[int]()
		for _, element := range elements {
			s.Add(element)
		}
		return set.Unmodifiable[int](s)
	})

	t.Run(
		"empty unmodifiable set: add to underlying set",
		func(t *testing.T) {
			t.Parallel()

			s := set.Of[int]()
			unmodSet := set.Unmodifiable[int](s)

			s.Add(1)

			internalsettest.Len(t, "set.Unmodifiable", unmodSet, 1)
			internalsettest.All(
				t,
				"set.Unmodifiable",
				unmodSet,
				[]int{1},
			)
			internalsettest.Contains(
				t,
				"set.Unmodifiable",
				unmodSet,
				[]int{1},
			)
			internalsettest.String(
				t,
				"set.Unmodifiable",
				unmodSet,
				[]int{1},
			)
		})

	t.Run(
		"empty unmodifiable set: "+
			"add x2 to underlying set: "+
			"has two-element string representation",
		func(t *testing.T) {
			t.Parallel()

			s := set.Of[string]()
			unmodSet := set.Unmodifiable[string](s)

			s.Add("link")
			s.Add("zelda")

			internalsettest.String(
				t,
				"set.Unmodifiable",
				unmodSet,
				[]string{"link", "zelda"},
			)
		})
}
