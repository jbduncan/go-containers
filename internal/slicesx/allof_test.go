package slicesx_test

import (
	"slices"
	"testing"

	"github.com/jbduncan/go-containers/internal/slicesx"
)

func TestAllOf(t *testing.T) {
	t.Parallel()

	type args struct {
		first int
		rest  []int
	}
	type testCase struct {
		name string
		args args
		want []int
	}
	tests := []testCase{
		{
			name: "1 + nil slice",
			args: args{
				first: 1,
			},
			want: []int{1},
		},
		{
			name: "1 + [2]",
			args: args{
				first: 1,
				rest:  []int{2},
			},
			want: []int{1, 2},
		},
		{
			name: "1 + [2, 3]",
			args: args{
				first: 1,
				rest:  []int{2, 3},
			},
			want: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := slicesx.AllOf(tt.args.first, tt.args.rest); !slices.Equal(got, tt.want) {
				t.Errorf("AllOf() = %v, want %v", got, tt.want)
			}
		})
	}
}
