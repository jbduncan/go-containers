package graph

import (
	"iter"
	"maps"

	"github.com/jbduncan/go-containers/set"
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

func (k keySet[T]) All() iter.Seq[T] {
	return maps.Keys(k.delegate)
}

func (k keySet[N]) String() string {
	return set.StringImpl[N](k)
}
