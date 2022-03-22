package container_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"go-containers/container"
	"reflect"
	"strings"
)

var _ = Describe("Set", func() {

	var set container.MutableSet[string]

	BeforeEach(func() {
		set = container.NewSet[string]()
	})

	Describe("when calling NewSet", func() {
		It("creates a new MutableSet with a Len of 0", func() {
			Expect(set).To(HaveSetLenOfZero())
		})
	})

	Describe("given a new MutableSet", func() {
		Context("when adding one element", func() {
			BeforeEach(func() {
				set.Add("link")
			})

			It("has a length of 1", func() {
				Expect(set).To(HaveSetLen(1))
			})

			It("contains the element", func() {
				Expect(set).To(SetContain("link"))
			})

			It("does not contain another element", func() {
				Expect(set).ToNot(SetContain("zelda"))
			})
		})

		Context("when adding two elements", func() {
			BeforeEach(func() {
				set.Add("link")
				set.Add("zelda")
			})

			It("has a length of 2", func() {
				Expect(set).To(HaveSetLen(2))
			})

			It("contains both elements", func() {
				Expect(set).To(SetContain("link"))
				Expect(set).To(SetContain("zelda"))
			})
		})

		Context("when adding three elements", func() {
			It("contains all elements", func() {
				set.Add("link")
				set.Add("zelda")
				set.Add("ganondorf")

				Expect(set).To(SetContain("link"))
				Expect(set).To(SetContain("zelda"))
				Expect(set).To(SetContain("ganondorf"))
			})
		})

		Context("when adding the same element twice", func() {
			It("has a length of 1", func() {
				set.Add("link")
				set.Add("link")

				Expect(set).To(HaveSetLen(1))
			})
		})
	})
})

func HaveSetLen(len int) types.GomegaMatcher {

	return WithTransform(
		func(set container.Set[string]) int {
			return set.Len()
		},
		Equal(len))
}

func isSet(value any) bool {
	return strings.Contains(strings.ToLower(reflect.TypeOf(value).String()), "container.set")
}

func callSetLen(set any) int {
	return int(reflect.ValueOf(set).MethodByName("Len").Call(nil)[0].Int())
}

func HaveSetLenOfZero() types.GomegaMatcher {
	return HaveSetLen(0)
}

func SetContain(elem string) types.GomegaMatcher {
	return WithTransform(
		func(set container.Set[string]) bool {
			return set.Contains(elem)
		},
		BeTrue())
}
