package settest

import (
	"testing"

	//lint:ignore ST1001 dot importing gomega matchers is best practice and
	// this package is used by test code only
	. "github.com/jbduncan/go-containers/internal/matchers"
	"github.com/jbduncan/go-containers/set"
	//lint:ignore ST1001 dot importing gomega matchers is best practice and
	// this package is used by test code only
	. "github.com/onsi/gomega"
)

// TestingT is an interface for the parts of *testing.T that settest.Set needs
// to run. Whenever you see this interface being used, pass in an instance of
// *testing.T.
type TestingT interface {
	Helper()
	Run(name string, f func(t *testing.T)) bool
}

// TODO: Document
func Set(t TestingT, setBuilder func(elements []string) set.Set[string]) {
	_, mutable := setBuilder(empty()).(set.MutableSet[string])

	emptySetHasLengthOfZero(t, setBuilder)

	emptySetContainsNothing(t, setBuilder)

	emptySetIterationDoesNothing(t, setBuilder)

	emptySetHasEmptyStringRepr(t, setBuilder)

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

	emptySetPlusOneMinusOneDoesNotContainAnything(t, setBuilder, mutable)

	twoElementSetHasLengthOfTwo(t, setBuilder)
	emptySetPlusTwoHasLengthOfTwo(t, setBuilder, mutable)

	twoElementSetContainsBothElements(t, setBuilder)
	emptySetPlusTwoContainsBothElements(t, setBuilder, mutable)

	twoElementSetReturnsBothElementsOnIteration(t, setBuilder)
	emptySetPlusTwoReturnsBothElementsOnIteration(t, setBuilder, mutable)

	twoElementSetHasTwoElementStringRepr(t, setBuilder)
	emptySetPlusTwoReturnsTwoElementStringRepr(t, setBuilder, mutable)

	emptySetPlusTwoMinusOneHasLengthOfOne(t, setBuilder, mutable)

	threeElementSetContainsAllThreeElements(t, setBuilder)
	emptySetPlusThreeContainsAllThreeElements(t, setBuilder, mutable)

	threeElementSetHasThreeElementStringRepr(t, setBuilder)
	emptySetPlusThreeHasThreeElementStringRepr(t, setBuilder, mutable)

	setInitializedFromTwoOfSameElementHasLengthOfOne(t, setBuilder)
	emptySetPlusSameElementTwiceHasLengthOfOne(t, setBuilder, mutable)

	emptySetPlusOneReturnsTrue(t, setBuilder, mutable)

	emptySetPlusSameElementTwiceReturnsFalse(t, setBuilder, mutable)

	emptySetPlusSameElementTwiceThenDifferentOnceReturnsTrue(
		t, setBuilder, mutable)

	emptySetMinusOneReturnsFalse(t, setBuilder, mutable)

	emptySetPlusOneMinusSameElementReturnsTrue(t, setBuilder, mutable)

	emptySetPlusOneMinusSameElementTwiceReturnsFalse(t, setBuilder, mutable)
}

func emptySetHasLengthOfZero(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
) {
	t.Helper()

	t.Run("empty set: has length of 0", func(t *testing.T) {
		g := NewWithT(t)
		s := setBuilder(empty())

		g.Expect(s).To(HaveLenOfZero())
	})
}

func emptySetContainsNothing(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
) {
	t.Helper()

	t.Run("empty set: contains nothing", func(t *testing.T) {
		g := NewWithT(t)
		s := setBuilder(empty())

		g.Expect(s).ToNot(Contain("link"))
	})
}

func emptySetIterationDoesNothing(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
) {
	t.Helper()

	t.Run("empty set: iteration does nothing", func(t *testing.T) {
		g := NewWithT(t)
		s := setBuilder(empty())

		g.Expect(s).To(HaveForEachThatEmitsNothing[string]())
	})
}

func emptySetHasEmptyStringRepr(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
) {
	t.Helper()

	t.Run("empty set: has empty string representation",
		func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(empty())

			g.Expect(s).To(HaveStringRepr("[]"))
		})
}

func emptySetRemoveDoesNothing(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
	mutable bool,
) {
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

func oneElementSetHasLengthOfOne(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
) {
	t.Helper()

	t.Run("one element set: has length of 1", func(t *testing.T) {
		g := NewWithT(t)
		s := setBuilder(oneElement())

		g.Expect(s).To(HaveLenOf(1))
	})
}

func emptySetPlusOneHasLengthOfOne(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
	mutable bool,
) {
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

func oneElementSetContainsPresentElement(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
) {
	t.Helper()

	t.Run("one element set: contains present element",
		func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(oneElement())

			g.Expect(s).To(Contain("link"))
		})
}

func emptySetPlusOneContainsPresentElement(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
	mutable bool,
) {
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

func oneElementSetDoesNotContainAbsentElement(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
) {
	t.Helper()

	t.Run("one element set: does not contain absent element",
		func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(oneElement())

			g.Expect(s).ToNot(Contain("zelda"))
		})
}

func emptySetPlusOneDoesNotContainAbsentElement(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
	mutable bool,
) {
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

func oneElementSetReturnsElementOnIteration(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
) {
	t.Helper()

	t.Run("one element set: returns element on iteration",
		func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(oneElement())

			g.Expect(s).To(HaveForEachThatConsistsOf[string]("link"))
		})
}

func emptySetPlusOneReturnsElementOnIteration(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
	mutable bool,
) {
	t.Helper()

	if mutable {
		t.Run("empty set: add: returns element on iteration",
			func(t *testing.T) {
				g := NewWithT(t)
				s := setBuilder(empty()).(set.MutableSet[string])

				s.Add("link")

				g.Expect(s).To(HaveForEachThatConsistsOf[string]("link"))
			})
	}
}

func oneElementSetHasOneElementStringRepr(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
) {
	t.Helper()

	t.Run("one element set: has one-element string representation",
		func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(oneElement())

			g.Expect(s).To(HaveStringRepr("[link]"))
		})
}

func emptySetPlusOneHasOneElementStringRepr(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
	mutable bool,
) {
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

func emptySetPlusOneMinusOneDoesNotContainAnything(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
	mutable bool,
) {
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

func twoElementSetHasLengthOfTwo(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
) {
	t.Helper()

	t.Run("two element set: has length of 2", func(t *testing.T) {
		g := NewWithT(t)
		s := setBuilder(twoElements())

		g.Expect(s).To(HaveLenOf(2))
	})
}

func emptySetPlusTwoHasLengthOfTwo(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
	mutable bool,
) {
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

func twoElementSetContainsBothElements(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
) {
	t.Helper()

	t.Run("two element set: contains both elements", func(t *testing.T) {
		g := NewWithT(t)
		s := setBuilder(twoElements())

		g.Expect(s).To(ContainAtLeast("link", "zelda"))
	})
}

func emptySetPlusTwoContainsBothElements(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
	mutable bool,
) {
	t.Helper()

	if mutable {
		t.Run("empty set: add x2: contains both elements", func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(empty()).(set.MutableSet[string])

			s.Add("link")
			s.Add("zelda")

			g.Expect(s).To(ContainAtLeast("link", "zelda"))
		})
	}
}

func twoElementSetReturnsBothElementsOnIteration(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
) {
	t.Helper()

	t.Run("two element set: returns both elements on iteration",
		func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(twoElements())

			g.Expect(s).To(HaveForEachThatConsistsOf[string]("link", "zelda"))
		})
}

func emptySetPlusTwoReturnsBothElementsOnIteration(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
	mutable bool,
) {
	t.Helper()

	if mutable {
		t.Run("empty set: add x2: returns both elements on iteration",
			func(t *testing.T) {
				g := NewWithT(t)
				s := setBuilder(empty()).(set.MutableSet[string])

				s.Add("link")
				s.Add("zelda")

				g.Expect(s).To(
					HaveForEachThatConsistsOf[string]("link", "zelda"))
			})
	}
}

func twoElementSetHasTwoElementStringRepr(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
) {
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

func emptySetPlusTwoReturnsTwoElementStringRepr(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
	mutable bool,
) {
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

func emptySetPlusTwoMinusOneHasLengthOfOne(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
	mutable bool,
) {
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

func emptySetPlusThreeContainsAllThreeElements(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
	mutable bool,
) {
	t.Helper()

	if mutable {
		t.Run("empty set: add x3: contains all three elements", func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(empty()).(set.MutableSet[string])

			s.Add("link")
			s.Add("zelda")
			s.Add("ganondorf")

			g.Expect(s).To(
				ContainAtLeast("link", "zelda", "ganondorf"))
		})
	}
}

func threeElementSetContainsAllThreeElements(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
) {
	t.Helper()

	t.Run("three element set: contains all three elements", func(t *testing.T) {
		g := NewWithT(t)
		s := setBuilder(threeElements())

		g.Expect(s).To(
			ContainAtLeast("link", "zelda", "ganondorf"))
	})
}

func emptySetPlusThreeHasThreeElementStringRepr(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
	mutable bool,
) {
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

func threeElementSetHasThreeElementStringRepr(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
) {
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

func setInitializedFromTwoOfSameElementHasLengthOfOne(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
) {
	t.Helper()

	t.Run("set initialized from two of same element: has length of 1",
		func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(twoSameElements())

			g.Expect(s).To(HaveLenOf(1))
		})
}

func emptySetPlusSameElementTwiceHasLengthOfOne(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
	mutable bool,
) {
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

func emptySetPlusOneReturnsTrue(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
	mutable bool,
) {
	t.Helper()

	if mutable {
		t.Run("empty set: add: returns true",
			func(t *testing.T) {
				g := NewWithT(t)
				s := setBuilder(empty()).(set.MutableSet[string])

				result := s.Add("link")

				g.Expect(result).To(BeTrue())
			})
	}
}

func emptySetPlusSameElementTwiceReturnsFalse(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
	mutable bool,
) {
	t.Helper()

	if mutable {
		t.Run("empty set: add same element x2: returns true",
			func(t *testing.T) {
				g := NewWithT(t)
				s := setBuilder(empty()).(set.MutableSet[string])

				s.Add("link")
				result := s.Add("link")

				g.Expect(result).To(BeFalse())
			})
	}
}

func emptySetPlusSameElementTwiceThenDifferentOnceReturnsTrue(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
	mutable bool,
) {
	t.Helper()

	if mutable {
		t.Run("empty set: add same element x2: add different element: returns true",
			func(t *testing.T) {
				g := NewWithT(t)
				s := setBuilder(empty()).(set.MutableSet[string])

				s.Add("link")
				s.Add("link")
				result := s.Add("zelda")

				g.Expect(result).To(BeTrue())
			})
	}
}

func emptySetMinusOneReturnsFalse(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
	mutable bool,
) {
	t.Helper()

	if mutable {
		t.Run("empty set: remove: returns false",
			func(t *testing.T) {
				g := NewWithT(t)
				s := setBuilder(empty()).(set.MutableSet[string])

				result := s.Remove("link")

				g.Expect(result).To(BeFalse())
			})
	}
}

func emptySetPlusOneMinusSameElementReturnsTrue(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
	mutable bool,
) {
	t.Helper()

	if mutable {
		t.Run("empty set: add: remove same element: returns true",
			func(t *testing.T) {
				g := NewWithT(t)
				s := setBuilder(empty()).(set.MutableSet[string])

				s.Add("link")
				result := s.Remove("link")

				g.Expect(result).To(BeTrue())
			})
	}
}

func emptySetPlusOneMinusSameElementTwiceReturnsFalse(
	t TestingT,
	setBuilder func(elements []string) set.Set[string],
	mutable bool,
) {
	t.Helper()

	if mutable {
		t.Run("empty set: add: remove same element x2: returns false",
			func(t *testing.T) {
				g := NewWithT(t)
				s := setBuilder(empty()).(set.MutableSet[string])

				s.Add("link")
				s.Remove("link")
				result := s.Remove("link")

				g.Expect(result).To(BeFalse())
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
