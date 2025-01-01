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
		a := set.Of[string]()
		b := set.Of[string]()

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
	a := set.Of[int]()
	b := set.Of[int]()
	union := set.Union[int](a, b)

	a.Add(1)

	g.Expect(union).To(Contain(1))
	g.Expect(union).To(HaveLenOf(1))
	g.Expect(union).To(HaveAllThatConsistsOf[int](1))
	g.Expect(union).To(HaveStringRepr("[1]"))
}

func FuzzUnion(f *testing.F) {
	addUnionFuzzSeedCorpuses(f)

	f.Fuzz(func(t *testing.T, a, b []byte) {
		g := NewWithT(t)
		setA := set.Of(a...)
		setB := set.Of(b...)

		got := set.Union[byte](setA, setB)

		g.Expect(got.Len()).To(BeNumerically(">=", 0))
		g.Expect(got.Len()).To(BeNumerically("<=", len(a)+len(b)))
		for _, elem := range a {
			g.Expect(got).To(Contain(elem))
		}
		for _, elem := range b {
			g.Expect(got).To(Contain(elem))
		}
	})
}

func FuzzUnionHasCommutativeProperty(f *testing.F) {
	addUnionFuzzSeedCorpuses(f)

	f.Fuzz(func(t *testing.T, a, b []byte) {
		g := NewWithT(t)
		setA := set.Of(a...)
		setB := set.Of(b...)

		g.Expect(set.Equal(set.Union[byte](setA, setB), set.Union[byte](setB, setA))).
			To(BeTrue(), "have commutative property")
	})
}

func FuzzUnionHasIdentityProperty(f *testing.F) {
	f.Add([]byte{}, true)
	f.Add([]byte{}, false)
	f.Add([]byte{0}, true)
	f.Add([]byte{0}, false)
	f.Add([]byte{1, 2, 3, 4, 5}, true)
	f.Add([]byte{1, 2, 3, 4, 5}, false)

	f.Fuzz(func(t *testing.T, bytes []byte, identityFirst bool) {
		g := NewWithT(t)
		s := set.Of(bytes...)

		var union set.UnionSet[byte]
		if identityFirst {
			union = set.Union[byte](set.Of[byte](), s)
		} else {
			union = set.Union[byte](s, set.Of[byte]())
		}

		g.Expect(set.Equal(union, s)).
			To(BeTrue(), "have identity property")
	})
}

func FuzzUnionHasIdempotentProperty(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{0})
	f.Add([]byte{1, 2, 3, 4, 5})

	f.Fuzz(func(t *testing.T, bytes []byte) {
		g := NewWithT(t)
		s := set.Of(bytes...)

		union := set.Union[byte](s, s)
		g.Expect(set.Equal(union, s)).
			To(BeTrue(), "have idempotent property")
	})
}

func addUnionFuzzSeedCorpuses(f *testing.F) {
	f.Helper()

	f.Add([]byte{}, []byte{})
	f.Add([]byte{1}, []byte{})
	f.Add([]byte{}, []byte{2})
	f.Add([]byte{3}, []byte{4})
	f.Add([]byte{5, 6}, []byte{7, 8, 9})
	f.Add([]byte{10, 20, 30, 50, 60, 70}, []byte{80, 90, 100})
	f.Add([]byte("0"), []byte("00"))
	f.Add([]byte("0"), []byte("127"))
	f.Add([]byte("01"), []byte("2"))
	f.Add([]byte("012789ABCXYZab"), []byte("cxy"))
	f.Add([]byte("089ABXYZ17"), []byte("2C"))
	f.Add([]byte("00000000000000"), []byte("0"))
	f.Add([]byte("0000"), []byte("0000000000000000000000000000000"))
	f.Add(
		[]byte("79ABX1YZayz \"c#$%&'2()*,\x0ft\b!x\xcbG\xe0\x0e\a-Cb\xa7\xc0\xf6\xc2\xd6\xe4u\x84@\v\x87\xcc.s\xdf\f"+
			"]}+\xc3\xe5\xda\xf9N\x8c"),
		[]byte("08"))
}
