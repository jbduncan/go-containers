package set_test

import (
	"testing"

	"github.com/jbduncan/go-containers/set"
)

func TestEqual(t *testing.T) {
	t.Parallel()

	t.Run("set a: nil; set b: nil; equal", func(t *testing.T) {
		t.Parallel()

		if got := set.Equal[string](nil, nil); !got {
			t.Errorf("set.Equal: got false, want true")
		}
	})

	t.Run("set a: [link]; set b: nil; not equal", func(t *testing.T) {
		t.Parallel()

		if got := set.Equal[string](set.Of("link"), nil); got {
			t.Errorf("set.Equal: got true, want false")
		}
	})

	t.Run("set a: [link]; set b: [zelda]; not equal", func(t *testing.T) {
		t.Parallel()

		if got := set.Equal[string](set.Of("link"), set.Of("zelda")); got {
			t.Errorf("set.Equal: got true, want false")
		}
	})

	t.Run(
		"set a: [link]; set b: [link, zelda]; not equal",
		func(t *testing.T) {
			t.Parallel()

			got := set.Equal[string](set.Of("link"), set.Of("link", "zelda"))
			if got {
				t.Errorf("set.Equal: got true, want false")
			}
		},
	)
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
		a := set.Of(bytesA...)
		b := set.Of(bytesB...)

		if got := set.Equal[byte](a, a); !got {
			t.Errorf("set.Equal(a, a): got false, want true")
		}
		if got := set.Equal[byte](a, b) == set.Equal[byte](b, a); !got {
			t.Errorf(
				"set.Equal(a, b) == set.Equal(b, a): got false, want true",
			)
		}
		// Purposefully checking that two logically identical expressions
		// should always return the same result.
		//nolint:revive,staticcheck
		if got := set.Equal[byte](a, b) == set.Equal[byte](a, b); !got {
			t.Errorf(
				"set.Equal(a, b) == set.Equal(a, b): got false, want true",
			)
		}
	})
}
