package orderagnostic

import (
	"maps"

	"github.com/google/go-cmp/cmp"
	"github.com/jbduncan/go-containers/internal/slicesx"
)

func Diff[T comparable](
	got []T,
	want []T,
	extraOptions ...cmp.Option,
) string {
	return cmp.Diff(
		want,
		got,
		slicesx.AllOf(
			cmp.Comparer(slicesEqual[T]),
			extraOptions,
		)...,
	)
}

func slicesEqual[T comparable](a []T, b []T) bool {
	x := make(map[T]int)
	for _, value := range a {
		x[value]++
	}
	y := make(map[T]int)
	for _, value := range b {
		y[value]++
	}
	return maps.Equal(x, y)
}
