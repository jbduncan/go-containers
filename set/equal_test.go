package set_test

import (
	"testing"

	"github.com/jbduncan/go-containers/set"
	. "github.com/onsi/gomega"
)

func TestEqual(t *testing.T) {
	t.Run("set a: nil; set b: nil; equal", func(t *testing.T) {
		g := NewWithT(t)

		g.Expect(set.Equal[string](nil, nil)).To(BeTrue())
	})

	t.Run("set a: [link]; set b: nil; not equal", func(t *testing.T) {
		g := NewWithT(t)

		g.Expect(set.Equal(set.Of("link"), nil)).To(BeFalse())
	})

	t.Run("set a: [link]; set b: [zelda]; not equal", func(t *testing.T) {
		g := NewWithT(t)

		g.Expect(set.Equal(set.Of("link"), set.Of("zelda"))).To(BeFalse())
	})

	t.Run("set a: [link]; set b: [link, zelda]; not equal", func(t *testing.T) {
		g := NewWithT(t)

		g.Expect(set.Equal(set.Of("link", "zelda"), set.Of("link"))).To(BeFalse())
	})
}

func FuzzEquals(f *testing.F) {
	f.Add([]byte{}, []byte{})
	f.Add([]byte{1}, []byte{})
	f.Add([]byte{}, []byte{2})
	f.Add([]byte{3}, []byte{4})
	f.Add([]byte{5, 6}, []byte{7, 8, 9})
	f.Add([]byte{10, 20, 30, 50, 60, 70}, []byte{80, 90, 100})
	f.Add([]byte("0"), []byte("00"))

	f.Fuzz(func(t *testing.T, bytesA []byte, bytesB []byte) {
		g := NewWithT(t)
		a := set.Of(bytesA...)
		b := set.Of(bytesB...)

		g.Expect(set.Equal(a, a)).To(BeTrue())
		g.Expect(set.Equal(a, b)).To(Equal(set.Equal(b, a)))
		g.Expect(set.Equal(a, b)).To(Equal(set.Equal(a, b)))
	})
}
