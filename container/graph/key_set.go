package graph

import (
	"fmt"
	"strings"

	"go-containers/container/set"
)

var _ set.Set[int] = (*keySet[int])(nil)

type keySet[N comparable] struct {
	delegate map[N]set.MutableSet[N]
}

func (k keySet[N]) Contains(elem N) bool {
	_, ok := k.delegate[elem]
	return ok
}

func (k keySet[N]) Len() int {
	return len(k.delegate)
}

func (k keySet[N]) ForEach(fn func(elem N)) {
	for key := range k.delegate {
		fn(key)
	}
}

func (k keySet[N]) ToSlice() []N {
	result := make([]N, 0, k.Len())

	for elem := range k.delegate {
		result = append(result, elem)
	}

	return result
}

func (k keySet[N]) String() string {
	var builder strings.Builder

	builder.WriteRune('[')
	index := 0
	for elem := range k.delegate {
		if index > 0 {
			builder.WriteString(", ")
		}

		builder.WriteString(fmt.Sprintf("%v", elem))
		index++
	}

	builder.WriteRune(']')
	return builder.String()
}
