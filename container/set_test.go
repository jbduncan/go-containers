package container_test

import (
	"github.com/onsi/gomega/types"
	"go-containers/container"
	"reflect"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Set", func() {

	var set container.MutableSet[string]

	BeforeEach(func() {
		set = container.NewSet[string]()
	})

	Describe("when creating a new set", func() {
		It("has a length of 0", func() {
			Expect(set).To(HaveLenOfZero())
		})

		It("does nothing on remove", func() {
			set.Remove("link")

			Expect(set).To(HaveLenOfZero())
		})

		It("returns nothing upon iteration", func() {
			count := 0
			set.ForEach(func(elem string) {
				count++
			})

			Expect(count).To(BeZero())
		})

		It("has an empty list string representation", func() {
			Expect(set.String()).To(Equal("[]"))
		})
	})

	Describe("given a new MutableSet", func() {
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
				var elems []string
				set.ForEach(func(elem string) {
					elems = append(elems, elem)
				})

				Expect(elems).To(ConsistOf("link"))
			})

			It("has a single element list string representation", func() {
				Expect(set.String()).To(Equal("[link]"))
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
				var elems []string
				set.ForEach(func(elem string) {
					elems = append(elems, elem)
				})

				Expect(elems).To(ConsistOf("link", "zelda"))
			})

			It("has a two element list string representation", func() {
				Expect(set.String()).To(BeElementOf("[link, zelda]", "[zelda, link]"))
			})
		})

		Context("when adding two elements and removing one", func() {
			It("has a length of 1", func() {
				set.Add("link")
				set.Add("zelda")
				set.Remove("link")

				Expect(set).To(HaveLenOf(1))
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
				Expect(set.String()).To(
					BeElementOf(
						"[link, zelda, ganondorf]",
						"[link, ganondorf, zelda]",
						"[zelda, link, ganondorf]",
						"[zelda, ganondorf, link]",
						"[ganondorf, link, zelda]",
						"[ganondorf, zelda, link]"))
			})
		})

		Context("when adding the same element twice", func() {
			It("has a length of 1", func() {
				set.Add("link")
				set.Add("link")

				Expect(set).To(HaveLenOf(1))
			})
		})
	})
})

func HaveLenOf(len int) types.GomegaMatcher {
	return And(
		WithTransform(
			func(value any) bool {
				return isSet(value)
			},
			BeTrue()),
		WithTransform(
			func(set any) int {
				return lenOf(set)
			},
			Equal(len)))
}

func HaveLenOfZero() types.GomegaMatcher {
	return HaveLenOf(0)
}

func Contain(elem any) types.GomegaMatcher {
	return WithTransform(
		func(set any) bool {
			return contain(set, elem)
		},
		BeTrue())
}

func isSet(value any) bool {
	return strings.Contains(
		strings.ToLower(reflect.TypeOf(value).String()),
		"container.set")
}

func lenOf(set any) int {
	return int(
		reflect.ValueOf(set).
			MethodByName("Len").
			Call(nil)[0].
			Int())
}

func contain(set any, elem any) bool {
	return reflect.ValueOf(set).
		MethodByName("Contains").
		Call([]reflect.Value{reflect.ValueOf(elem)})[0].
		Bool()
}
