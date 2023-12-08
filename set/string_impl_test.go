package set_test

import (
	"testing"

	"github.com/jbduncan/go-containers/set"
	"golang.org/x/exp/slices"
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
			arg:     set.New[string](),
			wantAny: []string{"[]"},
		},
		{
			name:    "one element set",
			arg:     oneElementSet(),
			wantAny: []string{"[link]"},
		},
		{
			name:    "two element set",
			arg:     twoElementSet(),
			wantAny: []string{"[link, zelda]", "[zelda, link]"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := set.StringImpl(tt.arg); !slices.Contains(tt.wantAny, got) {
				t.Errorf("StringImpl() = %v, wantAny any of %v", got, tt.wantAny)
			}
		})
	}
}
