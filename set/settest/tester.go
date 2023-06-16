package settest

import (
	"testing"

	. "github.com/jbduncan/go-containers/internal/matchers"
	"github.com/jbduncan/go-containers/set"
	. "github.com/onsi/gomega"
)

// TODO: Document
func Set(t *testing.T, setBuilder func(elements []string) set.Set[string]) {
	_, mutable := setBuilder(empty()).(set.MutableSet[string])

	emptySetHasLengthOfZero(t, setBuilder)

	emptySetContainsNothing(t, setBuilder)

	emptySetIterationDoesNothing(t, setBuilder)

	emptySetHasEmptyStringRepr(t, setBuilder)

	emptySetReturnsEmptySlice(t, setBuilder)

	emptySetRemoveDoesNothing(t, setBuilder, mutable)

	oneElementSetHasLengthOfOne(t, setBuilder)
	emptySetPlusOneHasLengthOfOne(t, setBuilder, mutable)

	oneElementSetContainsPresentElement(t, setBuilder)
	emptySetPlusOneContainsPresentElement(t, setBuilder, mutable)

	oneElementSetDoesNotContainAbsentElement(t, setBuilder)
	emptySetPlusOneDoesNotContainAbsentElement(t, setBuilder, mutable)

	oneElementSetReturnsElementOnIteration(t, setBuilder)
	emptySetPlusOneReturnsElementOnIteration(t, setBuilder, mutable)

	oneElementSetHasOneElementStringRepr(t, setBuilder)
	emptySetPlusOneHasOneElementStringRepr(t, setBuilder, mutable)

	oneElementSetReturnsOneElementSlice(t, setBuilder)
	emptySetPlusOneReturnsOneElementSlice(t, setBuilder, mutable)

	emptySetPlusOneMinusOneDoesNotContainAnything(t, setBuilder, mutable)

	twoElementSetHasLengthOfTwo(t, setBuilder)
	emptySetPlusTwoHasLengthOfTwo(t, setBuilder, mutable)

	twoElementSetContainsBothElements(t, setBuilder)
	emptySetPlusTwoContainsBothElements(t, setBuilder, mutable)

	twoElementSetReturnsBothElementsOnIteration(t, setBuilder)
	emptySetPlusTwoReturnsBothElementsOnIteration(t, setBuilder, mutable)

	twoElementSetHasTwoElementStringRepr(t, setBuilder)
	emptySetPlusTwoReturnsTwoElementStringRepr(t, setBuilder, mutable)

	twoElementSetReturnsTwoElementSlice(t, setBuilder)
	emptySetPlusTwoReturnsTwoElementSlice(t, setBuilder, mutable)

	emptySetPlusTwoMinusOneHasLengthOfOne(t, setBuilder, mutable)

	threeElementSetContainsAllThreeElements(t, setBuilder)
	emptySetPlusThreeContainsAllThreeElements(t, setBuilder, mutable)

	threeElementSetHasThreeElementStringRepr(t, setBuilder)
	emptySetPlusThreeHasThreeElementStringRepr(t, setBuilder, mutable)

	setInitializedFromTwoOfSameElementHasLengthOfOne(t, setBuilder)
	emptySetPlusSameElementTwiceHasLengthOfOne(t, setBuilder, mutable)
}

func emptySetHasLengthOfZero(t *testing.T, setBuilder func(elements []string) set.Set[string]) {
	t.Helper()

	t.Run("empty set: has length of 0", func(t *testing.T) {
		g := NewWithT(t)
		s := setBuilder(empty())

		g.Expect(s).To(HaveLenOfZero())
	})
}

func emptySetContainsNothing(t *testing.T, setBuilder func(elements []string) set.Set[string]) {
	t.Helper()

	t.Run("empty set: contains nothing", func(t *testing.T) {
		g := NewWithT(t)
		s := setBuilder(empty())

		g.Expect(s).ToNot(Contain("link"))
	})
}

func emptySetIterationDoesNothing(t *testing.T, setBuilder func(elements []string) set.Set[string]) {
	t.Helper()

	t.Run("empty set: iteration does nothing", func(t *testing.T) {
		g := NewWithT(t)
		s := setBuilder(empty())

		g.Expect(s).To(BeSetWithForEachThatProducesNothing())
	})
}

func emptySetHasEmptyStringRepr(t *testing.T, setBuilder func(elements []string) set.Set[string]) {
	t.Helper()

	t.Run("empty set: has empty string representation",
		func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(empty())

			g.Expect(s).To(HaveStringRepr("[]"))
		})
}

func emptySetReturnsEmptySlice(t *testing.T, setBuilder func(elements []string) set.Set[string]) {
	t.Helper()

	t.Run("empty set: returns empty slice", func(t *testing.T) {
		g := NewWithT(t)
		s := setBuilder(empty())

		g.Expect(s).To(HaveEmptyToSlice[string]())
	})
}

func emptySetRemoveDoesNothing(t *testing.T, setBuilder func(elements []string) set.Set[string], mutable bool) {
	t.Helper()

	if mutable {
		t.Run("empty set: remove does nothing",
			func(t *testing.T) {
				g := NewWithT(t)
				s := setBuilder(empty()).(set.MutableSet[string])

				s.Remove("link")

				g.Expect(s).To(HaveLenOfZero())
			})
	}
}

func oneElementSetHasLengthOfOne(t *testing.T, setBuilder func(elements []string) set.Set[string]) {
	t.Helper()

	t.Run("one element set: has length of 1", func(t *testing.T) {
		g := NewWithT(t)
		s := setBuilder(oneElement())

		g.Expect(s).To(HaveLenOf(1))
	})
}

func emptySetPlusOneHasLengthOfOne(t *testing.T, setBuilder func(elements []string) set.Set[string], mutable bool) {
	t.Helper()

	if mutable {
		t.Run("empty set: add: has length of 1", func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(empty()).(set.MutableSet[string])

			s.Add("link")

			g.Expect(s).To(HaveLenOf(1))
		})
	}
}

func oneElementSetContainsPresentElement(t *testing.T, setBuilder func(elements []string) set.Set[string]) {
	t.Helper()

	t.Run("one element set: contains present element",
		func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(oneElement())

			g.Expect(s).To(Contain("link"))
		})
}

func emptySetPlusOneContainsPresentElement(t *testing.T, setBuilder func(elements []string) set.Set[string], mutable bool) {
	t.Helper()

	if mutable {
		t.Run("empty set: add: contains present element",
			func(t *testing.T) {
				g := NewWithT(t)
				s := setBuilder(empty()).(set.MutableSet[string])

				s.Add("link")

				g.Expect(s).To(Contain("link"))
			})
	}
}

func oneElementSetDoesNotContainAbsentElement(t *testing.T, setBuilder func(elements []string) set.Set[string]) {
	t.Helper()

	t.Run("one element set: does not contain absent element",
		func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(oneElement())

			g.Expect(s).ToNot(Contain("zelda"))
		})
}

func emptySetPlusOneDoesNotContainAbsentElement(t *testing.T, setBuilder func(elements []string) set.Set[string], mutable bool) {
	t.Helper()

	if mutable {
		t.Run("empty set: add: does not contain absent element",
			func(t *testing.T) {
				g := NewWithT(t)
				s := setBuilder(empty()).(set.MutableSet[string])

				s.Add("link")

				g.Expect(s).ToNot(Contain("zelda"))
			})
	}
}

func oneElementSetReturnsElementOnIteration(t *testing.T, setBuilder func(elements []string) set.Set[string]) {
	t.Helper()

	t.Run("one element set: returns element on iteration",
		func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(oneElement())

			g.Expect(s).To(BeSetWithForEachThatProduces("link"))
		})
}

func emptySetPlusOneReturnsElementOnIteration(t *testing.T, setBuilder func(elements []string) set.Set[string], mutable bool) {
	t.Helper()

	if mutable {
		t.Run("empty set: add: returns element on iteration",
			func(t *testing.T) {
				g := NewWithT(t)
				s := setBuilder(empty()).(set.MutableSet[string])

				s.Add("link")

				g.Expect(s).To(BeSetWithForEachThatProduces("link"))
			})
	}
}

func oneElementSetHasOneElementStringRepr(t *testing.T, setBuilder func(elements []string) set.Set[string]) {
	t.Helper()

	t.Run("one element set: has one-element string representation",
		func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(oneElement())

			g.Expect(s).To(HaveStringRepr("[link]"))
		})
}

func emptySetPlusOneHasOneElementStringRepr(t *testing.T, setBuilder func(elements []string) set.Set[string], mutable bool) {
	t.Helper()

	if mutable {
		t.Run("empty set: add: has one-element string representation",
			func(t *testing.T) {
				g := NewWithT(t)
				s := setBuilder(empty()).(set.MutableSet[string])

				s.Add("link")

				g.Expect(s).To(HaveStringRepr("[link]"))
			})
	}
}

func oneElementSetReturnsOneElementSlice(t *testing.T, setBuilder func(elements []string) set.Set[string]) {
	t.Helper()

	t.Run("one element set: returns one-element slice",
		func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(oneElement())

			g.Expect(s).To(HaveToSliceThatConsistsOf("link"))
		})
}

func emptySetPlusOneReturnsOneElementSlice(t *testing.T, setBuilder func(elements []string) set.Set[string], mutable bool) {
	t.Helper()

	if mutable {
		t.Run("empty set: add: returns one-element slice",
			func(t *testing.T) {
				g := NewWithT(t)
				s := setBuilder(empty()).(set.MutableSet[string])

				s.Add("link")

				g.Expect(s).To(HaveToSliceThatConsistsOf("link"))
			})
	}
}

func emptySetPlusOneMinusOneDoesNotContainAnything(t *testing.T, setBuilder func(elements []string) set.Set[string], mutable bool) {
	t.Helper()

	if mutable {
		t.Run("empty set: add: remove: does not contain anything",
			func(t *testing.T) {
				g := NewWithT(t)
				s := setBuilder(empty()).(set.MutableSet[string])

				s.Add("link")
				s.Remove("link")

				g.Expect(s).ToNot(Contain("link"))
			})
	}
}

func twoElementSetHasLengthOfTwo(t *testing.T, setBuilder func(elements []string) set.Set[string]) {
	t.Helper()

	t.Run("two element set: has length of 2", func(t *testing.T) {
		g := NewWithT(t)
		s := setBuilder(twoElements())

		g.Expect(s).To(HaveLenOf(2))
	})
}

func emptySetPlusTwoHasLengthOfTwo(t *testing.T, setBuilder func(elements []string) set.Set[string], mutable bool) {
	t.Helper()

	if mutable {
		t.Run("empty set: add x2: has length of 2", func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(empty()).(set.MutableSet[string])

			s.Add("link")
			s.Add("zelda")

			g.Expect(s).To(HaveLenOf(2))
		})
	}
}

func twoElementSetContainsBothElements(t *testing.T, setBuilder func(elements []string) set.Set[string]) {
	t.Helper()

	t.Run("two element set: contains both elements", func(t *testing.T) {
		g := NewWithT(t)
		s := setBuilder(twoElements())

		// TODO: Introduce "ContainAtLeast"
		g.Expect(s).To(And(Contain("link"), Contain("zelda")))
	})
}

func emptySetPlusTwoContainsBothElements(t *testing.T, setBuilder func(elements []string) set.Set[string], mutable bool) {
	t.Helper()

	if mutable {
		t.Run("empty set: add x2: contains both elements", func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(empty()).(set.MutableSet[string])

			s.Add("link")
			s.Add("zelda")

			// TODO: Introduce "ContainAtLeast"
			g.Expect(s).To(And(Contain("link"), Contain("zelda")))
		})
	}
}

func twoElementSetReturnsBothElementsOnIteration(t *testing.T, setBuilder func(elements []string) set.Set[string]) {
	t.Helper()

	t.Run("two element set: returns both elements on iteration",
		func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(twoElements())

			g.Expect(s).To(BeSetWithForEachThatProduces("link", "zelda"))
		})
}

func emptySetPlusTwoReturnsBothElementsOnIteration(t *testing.T, setBuilder func(elements []string) set.Set[string], mutable bool) {
	t.Helper()

	if mutable {
		t.Run("empty set: add x2: returns both elements on iteration",
			func(t *testing.T) {
				g := NewWithT(t)
				s := setBuilder(empty()).(set.MutableSet[string])

				s.Add("link")
				s.Add("zelda")

				g.Expect(s).To(
					BeSetWithForEachThatProduces("link", "zelda"))
			})
	}
}

func twoElementSetHasTwoElementStringRepr(t *testing.T, setBuilder func(elements []string) set.Set[string]) {
	t.Helper()

	t.Run("two element set: has two-element string representation",
		func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(twoElements())

			g.Expect(s).To(
				HaveStringRepr(
					BeElementOf("[link, zelda]", "[zelda, link]")))
		})
}

func emptySetPlusTwoReturnsTwoElementStringRepr(t *testing.T, setBuilder func(elements []string) set.Set[string], mutable bool) {
	t.Helper()

	if mutable {
		t.Run("empty set: add x2: has two-element string representation",
			func(t *testing.T) {
				g := NewWithT(t)
				s := setBuilder(empty()).(set.MutableSet[string])

				s.Add("link")
				s.Add("zelda")

				g.Expect(s).To(
					HaveStringRepr(
						BeElementOf("[link, zelda]", "[zelda, link]")))
			})
	}
}

func twoElementSetReturnsTwoElementSlice(t *testing.T, setBuilder func(elements []string) set.Set[string]) {
	t.Helper()

	t.Run("two element set: returns two-element slice", func(t *testing.T) {
		g := NewWithT(t)
		s := setBuilder(twoElements())

		g.Expect(s).To(HaveToSliceThatConsistsOf("link", "zelda"))
	})
}

func emptySetPlusTwoReturnsTwoElementSlice(t *testing.T, setBuilder func(elements []string) set.Set[string], mutable bool) {
	t.Helper()

	if mutable {
		t.Run("empty set: add x2: returns two-element slice",
			func(t *testing.T) {
				g := NewWithT(t)
				s := setBuilder(empty()).(set.MutableSet[string])

				s.Add("link")
				s.Add("zelda")

				g.Expect(s).To(HaveToSliceThatConsistsOf("link", "zelda"))
			})
	}
}

func emptySetPlusTwoMinusOneHasLengthOfOne(t *testing.T, setBuilder func(elements []string) set.Set[string], mutable bool) {
	t.Helper()

	if mutable {
		t.Run("empty set: add x2: remove x1: has length of 1",
			func(t *testing.T) {
				g := NewWithT(t)
				s := setBuilder(empty()).(set.MutableSet[string])

				s.Add("link")
				s.Add("zelda")
				s.Remove("link")

				g.Expect(s).To(HaveLenOf(1))
			})
	}
}

func emptySetPlusThreeContainsAllThreeElements(t *testing.T, setBuilder func(elements []string) set.Set[string], mutable bool) {
	t.Helper()

	if mutable {
		t.Run("empty set: add x3: contains all three elements", func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(empty()).(set.MutableSet[string])

			s.Add("link")
			s.Add("zelda")
			s.Add("ganondorf")

			// TODO: Introduce "ContainAtLeast"
			g.Expect(s).To(
				And(
					Contain("link"),
					Contain("zelda"),
					Contain("ganondorf")))
		})
	}
}

func threeElementSetContainsAllThreeElements(t *testing.T, setBuilder func(elements []string) set.Set[string]) {
	t.Helper()

	t.Run("three element set: contains all three elements", func(t *testing.T) {
		g := NewWithT(t)
		s := setBuilder(threeElements())

		// TODO: Introduce "ContainAtLeast"
		g.Expect(s).To(
			And(
				Contain("link"),
				Contain("zelda"),
				Contain("ganondorf")))
	})
}

func emptySetPlusThreeHasThreeElementStringRepr(t *testing.T, setBuilder func(elements []string) set.Set[string], mutable bool) {
	t.Helper()

	if mutable {
		t.Run("empty set: add x3: has three-element string representation",
			func(t *testing.T) {
				g := NewWithT(t)
				s := setBuilder(empty()).(set.MutableSet[string])

				s.Add("link")
				s.Add("zelda")
				s.Add("ganondorf")

				g.Expect(s).To(
					HaveStringRepr(
						BeElementOf(
							"[link, zelda, ganondorf]",
							"[link, ganondorf, zelda]",
							"[zelda, link, ganondorf]",
							"[zelda, ganondorf, link]",
							"[ganondorf, link, zelda]",
							"[ganondorf, zelda, link]")))
			})
	}
}

func threeElementSetHasThreeElementStringRepr(t *testing.T, setBuilder func(elements []string) set.Set[string]) {
	t.Helper()

	t.Run("three element set: has three-element string representation",
		func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(threeElements())

			g.Expect(s).To(
				HaveStringRepr(
					BeElementOf(
						"[link, zelda, ganondorf]",
						"[link, ganondorf, zelda]",
						"[zelda, link, ganondorf]",
						"[zelda, ganondorf, link]",
						"[ganondorf, link, zelda]",
						"[ganondorf, zelda, link]")))
		})
}

func setInitializedFromTwoOfSameElementHasLengthOfOne(t *testing.T, setBuilder func(elements []string) set.Set[string]) bool {
	t.Helper()

	return t.Run("set initialized from two of same element: has length of 1",
		func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(twoSameElements())

			g.Expect(s).To(HaveLenOf(1))
		})
}

func emptySetPlusSameElementTwiceHasLengthOfOne(t *testing.T, setBuilder func(elements []string) set.Set[string], mutable bool) {
	t.Helper()

	if mutable {
		t.Run("empty set: add same element x2: has length of 1",
			func(t *testing.T) {
				g := NewWithT(t)
				s := setBuilder(empty()).(set.MutableSet[string])

				s.Add("link")
				s.Add("link")

				g.Expect(s).To(HaveLenOf(1))
			})
	}
}

func empty() []string {
	return make([]string, 0)
}

func oneElement() []string {
	return []string{"link"}
}

func twoElements() []string {
	return []string{"link", "zelda"}
}

func threeElements() []string {
	return []string{"link", "zelda", "ganondorf"}
}

func twoSameElements() []string {
	return []string{"link", "link"}
}