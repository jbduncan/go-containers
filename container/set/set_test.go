package set_test

import (
	"github.com/onsi/gomega/types"
	"go-containers/container/set"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "go-containers/internal/matchers"
)

var _ = Describe("Set", func() {

	var mutSet set.MutableSet[string]

	BeforeEach(func() {
		mutSet = set.New[string]()
	})

	Describe("given a new set", func() {
		It("has a length of 0", func() {
			Expect(mutSet).To(haveLenOfZero())
		})

		It("does nothing on remove", func() {
			mutSet.Remove("link")

			Expect(mutSet).To(haveLenOfZero())
		})

		It("returns nothing upon iteration", func() {
			Expect(mutSet).To(haveForEachThatProducesNothing())
		})

		It("has an empty list string representation", func() {
			Expect(mutSet).To(HaveStringRepr("[]"))
		})

		Context("when adding one element", func() {
			BeforeEach(func() {
				mutSet.Add("link")
			})

			It("has a length of 1", func() {
				Expect(mutSet).To(haveLenOf(1))
			})

			It("contains the element", func() {
				Expect(mutSet).To(contain("link"))
			})

			It("does not contain any other element", func() {
				Expect(mutSet).ToNot(contain("zelda"))
			})

			It("returns element upon iteration", func() {
				Expect(mutSet).To(haveForEachThatProduces("link"))
			})

			It("has a single element list string representation", func() {
				Expect(mutSet).To(HaveStringRepr("[link]"))
			})
		})

		Context("when adding and removing one element", func() {
			It("has a length of 0", func() {
				mutSet.Add("link")
				mutSet.Remove("link")

				Expect(mutSet).ToNot(contain("link"))
			})
		})

		Context("when adding two elements", func() {
			BeforeEach(func() {
				mutSet.Add("link")
				mutSet.Add("zelda")
			})

			It("has a length of 2", func() {
				Expect(mutSet).To(haveLenOf(2))
			})

			It("contains both elements", func() {
				Expect(mutSet).To(contain("link"))
				Expect(mutSet).To(contain("zelda"))
			})

			It("returns both elements upon iteration", func() {
				Expect(mutSet).To(haveForEachThatProduces("link", "zelda"))
			})

			It("has a two element list string representation", func() {
				Expect(mutSet).To(HaveStringRepr(BeElementOf("[link, zelda]", "[zelda, link]")))
			})

			Context("and removing one", func() {
				It("has a length of 1", func() {
					mutSet.Add("link")
					mutSet.Add("zelda")
					mutSet.Remove("link")

					Expect(mutSet).To(haveLenOf(1))
				})
			})
		})

		Context("when adding three elements", func() {
			BeforeEach(func() {
				mutSet.Add("link")
				mutSet.Add("zelda")
				mutSet.Add("ganondorf")
			})

			It("contains all elements", func() {
				Expect(mutSet).To(contain("link"))
				Expect(mutSet).To(contain("zelda"))
				Expect(mutSet).To(contain("ganondorf"))
			})

			It("has a three element list string representation", func() {
				Expect(mutSet).To(
					HaveStringRepr(
						BeElementOf(
							"[link, zelda, ganondorf]",
							"[link, ganondorf, zelda]",
							"[zelda, link, ganondorf]",
							"[zelda, ganondorf, link]",
							"[ganondorf, link, zelda]",
							"[ganondorf, zelda, link]")))
			})
		})

		Context("when adding the same element twice", func() {
			It("has a length of 1", func() {
				mutSet.Add("link")
				mutSet.Add("link")

				Expect(mutSet).To(haveLenOf(1))
			})
		})

		Context("when wrapping it in an unmodifiable set", func() {
			var unmodSet set.Set[string]

			BeforeEach(func() {
				unmodSet = set.Unmodifiable(mutSet)
			})

			It("has a length of 0", func() {
				Expect(unmodSet).To(haveLenOfZero())
			})

			It("returns nothing upon iteration", func() {
				Expect(unmodSet).To(haveForEachThatProducesNothing())
			})

			It("has an empty list string representation", func() {
				Expect(unmodSet).To(HaveStringRepr("[]"))
			})

			Context("and adding one element afterwards", func() {
				BeforeEach(func() {
					mutSet.Add("link")
				})

				It("has a length of 1", func() {
					Expect(unmodSet).To(haveLenOf(1))
				})

				It("contains the element", func() {
					Expect(unmodSet).To(contain("link"))
				})

				It("returns element upon iteration", func() {
					Expect(unmodSet).To(haveForEachThatProduces("link"))
				})

				It("has a single element list string representation", func() {
					Expect(unmodSet).To(HaveStringRepr("[link]"))
				})
			})

			Context("and adding two elements afterwards", func() {
				It("has a two element list string representation", func() {
					mutSet.Add("link")
					mutSet.Add("zelda")

					Expect(unmodSet).
						To(HaveStringRepr(BeElementOf("[link, zelda]", "[zelda, link]")))
				})
			})
		})
	})
})

func haveLenOf(len int) types.GomegaMatcher {
	return WithTransform(
		func(set set.Set[string]) int {
			return set.Len()
		},
		Equal(len))
}

func haveLenOfZero() types.GomegaMatcher {
	return haveLenOf(0)
}

func contain(elem string) types.GomegaMatcher {
	return WithTransform(
		func(set set.Set[string]) bool {
			return set.Contains(elem)
		},
		BeTrue())
}

func haveForEachThatProduces(first string, others ...string) types.GomegaMatcher {
	all := []string{first}
	all = append(all, others...)

	return WithTransform(forEachResults, ConsistOf(all))
}

func haveForEachThatProducesNothing() types.GomegaMatcher {
	return WithTransform(forEachResults, BeEmpty())
}

func forEachResults(set set.Set[string]) []string {
	var result []string

	set.ForEach(func(elem string) {
		result = append(result, elem)
	})

	return result
}
