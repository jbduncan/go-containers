package set_test

import (
	"github.com/onsi/gomega/types"
	"go-containers/container/set"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "go-containers/internal/matchers"
)

var _ = Describe("Sets", func() {

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

		It("returns nothing on iteration", func() {
			Expect(mutSet).To(haveForEachThatProducesNothing())
		})

		It("has an empty list string representation", func() {
			Expect(mutSet).To(HaveStringRepr("[]"))
		})

		Context("when returning a slice representation", func() {
			It("returns an empty slice", func() {
				Expect(mutSet).To(HaveEmptyToSlice())
			})
		})

		Context("when adding one element", func() {
			BeforeEach(func() {
				mutSet.Add("link")
			})

			It("has a length of 1", func() {
				Expect(mutSet).To(HaveLenOf(1))
			})

			It("contains the element", func() {
				Expect(mutSet).To(Contain("link"))
			})

			It("does not contain any other element", func() {
				Expect(mutSet).ToNot(Contain("zelda"))
			})

			It("returns element on iteration", func() {
				Expect(mutSet).To(haveForEachThatProduces("link"))
			})

			Context("when returning a slice representation", func() {
				It("returns a single element slice", func() {
					Expect(mutSet).To(HaveToSliceThatConsistsOf("link"))
				})
			})

			It("has a single element list string representation", func() {
				Expect(mutSet).To(HaveStringRepr("[link]"))
			})
		})

		Context("when adding and removing one element", func() {
			It("has a length of 0", func() {
				mutSet.Add("link")
				mutSet.Remove("link")

				Expect(mutSet).ToNot(Contain("link"))
			})
		})

		Context("when adding two elements", func() {
			BeforeEach(func() {
				mutSet.Add("link")
				mutSet.Add("zelda")
			})

			It("has a length of 2", func() {
				Expect(mutSet).To(HaveLenOf(2))
			})

			It("contains both elements", func() {
				Expect(mutSet).To(Contain("link"))
				Expect(mutSet).To(Contain("zelda"))
			})

			It("returns both elements upon iteration", func() {
				Expect(mutSet).To(haveForEachThatProduces("link", "zelda"))
			})

			Context("when returning a slice representation", func() {
				It("returns a two element slice", func() {
					Expect(mutSet).To(HaveToSliceThatConsistsOf("link", "zelda"))
				})
			})

			It("has a two element list string representation", func() {
				Expect(mutSet).To(HaveStringRepr(BeElementOf("[link, zelda]", "[zelda, link]")))
			})

			Context("and removing one", func() {
				It("has a length of 1", func() {
					mutSet.Add("link")
					mutSet.Add("zelda")
					mutSet.Remove("link")

					Expect(mutSet).To(HaveLenOf(1))
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
				Expect(mutSet).To(Contain("link"))
				Expect(mutSet).To(Contain("zelda"))
				Expect(mutSet).To(Contain("ganondorf"))
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

				Expect(mutSet).To(HaveLenOf(1))
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

			It("returns nothing on iteration", func() {
				Expect(unmodSet).To(haveForEachThatProducesNothing())
			})

			Context("when returning a slice representation", func() {
				It("returns an empty slice", func() {
					Expect(mutSet).To(HaveEmptyToSlice())
				})
			})

			It("has an empty list string representation", func() {
				Expect(unmodSet).To(HaveStringRepr("[]"))
			})

			Context("and adding one element afterwards", func() {
				BeforeEach(func() {
					mutSet.Add("link")
				})

				It("has a length of 1", func() {
					Expect(unmodSet).To(HaveLenOf(1))
				})

				It("contains the element", func() {
					Expect(unmodSet).To(Contain("link"))
				})

				It("returns element on iteration", func() {
					Expect(unmodSet).To(haveForEachThatProduces("link"))
				})

				Context("when returning a slice representation", func() {
					It("returns a single element slice", func() {
						Expect(mutSet).To(HaveToSliceThatConsistsOf("link"))
					})
				})

				It("has a single element list string representation", func() {
					Expect(unmodSet).To(HaveStringRepr("[link]"))
				})
			})

			Context("and adding two elements afterwards", func() {
				BeforeEach(func() {
					mutSet.Add("link")
					mutSet.Add("zelda")
				})

				Context("when returning a slice representation", func() {
					It("returns a two element slice", func() {
						Expect(mutSet).To(HaveToSliceThatConsistsOf("link", "zelda"))
					})
				})

				It("has a two element list string representation", func() {
					Expect(unmodSet).
						To(HaveStringRepr(BeElementOf("[link, zelda]", "[zelda, link]")))
				})
			})
		})
	})
})

func haveLenOfZero() types.GomegaMatcher {
	return HaveLenOf(0)
}

func haveForEachThatProduces(first string, others ...string) types.GomegaMatcher {
	all := []string{first}
	all = append(all, others...)

	return WithTransform(ForEachToSlice[string], ConsistOf(all))
}

func haveForEachThatProducesNothing() types.GomegaMatcher {
	return WithTransform(ForEachToSlice[string], BeEmpty())
}
