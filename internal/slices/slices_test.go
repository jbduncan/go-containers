package slices_test

import (
	"testing"

	"github.com/jbduncan/go-containers/internal/slices"
	. "github.com/onsi/gomega"
)

func TestCartesianProduct(t *testing.T) {
	type args[T any] struct {
		values [][]T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want [][]T
	}
	tests := []testCase[int]{
		{
			name: "empty",
			args: args[int]{
				values: make([][]int, 0),
			},
			want: [][]int{{}},
		},
		{
			name: "one list with one element",
			args: args[int]{
				values: [][]int{{1}},
			},
			want: [][]int{{1}},
		},
		{
			name: "one list with three elements",
			args: args[int]{
				values: [][]int{{1, 2, 3}},
			},
			want: [][]int{{1}, {2}, {3}},
		},
		{
			name: "two lists with one element each",
			args: args[int]{
				values: [][]int{{1}, {2}},
			},
			want: [][]int{{1, 2}},
		},
		{
			name: "two lists with different number elements",
			args: args[int]{
				values: [][]int{{1, 2, 3}, {8, 9}},
			},
			want: [][]int{
				{1, 8},
				{1, 9},
				{2, 8},
				{2, 9},
				{3, 8},
				{3, 9},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			got := slices.CartesianProduct(tt.args.values)

			g.Expect(got).To(Equal(tt.want))
		})
	}
}

func TestCopyToNonNilSlice(t *testing.T) {
	type args struct {
		values []int
	}
	type testCase struct {
		name string
		args args
		want []int
	}
	tests := []testCase{
		{
			name: "nil",
			args: args{
				values: nil,
			},
			want: make([]int, 0),
		},
		{
			name: "empty",
			args: args{
				values: make([]int, 0),
			},
			want: make([]int, 0),
		},
		{
			name: "one element",
			args: args{
				values: []int{1},
			},
			want: []int{1},
		},
		{
			name: "three elements",
			args: args{
				values: []int{1, 2, 3},
			},
			want: []int{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			got := slices.CopyToNonNilSlice(tt.args.values)

			g.Expect(got).ToNot(BeIdenticalTo(tt.args.values))
			g.Expect(got).To(Equal(tt.want))
		})
	}
}

func TestRepeat(t *testing.T) {
	type args struct {
		value int
		times int
	}
	type testCase struct {
		name string
		args args
		want []int
	}
	tests := []testCase{
		{
			name: "zero times",
			args: args{
				value: 1,
				times: 0,
			},
			want: make([]int, 0),
		},
		{
			name: "one time",
			args: args{
				value: 1,
				times: 1,
			},
			want: []int{1},
		},
		{
			name: "three times",
			args: args{
				value: 1,
				times: 3,
			},
			want: []int{1, 1, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			got := slices.Repeat(tt.args.value, tt.args.times)

			g.Expect(got).To(Equal(tt.want))
		})
	}
}

func TestRepeatAndCartesianProduct(t *testing.T) {
	type args struct {
		value []int
		times int
	}
	type testCase struct {
		name string
		args args
		want [][]int
	}
	tests := []testCase{
		{
			name: "two elems, zero times",
			args: args{
				value: []int{1, 2},
				times: 0,
			},
			want: [][]int{{}},
		},
		{
			name: "two elems, one time",
			args: args{
				value: []int{1, 2},
				times: 1,
			},
			want: [][]int{{1}, {2}},
		},
		{
			name: "two elems, two times",
			args: args{
				value: []int{1, 2},
				times: 2,
			},
			want: [][]int{
				{1, 1}, {1, 2}, {2, 1}, {2, 2},
			},
		},
		{
			name: "two elems, three times",
			args: args{
				value: []int{1, 2},
				times: 3,
			},
			want: [][]int{
				{1, 1, 1},
				{1, 1, 2},
				{1, 2, 1},
				{1, 2, 2},
				{2, 1, 1},
				{2, 1, 2},
				{2, 2, 1},
				{2, 2, 2},
			},
		},
		{
			name: "three elems, two times",
			args: args{
				value: []int{1, 2, 3},
				times: 2,
			},
			want: [][]int{
				{1, 1},
				{1, 2},
				{1, 3},
				{2, 1},
				{2, 2},
				{2, 3},
				{3, 1},
				{3, 2},
				{3, 3},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			got := slices.CartesianProduct(slices.Repeat(tt.args.value, tt.args.times))

			g.Expect(got).To(Equal(tt.want))
		})
	}
}
