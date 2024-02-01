package graph

import (
	"github.com/jbduncan/go-containers/set"
)

var _ set.Set[EndpointPair[int]] = (*incidentEdgeSet[int])(nil)

type incidentEdgeSet[N comparable] struct {
	node     N
	delegate Graph[N]
}

func (i incidentEdgeSet[N]) Contains(elem EndpointPair[N]) bool {
	return i.node == elem.Source() && i.delegate.AdjacentNodes(i.node).Contains(elem.Target()) ||
		i.node == elem.Target() && i.delegate.AdjacentNodes(i.node).Contains(elem.Source())
}

func (i incidentEdgeSet[N]) Len() int {
	return i.delegate.AdjacentNodes(i.node).Len()
}

func (i incidentEdgeSet[N]) ForEach(fn func(elem EndpointPair[N])) {
	i.delegate.AdjacentNodes(i.node).ForEach(
		func(adjNode N) {
			fn(EndpointPairOf(i.node, adjNode))
		})
}

func (i incidentEdgeSet[N]) String() string {
	return set.StringImpl[EndpointPair[N]](i)
}
