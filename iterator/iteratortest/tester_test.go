package iteratortest_test

import (
	"github.com/jbduncan/go-containers/iterator"
	"github.com/jbduncan/go-containers/iterator/iteratortest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Iterator testers", func() {
	Describe("given a new iterator tester", func() {
		Context("with a slice and a known-order iterator with the same elements and order", func() {
			Context("when running the tester", func() {
				It("passes", func() {
					tester := iteratortest.ForIteratorWithKnownOrder(
						"Known-order slice iterator",
						func() iterator.Iterator[int] {
							return sliceIter([]int{1, 2, 3})
						},
						[]int{1, 2, 3})

					err := tester.Test()

					Expect(err).ToNot(HaveOccurred())
				})
			})
		})

		Context("with a slice and a known-order iterator with different orders", func() {
			Context("when running the tester", func() {
				It("fails", func() {
					tester := iteratortest.ForIteratorWithKnownOrder(
						"Slice iterator in different order",
						func() iterator.Iterator[int] {
							return sliceIter([]int{1, 3, 2})
						},
						[]int{1, 2, 3})

					err := tester.Test()

					Expect(err).
						To(MatchError(
							ContainSubstring(
								"iterator 'Slice iterator in different order' misbehaves when " +
									"running operations [Value() Value() Value() Value() Value()]")))
				})
			})
		})

		Context("with a slice and a known-order iterator with different lengths", func() {
			Context("when running the tester", func() {
				It("fails", func() {
					tester := iteratortest.ForIteratorWithKnownOrder(
						"Slice iterator with too many elements",
						func() iterator.Iterator[int] {
							return sliceIter([]int{1, 2, 3, 4})
						},
						[]int{1, 2, 3})

					err := tester.Test()

					Expect(err).
						To(MatchError(
							ContainSubstring(
								"iterator 'Slice iterator with too many elements' misbehaves when " +
									"running operations [Value() Value() Value() Value() Value()]")))
				})
			})
		})

		Context("with a slice and an unknown-order iterator with different orders", func() {
			Context("when running the tester", func() {
				It("passes", func() {
					tester := iteratortest.ForIteratorWithUnknownOrder(
						"Unknown-order slice iterator",
						func() iterator.Iterator[int] {
							return sliceIter([]int{1, 3, 2})
						},
						[]int{1, 2, 3})

					err := tester.Test()

					Expect(err).ToNot(HaveOccurred())
				})
			})
		})
	})
})

func sliceIter(ints []int) iterator.Iterator[int] {
	return &intSliceIterator{ints, 0}
}

type intSliceIterator struct {
	ints  []int
	index int
}

func (i *intSliceIterator) Value() int {
	if !i.Next() {
		panic("no more elements")
	}

	result := i.ints[i.index]
	i.index += 1
	return result
}

func (i *intSliceIterator) Next() bool {
	return i.index < len(i.ints)
}
