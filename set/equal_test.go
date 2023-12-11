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

		g.Expect(set.Equal(setOf("link"), nil)).To(BeFalse())
	})

	t.Run("set a: [link]; set b: [zelda]; not equal", func(t *testing.T) {
		g := NewWithT(t)

		g.Expect(set.Equal(setOf("link"), setOf("zelda"))).To(BeFalse())
	})

	t.Run("set a: [link]; set b: [link, zelda]; not equal", func(t *testing.T) {
		g := NewWithT(t)

		g.Expect(set.Equal(setOf("link", "zelda"), setOf("link"))).To(BeFalse())
	})
}
