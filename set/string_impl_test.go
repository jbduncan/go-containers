package set_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/jbduncan/go-containers/set"
	. "github.com/onsi/gomega"
)

func TestStringImpl(t *testing.T) {
	type testCase struct {
		name    string
		arg     set.Set[string]
		wantAny []string
	}
	tests := []testCase{
		{
			name:    "empty set",
			arg:     set.Of[string](),
			wantAny: []string{"[]"},
		},
		{
			name:    "one element set",
			arg:     set.Of("link"),
			wantAny: []string{"[link]"},
		},
		{
			name:    "two element set",
			arg:     set.Of("link", "zelda"),
			wantAny: []string{"[link, zelda]", "[zelda, link]"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := set.StringImpl(tt.arg); !slices.Contains(tt.wantAny, got) {
				t.Errorf("StringImpl() = %v, want any of %v", got, tt.wantAny)
			}
		})
	}
}

func FuzzStringImpl(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{0})
	f.Add([]byte{1, 2, 3})
	f.Add([]byte{7, 1})
	f.Add([]byte{255, 123, 4})
	f.Add([]byte{0, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100})
	f.Add(slices.Repeat([]byte{byte(0)}, 10_000))
	f.Add([]byte("x829"))
	f.Add([]byte("0127"))
	f.Add([]byte("78091"))
	f.Add([]byte("0028C17YZ \x10\xda&+\xa8xzyA\x12\xe3a\xfc\xe9\x974c\xffB'8\x90\x90\xd3\x13"))
	f.Add([]byte("\xffy\x000"))

	f.Fuzz(func(t *testing.T, bytes []byte) {
		g := NewWithT(t)
		s := set.Of(bytes...)

		got := set.StringImpl(s)

		g.Expect(got).To(HavePrefix("["))
		g.Expect(got).To(HaveSuffix("]"))
		for elem := range s.All() {
			g.Expect(got).To(ContainSubstring(fmt.Sprintf("%v", elem)))
		}
	})
}
