package set_test

import (
	"testing"

	slices2 "github.com/jbduncan/go-containers/internal/slices"
	"github.com/jbduncan/go-containers/set"
	. "github.com/onsi/gomega"
	"golang.org/x/exp/slices"
)

func TestToSlice(t *testing.T) {
	type testCase struct {
		name string
		s    set.Set[string]
		want []string
	}
	tests := []testCase{
		{
			name: "empty set",
			s:    set.Of[string](),
			want: make([]string, 0),
		},
		{
			name: "one-element set",
			s:    set.Of("link"),
			want: []string{"link"},
		},
		{
			name: "two-element set",
			s:    set.Of("link", "zelda"),
			want: []string{"link", "zelda"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(set.ToSlice(tt.s)).To(ConsistOf(tt.want))
		})
	}
}

func FuzzToSlice(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{0})
	f.Add([]byte{1, 2, 3})
	f.Add([]byte{7, 1})
	f.Add([]byte{255, 123, 4})
	f.Add([]byte{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100})
	f.Add(slices2.Repeat(byte(0), 10_000))
	f.Add([]byte("0000"))
	f.Add([]byte("0000000000000"))
	f.Add([]byte("0000000000000000"))
	f.Add([]byte("0000000000000000000000000000000000000000000000000000000000000000"))
	f.Add([]byte("000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
		"00000000000000000000000000"))

	f.Fuzz(func(t *testing.T, bytes []byte) {
		g := NewWithT(t)
		s := set.Of(bytes...)

		got := set.ToSlice(s)

		g.Expect(len(got)).To(BeNumerically("<=", len(bytes)))
		g.Expect(len(got)).To(BeNumerically(">=", 0))
		g.Expect(got).To(HaveLen(s.Len()))
		for _, b := range bytes {
			g.Expect(slices.Contains(got, b)).To(BeTrue())
		}
	})
}
