package graph

import "github.com/jbduncan/go-containers/set"

var _ set.Set[int] = (*neighborSet[int])(nil)

type neighborSet[N comparable] struct {
	node            N
	nodeToNeighbors map[N]set.MutableSet[N]
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

func (a neighborSet[N]) ForEach(fn func(elem N)) {
	if neighbors, ok := a.nodeToNeighbors[a.node]; ok {
		neighbors.ForEach(fn)
	}
}

func (a neighborSet[N]) String() string {
	return set.StringImpl[N](a)
}
