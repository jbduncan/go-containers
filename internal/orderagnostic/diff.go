package orderagnostic

import "github.com/google/go-cmp/cmp"

func Diff[T comparable](got []T, want []T) string {
	return cmp.Diff(want, got, cmp.Comparer(SlicesEqual[T]))
}
