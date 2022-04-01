package container_test

import (
	"github.com/onsi/gomega/types"
	"go-containers/container"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "go-containers/internal/matchers"
)

var _ = Describe("Set", func() {

	var set container.MutableSet[string]

	BeforeEach(func() {
		set = container.NewSet[string]()
	})

	Describe("given a new set", func() {
		It("has a length of 0", func() {
			Expect(set).To(HaveLenOfZero())
		})

		It("does nothing on remove", func() {
			set.Remove("link")

			Expect(set).To(HaveLenOfZero())
		})

		It("returns nothing upon iteration", func() {
			Expect(set).To(HaveForEachThatProducesNothing())
		})

		It("has an empty list string representation", func() {
			Expect(set).To(HaveStringRepr("[]"))
		})

		Context("when adding one element", func() {
			BeforeEach(func() {
				set.Add("link")
			})

			It("has a length of 1", func() {
				Expect(set).To(HaveLenOf(1))
			})

			It("contains the element", func() {
				Expect(set).To(Contain("link"))
			})

			It("does not contain any other element", func() {
				Expect(set).ToNot(Contain("zelda"))
			})

			It("returns element upon iteration", func() {
				Expect(set).To(HaveForEachThatProduces("link"))
			})

			It("has a single element list string representation", func() {
				Expect(set).To(HaveStringRepr("[link]"))
			})
		})

		Context("when adding and removing one element", func() {
			It("has a length of 0", func() {
				set.Add("link")
				set.Remove("link")

				Expect(set).ToNot(Contain("link"))
			})
		})

		Context("when adding two elements", func() {
			BeforeEach(func() {
				set.Add("link")
				set.Add("zelda")
			})

			It("has a length of 2", func() {
				Expect(set).To(HaveLenOf(2))
			})

			It("contains both elements", func() {
				Expect(set).To(Contain("link"))
				Expect(set).To(Contain("zelda"))
			})

			It("returns both elements upon iteration", func() {
				Expect(set).To(HaveForEachThatProduces("link", "zelda"))
			})

			It("has a two element list string representation", func() {
				Expect(set).To(HaveStringRepr(BeElementOf("[link, zelda]", "[zelda, link]")))
			})

			Context("and removing one", func() {
				It("has a length of 1", func() {
					set.Add("link")
					set.Add("zelda")
					set.Remove("link")

					Expect(set).To(HaveLenOf(1))
				})
			})
		})

		Context("when adding three elements", func() {
			BeforeEach(func() {
				set.Add("link")
				set.Add("zelda")
				set.Add("ganondorf")
			})

			It("contains all elements", func() {
				Expect(set).To(Contain("link"))
				Expect(set).To(Contain("zelda"))
				Expect(set).To(Contain("ganondorf"))
			})

			It("has a three element list string representation", func() {
				Expect(set).To(
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
				set.Add("link")
				set.Add("link")

				Expect(set).To(HaveLenOf(1))
			})
		})

		Context("when wrapping it in an unmodifiable set", func() {
			var unmodSet container.Set[string]

			BeforeEach(func() {
				unmodSet = container.UnmodifiableSet(set)
			})

			It("has a length of 0", func() {
				Expect(unmodSet).To(HaveLenOfZero())
			})

			It("returns nothing upon iteration", func() {
				Expect(unmodSet).To(HaveForEachThatProducesNothing())
			})

			It("has an empty list string representation", func() {
				Expect(unmodSet).To(HaveStringRepr("[]"))
			})

			Context("and adding one element afterwards", func() {
				BeforeEach(func() {
					set.Add("link")
				})

				It("has a length of 1", func() {
					Expect(unmodSet).To(HaveLenOf(1))
				})

				It("contains the element", func() {
					Expect(unmodSet).To(Contain("link"))
				})

				It("returns element upon iteration", func() {
					Expect(unmodSet).To(HaveForEachThatProduces("link"))
				})

				It("has a single element list string representation", func() {
					Expect(unmodSet).To(HaveStringRepr("[link]"))
				})
			})

			Context("and adding two elements afterwards", func() {
				It("has a two element list string representation", func() {
					set.Add("link")
					set.Add("zelda")

					Expect(unmodSet).
						To(HaveStringRepr(BeElementOf("[link, zelda]", "[zelda, link]")))
				})
			})
		})
	})
})

func HaveLenOf(len int) types.GomegaMatcher {
	return WithTransform(
		func(set container.Set[string]) int {
			return set.Len()
		},
		Equal(len))
}

func HaveLenOfZero() types.GomegaMatcher {
	return HaveLenOf(0)
}

func Contain(elem string) types.GomegaMatcher {
	return WithTransform(
		func(set container.Set[string]) bool {
			return set.Contains(elem)
		},
		BeTrue())
}

func HaveForEachThatProduces(first string, others ...string) types.GomegaMatcher {
	all := []string{first}
	all = append(all, others...)

	return WithTransform(forEachResults, ConsistOf(all))
}

func HaveForEachThatProducesNothing() types.GomegaMatcher {
	return WithTransform(forEachResults, BeEmpty())
}

func forEachResults(set container.Set[string]) []string {
	var result []string
	set.ForEach(func(elem string) {
		result = append(result, elem)
	})
	return result
}
