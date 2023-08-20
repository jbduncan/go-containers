package set_test

import (
	"testing"

	. "github.com/jbduncan/go-containers/internal/matchers"
	"github.com/jbduncan/go-containers/set"
	"github.com/jbduncan/go-containers/set/settest"
	. "github.com/onsi/gomega"
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

func TestSetUnmodifiable(t *testing.T) {
	settest.Set(t, func(elements []string) set.Set[string] {
		s := set.New[string]()
		for _, element := range elements {
			s.Add(element)
		}
		return set.Unmodifiable[string](s)
	})

	t.Run(
		"empty unmodifiable set: add to underlying set: has length of 1",
		func(t *testing.T) {
			g := NewWithT(t)
			s := set.New[string]()
			unmodSet := set.Unmodifiable(s)

			s.Add("link")

			g.Expect(unmodSet).To(HaveLenOf(1))
		})

	t.Run(
		"empty unmodifiable set: add to underlying set: contains element",
		func(t *testing.T) {
			g := NewWithT(t)
			s := set.New[string]()
			unmodSet := set.Unmodifiable(s)

			s.Add("link")

			g.Expect(unmodSet).To(Contain("link"))
		})

	t.Run(
		"empty unmodifiable set: "+
			"add to underlying set: "+
			"returns element on iteration",
		func(t *testing.T) {
			g := NewWithT(t)
			s := set.New[string]()
			unmodSet := set.Unmodifiable(s)

			s.Add("link")

			g.Expect(unmodSet).To(HaveForEachThatConsistsOf[string]("link"))
		})

	t.Run(
		"empty unmodifiable set: "+
			"add to underlying set: "+
			"has one-element string representation",
		func(t *testing.T) {
			g := NewWithT(t)
			s := set.New[string]()
			unmodSet := set.Unmodifiable(s)

			s.Add("link")

			g.Expect(unmodSet).To(HaveStringRepr("[link]"))
		})

	t.Run(
		"empty unmodifiable set: "+
			"add to underlying set: "+
			"returns one-element slice",
		func(t *testing.T) {
			g := NewWithT(t)
			s := set.New[string]()
			unmodSet := set.Unmodifiable(s)

			s.Add("link")

			g.Expect(unmodSet).To(HaveToSliceThatConsistsOf("link"))
		})

	t.Run(
		"empty unmodifiable set: "+
			"add x2 to underlying set: "+
			"has two-element string representation",
		func(t *testing.T) {
			g := NewWithT(t)
			s := set.New[string]()
			unmodSet := set.Unmodifiable(s)

			s.Add("link")
			s.Add("zelda")

			g.Expect(unmodSet).To(HaveStringRepr(BeElementOf("[link, zelda]", "[zelda, link]")))
		})

	t.Run(
		"empty unmodifiable set: "+
			"add x2 to underlying set: "+
			"returns two-element slice",
		func(t *testing.T) {
			g := NewWithT(t)
			s := set.New[string]()
			unmodSet := set.Unmodifiable(s)

			s.Add("link")
			s.Add("zelda")

			g.Expect(unmodSet).To(HaveToSliceThatConsistsOf("link", "zelda"))
		})
}
