package set_test

import (
	"testing"

	"github.com/jbduncan/go-containers/set"
	. "github.com/onsi/gomega"
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
