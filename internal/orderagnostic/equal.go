package orderagnostic

import "maps"

func SlicesEqual[T comparable](a []T, b []T) bool {
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
