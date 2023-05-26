package settest

import (
	"testing"

	. "github.com/onsi/gomega"
	"go-containers/container/set"
	. "go-containers/internal/matchers"
)

// TODO: Document
func Set(t *testing.T, setBuilder func(elements []string) set.Set[string]) {
	_, mutable := setBuilder(empty()).(set.MutableSet[string])

	t.Run("empty set: has length of 0", func(t *testing.T) {
		g := NewWithT(t)
		s := setBuilder(empty())

		g.Expect(s).To(HaveLenOfZero())
	})

	t.Run("empty set: contains nothing", func(t *testing.T) {
		g := NewWithT(t)
		s := setBuilder(empty())

		g.Expect(s).ToNot(Contain("link"))
	})

	t.Run("empty set: iteration does nothing", func(t *testing.T) {
		g := NewWithT(t)
		s := setBuilder(empty())

		g.Expect(s).To(BeSetWithForEachThatProducesNothing())
	})

	t.Run("empty set: has empty string representation",
		func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(empty())

			g.Expect(s).To(HaveStringRepr("[]"))
		})

	t.Run("empty set: returns empty slice", func(t *testing.T) {
		g := NewWithT(t)
		s := setBuilder(empty())

		g.Expect(s).To(HaveEmptyToSlice[string]())
	})

	if mutable {
		t.Run("empty set: remove does nothing",
			func(t *testing.T) {
				g := NewWithT(t)
				s := setBuilder(empty()).(set.MutableSet[string])

				s.Remove("link")

				g.Expect(s).To(HaveLenOfZero())
			})
	}

	t.Run("one element set: has length of 1", func(t *testing.T) {
		g := NewWithT(t)
		s := setBuilder(oneElement())

		g.Expect(s).To(HaveLenOf(1))
	})

	if mutable {
		t.Run("empty set: add: has length of 1", func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(empty()).(set.MutableSet[string])

			s.Add("link")

			g.Expect(s).To(HaveLenOf(1))
		})
	}

	t.Run("one element set: contains present element",
		func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(oneElement())

			g.Expect(s).To(Contain("link"))
		})

	if mutable {
		t.Run("empty set: add: contains present element",
			func(t *testing.T) {
				g := NewWithT(t)
				s := setBuilder(empty()).(set.MutableSet[string])

				s.Add("link")

				g.Expect(s).To(Contain("link"))
			})
	}

	t.Run("one element set: does not contain absent element",
		func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(oneElement())

			g.Expect(s).ToNot(Contain("zelda"))
		})

	if mutable {
		t.Run("empty set: add: does not contain absent element",
			func(t *testing.T) {
				g := NewWithT(t)
				s := setBuilder(empty()).(set.MutableSet[string])

				s.Add("link")

				g.Expect(s).ToNot(Contain("zelda"))
			})
	}

	t.Run("one element set: returns element on iteration",
		func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(oneElement())

			g.Expect(s).To(BeSetWithForEachThatProduces("link"))
		})

	if mutable {
		t.Run("empty set: add: returns element on iteration",
			func(t *testing.T) {
				g := NewWithT(t)
				s := setBuilder(empty()).(set.MutableSet[string])

				s.Add("link")

				g.Expect(s).To(BeSetWithForEachThatProduces("link"))
			})
	}

	t.Run("one element set: has one-element string representation",
		func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(oneElement())

			g.Expect(s).To(HaveStringRepr("[link]"))
		})

	if mutable {
		t.Run("empty set: add: has one-element string representation",
			func(t *testing.T) {
				g := NewWithT(t)
				s := setBuilder(empty()).(set.MutableSet[string])

				s.Add("link")

				g.Expect(s).To(HaveStringRepr("[link]"))
			})
	}

	t.Run("one element set: returns one-element slice",
		func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(oneElement())

			g.Expect(s).To(HaveToSliceThatConsistsOf("link"))
		})

	if mutable {
		t.Run("empty set: add: returns one-element slice",
			func(t *testing.T) {
				g := NewWithT(t)
				s := setBuilder(empty()).(set.MutableSet[string])

				s.Add("link")

				g.Expect(s).To(HaveToSliceThatConsistsOf("link"))
			})
	}

	if mutable {
		t.Run("empty set: add: remove: does not contain element",
			func(t *testing.T) {
				g := NewWithT(t)
				s := setBuilder(empty()).(set.MutableSet[string])

				s.Add("link")
				s.Remove("link")

				g.Expect(s).ToNot(HaveToSliceThatConsistsOf("link"))
			})
	}

	t.Run("two element set: has length of 2", func(t *testing.T) {
		g := NewWithT(t)
		s := setBuilder(twoElements()).(set.MutableSet[string])

		g.Expect(s).To(HaveLenOf(2))
	})

	if mutable {
		t.Run("empty set: add x2: has length of 2", func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(empty()).(set.MutableSet[string])

			s.Add("link")
			s.Add("zelda")

			g.Expect(s).To(HaveLenOf(2))
		})
	}

	t.Run("two element set: contains both elements", func(t *testing.T) {
		g := NewWithT(t)
		s := setBuilder(twoElements()).(set.MutableSet[string])

		// TODO: Introduce "ContainAll"
		g.Expect(s).To(And(Contain("link"), Contain("zelda")))
	})

	if mutable {
		t.Run("empty set: add x2: contains both elements", func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(empty()).(set.MutableSet[string])

			s.Add("link")
			s.Add("zelda")

			// TODO: Introduce "ContainAll"
			g.Expect(s).To(And(Contain("link"), Contain("zelda")))
		})
	}

	t.Run("two element set: returns both elements on iteration",
		func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(twoElements()).(set.MutableSet[string])

			g.Expect(s).To(BeSetWithForEachThatProduces("link", "zelda"))
		})

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

	t.Run("two element set: has two-element string representation",
		func(t *testing.T) {
			g := NewWithT(t)
			s := setBuilder(twoElements())

			g.Expect(s).To(
				HaveStringRepr(
					BeElementOf("[link, zelda]", "[zelda, link]")))
		})

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

	t.Run("two element set: returns two-element slice", func(t *testing.T) {
		g := NewWithT(t)
		s := setBuilder(twoElements()).(set.MutableSet[string])

		g.Expect(s).To(HaveToSliceThatConsistsOf("link", "zelda"))
	})

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

	if mutable {
		t.Run("empty set: add x2: remove x1: contains just one element",
			func(t *testing.T) {
				g := NewWithT(t)
				s := setBuilder(empty()).(set.MutableSet[string])

				s.Add("link")
				s.Add("zelda")
				s.Remove("link")

				g.Expect(s).To(BeSetThatConsistsOf[string]("zelda"))
			})
	}

	t.Run("three element set: contains all elements", func(t *testing.T) {
		// TODO
	})
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
