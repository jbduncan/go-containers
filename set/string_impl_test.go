package set_test

import (
	"testing"

	"github.com/jbduncan/go-containers/set"
)

func TestStringImpl(t *testing.T) {
	type args struct {
		s set.Set[string]
	}
	type testCase struct {
		name string
		args args
		want string
	}
	tests := []testCase{
		{
			name: "empty set",
			args: args{
				s: set.New[string](),
			},
			want: "[]",
		},
		{
			name: "one element set",
			args: args{
				s: oneElementSet(),
			},
			want: "[link]",
		},
		{
			name: "two element set",
			args: args{
				s: twoElementSet(),
			},
			want: "[link, zelda]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := set.StringImpl(tt.args.s); got != tt.want {
				t.Errorf("StringImpl() = %v, want %v", got, tt.want)
			}
		})
	}
}
