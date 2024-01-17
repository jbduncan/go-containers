package graph

import (
	"github.com/jbduncan/go-containers/set"
)

var _ set.Set[EndpointPair[int]] = (*incidentEdgeSet[int])(nil)

type incidentEdgeSet[N comparable] struct {
	node          N
	adjacencyList map[N]set.MutableSet[N]
}

func (i incidentEdgeSet[N]) Contains(elem EndpointPair[N]) bool {
	adjacentNodes, ok := i.adjacencyList[i.node]
	if !ok {
		return false
	}

	return i.node == elem.Source() && adjacentNodes.Contains(elem.Target()) ||
		i.node == elem.Target() && adjacentNodes.Contains(elem.Source())
}

func (i incidentEdgeSet[N]) Len() int {
	adjacentNodes, ok := i.adjacencyList[i.node]
	if !ok {
		return 0
	}

	return adjacentNodes.Len()
}

func (i incidentEdgeSet[N]) ForEach(fn func(elem EndpointPair[N])) {
	adjacentNodes, ok := i.adjacencyList[i.node]
	if !ok {
		return
	}

	adjacentNodes.ForEach(
		func(adjNode N) {
			fn(EndpointPairOf(i.node, adjNode))
		})
}

func (i incidentEdgeSet[N]) String() string {
	return set.StringImpl[EndpointPair[N]](i)
}
