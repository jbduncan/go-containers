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
	if elem.IsOrdered() {
		return false
	}

	adjacentNodes, ok := i.adjacencyList[i.node]
	if !ok {
		return false
	}

	return i.node == elem.NodeU() && adjacentNodes.Contains(elem.NodeV()) ||
		i.node == elem.NodeV() && adjacentNodes.Contains(elem.NodeU())
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
			fn(UnorderedEndpointPair(i.node, adjNode))
		})
}

func (i incidentEdgeSet[N]) String() string {
	return set.StringImpl[EndpointPair[N]](i)
}
