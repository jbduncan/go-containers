package set_test

import (
	"testing"

	internalsettest "github.com/jbduncan/go-containers/internal/settest"
	"github.com/jbduncan/go-containers/set"
	"github.com/jbduncan/go-containers/set/settest"
)

func TestUnion(t *testing.T) {
	t.Parallel()

	settest.TestReadOnly(t, func(elements []int) settest.Set[int] {
		a := set.Of[int]()
		b := set.Of[int]()

		for i, element := range elements {
			if i%2 == 0 {
				a.Add(element)
			} else {
				b.Add(element)
			}
		}

		return set.Union[int](a, b)
	})

	t.Run("union is unmodifiable", func(t *testing.T) {
		t.Parallel()

		union := set.Union[int](set.Of[int](), set.Of[int]())

		internalsettest.IsMutable(t, "set.Union", union)
	})

	t.Run("union is view", func(t *testing.T) {
		t.Parallel()

		a := set.Of[int]()
		b := set.Of[int]()
		union := set.Union[int](a, b)

		a.Add(1)

		internalsettest.Len(t, "set.Union", union, 1)
		internalsettest.All(t, "set.Union", union, []int{1})
		internalsettest.Contains(t, "set.Union", union, []int{1})
		internalsettest.String(t, "set.Union", union, []int{1})
	})
}

func FuzzUnion(f *testing.F) {
	addUnionFuzzSeedCorpuses(f)

	f.Fuzz(func(t *testing.T, a, b []byte) {
		setA := set.Of(a...)
		setB := set.Of(b...)

		union := set.Union[byte](setA, setB)

		if got := 0 <= union.Len() && union.Len() <= len(a)+len(b); !got {
			t.Errorf(
				"set.Union: got Set.Len of %d, want in range [0-%d]",
				union.Len(),
				len(a)+len(b),
			)
		}
		internalsettest.Contains(t, "set.Union", union, a)
		internalsettest.Contains(t, "set.Union", union, b)
	})
}

func FuzzUnionHasCommutativeProperty(f *testing.F) {
	addUnionFuzzSeedCorpuses(f)

	f.Fuzz(func(t *testing.T, a, b []byte) {
		setA := set.Of(a...)
		setB := set.Of(b...)

		if got := set.Equal[byte](
			set.Union[byte](setA, setB),
			set.Union[byte](setB, setA),
		); !got {
			t.Error("set.Union: have commutative property: " +
				"got false, want true")
		}
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
		s := set.Of(bytes...)

		var union set.UnionSet[byte]
		if identityFirst {
			union = set.Union[byte](set.Of[byte](), s)
		} else {
			union = set.Union[byte](s, set.Of[byte]())
		}

		if got := set.Equal[byte](union, s); !got {
			t.Error("set.Union: have identity property: got false, want true")
		}
	})
}

func FuzzUnionHasIdempotentProperty(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{0})
	f.Add([]byte{1, 2, 3, 4, 5})

	f.Fuzz(func(t *testing.T, bytes []byte) {
		s := set.Of(bytes...)

		union := set.Union[byte](s, s)
		if got := set.Equal[byte](union, s); !got {
			t.Error("set.Union: have idempotent property: " +
				"got false, want true")
		}
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
