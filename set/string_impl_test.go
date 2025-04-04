package set_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/jbduncan/go-containers/set"
)

func TestStringImpl(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

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
		s := set.Of(bytes...)

		got := s.String()
		prefixFound := strings.HasPrefix(got, "[")
		if !prefixFound {
			t.Fatalf(
				`got Set.String of %q, want to have prefix "["`,
				got,
			)
		}
		suffixFound := strings.HasSuffix(got, "]")
		if !suffixFound {
			t.Fatalf(
				`got Set.String of %q, want to have suffix "]"`,
				got,
			)
		}
		for elem := range s.All() {
			if !strings.Contains(got, fmt.Sprintf("%v", elem)) {
				t.Fatalf(
					`got Set.String of %q, want to contain %q`,
					got,
					fmt.Sprintf("%v", elem),
				)
			}
		}
	})
}
