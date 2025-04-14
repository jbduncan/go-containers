package graph

import (
	"iter"

	"github.com/jbduncan/go-containers/set"
)

type neighborSet[N comparable] struct {
	node            N
	nodeToNeighbors map[N]set.Set[N]
}

func (a neighborSet[N]) Contains(elem N) bool {
	if neighbors, ok := a.nodeToNeighbors[a.node]; ok {
		return neighbors.Contains(elem)
	}

	return false
}

func (a neighborSet[N]) Len() int {
	if neighbors, ok := a.nodeToNeighbors[a.node]; ok {
		return neighbors.Len()
	}

	return 0
}

func (a neighborSet[N]) All() iter.Seq[N] {
	return func(yield func(N) bool) {
		neighbors, ok := a.nodeToNeighbors[a.node]
		if !ok {
			return
		}

		for neighbor := range neighbors.All() {
			if !yield(neighbor) {
				return
			}
		}
	}
}

func (a neighborSet[N]) String() string {
	return set.StringImpl[N](a)
}
