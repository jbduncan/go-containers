package orderagnostic_test

import (
	"testing"

	"github.com/jbduncan/go-containers/internal/orderagnostic"
)

func TestDiff(t *testing.T) {
	type args struct {
		got  []int
		want []int
	}
	type testCase struct {
		name      string
		args      args
		wantEmpty bool
	}
	tests := []testCase{
		{
			name:      "nil slices",
			args:      args{},
			wantEmpty: true,
		},
		{
			name: "empty slices",
			args: args{
				got:  []int{},
				want: []int{},
			},
			wantEmpty: true,
		},
		{
			name: "identical one-element slices",
			args: args{
				got:  []int{1},
				want: []int{1},
			},
			wantEmpty: true,
		},
		{
			name: "different one-element slices",
			args: args{
				got:  []int{1},
				want: []int{2},
			},
			wantEmpty: false,
		},
		{
			name: "slices with the same elements in a different order",
			args: args{
				got:  []int{1, 2},
				want: []int{2, 1},
			},
			wantEmpty: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := orderagnostic.Diff(
				tt.args.got,
				tt.args.want,
			); len(got) == 0 != tt.wantEmpty {
				diffKind := "empty string diff"
				if !tt.wantEmpty {
					diffKind = "non-" + diffKind
				}
				t.Errorf("Diff(): got %v, want %v", got, diffKind)
			}
		})
	}
}
