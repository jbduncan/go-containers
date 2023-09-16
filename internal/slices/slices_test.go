package slices

import (
	"testing"

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
				values: oneListWithOneElement(),
			},
			want: oneListWithOneElementCartesianProduct(),
		},
		{
			name: "one list with three elements",
			args: args[int]{
				values: oneListWithThreeElements(),
			},
			want: oneListWithThreeElementsCartesianProduct(),
		},
		{
			name: "two lists with one element each",
			args: args[int]{
				values: twoListsWithOneElementEach(),
			},
			want: twoListsWithOneElementEachCartesianProduct(),
		},
		{
			name: "two lists with different number elements",
			args: args[int]{
				values: twoListsWithDifferentNumberElements(),
			},
			want: twoListsWithDifferentNumberElementsCartesianProduct(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			got := CartesianProduct(tt.args.values)

			g.Expect(got).To(Equal(tt.want))
		})
	}
}

func TestCopyToNonNilSlice(t *testing.T) {
	type args[T any] struct {
		values []T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "nil",
			args: args[int]{
				values: nil,
			},
			want: make([]int, 0),
		},
		{
			name: "empty",
			args: args[int]{
				values: make([]int, 0),
			},
			want: make([]int, 0),
		},
		{
			name: "one element",
			args: args[int]{
				values: oneElement(),
			},
			want: oneElement(),
		},
		{
			name: "three elements",
			args: args[int]{threeElements()},
			want: threeElements(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			got := CopyToNonNilSlice(tt.args.values)

			g.Expect(got).ToNot(BeIdenticalTo(tt.args.values))
			g.Expect(got).To(Equal(tt.want))
		})
	}
}

func TestRepeat(t *testing.T) {
	type args[T any] struct {
		value T
		times int
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "zero times",
			args: args[int]{
				value: 1,
				times: 0,
			},
			want: make([]int, 0),
		},
		{
			name: "one time",
			args: args[int]{
				value: 1,
				times: 1,
			},
			want: []int{1},
		},
		{
			name: "three times",
			args: args[int]{
				value: 1,
				times: 3,
			},
			want: []int{1, 1, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewWithT(t)

			got := Repeat(tt.args.value, tt.args.times)

			g.Expect(got).To(Equal(tt.want))
		})
	}
}

func TestRepeatAndCartesianProduct(t *testing.T) {
	type args[T any] struct {
		value T
		times int
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[[]int]{
		{
			name: "two elems, zero times",
			args: args[[]int]{
				value: []int{1, 2},
				times: 0,
			},
			want: [][]int{{}},
		},
		{
			name: "two elems, one time",
			args: args[[]int]{
				value: []int{1, 2},
				times: 1,
			},
			want: [][]int{{1}, {2}},
		},
		{
			name: "two elems, two times",
			args: args[[]int]{
				value: []int{1, 2},
				times: 2,
			},
			want: [][]int{
				{1, 1}, {1, 2}, {2, 1}, {2, 2},
			},
		},
		{
			name: "two elems, three times",
			args: args[[]int]{
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
			args: args[[]int]{
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

			got := CartesianProduct(Repeat(tt.args.value, tt.args.times))

			g.Expect(got).To(Equal(tt.want))
		})
	}
}

func oneListWithOneElement() [][]int {
	return [][]int{
		{
			1,
		},
	}
}

func oneListWithOneElementCartesianProduct() [][]int {
	return [][]int{
		{
			1,
		},
	}
}

func oneListWithThreeElements() [][]int {
	return [][]int{
		{
			1, 2, 3,
		},
	}
}

func oneListWithThreeElementsCartesianProduct() [][]int {
	return [][]int{
		{
			1,
		},
		{
			2,
		},
		{
			3,
		},
	}
}

func twoListsWithDifferentNumberElements() [][]int {
	return [][]int{
		{
			1, 2, 3,
		},
		{
			8, 9,
		},
	}
}

func twoListsWithDifferentNumberElementsCartesianProduct() [][]int {
	return [][]int{
		{
			1, 8,
		},
		{
			1, 9,
		},
		{
			2, 8,
		},
		{
			2, 9,
		},
		{
			3, 8,
		},
		{
			3, 9,
		},
	}
}

func twoListsWithOneElementEach() [][]int {
	return [][]int{
		{
			1,
		},
		{
			2,
		},
	}
}

func twoListsWithOneElementEachCartesianProduct() [][]int {
	return [][]int{
		{
			1, 2,
		},
	}
}

func oneElement() []int {
	return []int{1}
}

func threeElements() []int {
	return []int{1, 2, 3}
}
