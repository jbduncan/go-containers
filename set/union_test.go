package set_test

import (
	"testing"

	. "github.com/jbduncan/go-containers/internal/matchers"
	"github.com/jbduncan/go-containers/set"
	"github.com/jbduncan/go-containers/set/settest"
	. "github.com/onsi/gomega"
)

func TestUnion(t *testing.T) {
	settest.Set(t, func(elems []string) set.Set[string] {
		a := set.NewMutable[string]()
		b := set.NewMutable[string]()

		for i, elem := range elems {
			if i%2 == 0 {
				a.Add(elem)
			} else {
				b.Add(elem)
			}
		}

		return set.Union[string](a, b)
	})
}

func TestUnionIsUnmodifiable(t *testing.T) {
	g := NewWithT(t)

	union := set.Union[int](set.Of[int](), set.Of[int]())

	g.Expect(union).To(BeNonMutableSet[int]())
}

func TestUnionIsView(t *testing.T) {
	g := NewWithT(t)
	a := set.NewMutable[int]()
	b := set.NewMutable[int]()
	union := set.Union[int](a, b)

	a.Add(1)

	g.Expect(union).To(Contain(1))
	g.Expect(union).To(HaveLenOf(1))
	g.Expect(union).To(HaveForEachThatConsistsOf[int](1))
	g.Expect(union).To(HaveStringRepr("[1]"))
}
