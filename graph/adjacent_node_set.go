package graph

import "github.com/jbduncan/go-containers/set"

var _ set.Set[int] = (*adjacentNodeSet[int])(nil)

type adjacentNodeSet[N comparable] struct {
	node          N
	adjacencyList map[N]set.MutableSet[N]
}

func (a adjacentNodeSet[N]) Contains(elem N) bool {
	adjacentNodes, ok := a.adjacencyList[a.node]
	if !ok {
		return false
	}

	return adjacentNodes.Contains(elem)
}

func (a adjacentNodeSet[N]) Len() int {
	adjacentNodes, ok := a.adjacencyList[a.node]
	if !ok {
		return 0
	}

	return adjacentNodes.Len()
}

func (a adjacentNodeSet[N]) ForEach(fn func(elem N)) {
	adjacentNodes, ok := a.adjacencyList[a.node]
	if !ok {
		return
	}

	adjacentNodes.ForEach(fn)
}

func (a adjacentNodeSet[N]) String() string {
	return set.StringImpl[N](a)
}
