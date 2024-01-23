package graph

import "github.com/jbduncan/go-containers/set"

var _ set.Set[int] = (*adjacentNodeSet[int])(nil)

type adjacentNodeSet[N comparable] struct {
	node                N
	nodeToAdjacentNodes map[N]set.MutableSet[N]
}

func (a adjacentNodeSet[N]) Contains(elem N) bool {
	if adjacentNodes, ok := a.nodeToAdjacentNodes[a.node]; ok {
		return adjacentNodes.Contains(elem)
	}

	return false
}

func (a adjacentNodeSet[N]) Len() int {
	if adjacentNodes, ok := a.nodeToAdjacentNodes[a.node]; ok {
		return adjacentNodes.Len()
	}

	return 0
}

func (a adjacentNodeSet[N]) ForEach(fn func(elem N)) {
	if adjacentNodes, ok := a.nodeToAdjacentNodes[a.node]; ok {
		adjacentNodes.ForEach(fn)
	}
}

func (a adjacentNodeSet[N]) String() string {
	return set.StringImpl[N](a)
}
