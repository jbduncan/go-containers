package graph

import "go-containers/container/set"

type incidentEdgeSet[N comparable] struct {
	node          N
	adjacentNodes set.MutableSet[N]
}

func (i incidentEdgeSet[N]) Contains(elem EndpointPair[N]) bool {
	if !elem.IsOrdered() {
		// TODO: Once EndpointPair.Equal() is introduced, replace this code with an
		//       EndpointPair.Equal() check
		return (i.node == elem.NodeU() && i.adjacentNodes.Contains(elem.NodeV())) ||
			(i.node == elem.NodeV() && i.adjacentNodes.Contains(elem.NodeU()))
	}
	return false
}

func (i incidentEdgeSet[N]) Len() int {
	return i.adjacentNodes.Len()
}

func (i incidentEdgeSet[N]) ForEach(fn func(elem EndpointPair[N])) {
	i.adjacentNodes.ForEach(
		func(adjNode N) {
			fn(NewUnorderedEndpointPair(i.node, adjNode))
		})
}

func (i incidentEdgeSet[N]) String() string {
	//TODO implement me
	panic("implement me")
}