package orderagnostic_test

import (
	"github.com/jbduncan/go-containers/internal/orderagnostic"
	"testing"
)

func TestSlicesEqual(t *testing.T) {
	type args struct {
		a []int
		b []int
	}
	type testCase struct {
		name string
		args args
		want bool
	}
	tests := []testCase{
		{
			name: "nil slices",
			args: args{},
			want: true,
		},
		{
			name: "empty slices",
			args: args{
				a: []int{},
				b: []int{},
			},
			want: true,
		},
		{
			name: "identical one-element slices",
			args: args{
				a: []int{1},
				b: []int{1},
			},
			want: true,
		},
		{
			name: "different one-element slices",
			args: args{
				a: []int{1},
				b: []int{2},
			},
			want: false,
		},
		{
			name: "slices with the same elements in a different order",
			args: args{
				a: []int{1, 2},
				b: []int{2, 1},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := orderagnostic.SlicesEqual(
				tt.args.a,
				tt.args.b,
			); got != tt.want {
				t.Errorf("SlicesEqual(): got %v, want %v", got, tt.want)
			}
		})
	}
}
