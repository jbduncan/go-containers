package graph

import "github.com/jbduncan/go-containers/set"

var _ set.Set[EndpointPair[int]] = (*incidentEdgeSet[int])(nil)

type incidentEdgeSet[N comparable] struct {
	node     N
	delegate Graph[N]
}

func (i incidentEdgeSet[N]) Contains(elem EndpointPair[N]) bool {
	return i.containsForwards(elem) || i.containsBackwards(elem)
}

func (i incidentEdgeSet[N]) containsForwards(elem EndpointPair[N]) bool {
	return i.node == elem.Source() &&
		i.delegate.Successors(i.node).Contains(elem.Target())
}

func (i incidentEdgeSet[N]) containsBackwards(elem EndpointPair[N]) bool {
	return i.node == elem.Target() &&
		i.delegate.Predecessors(i.node).Contains(elem.Source())
}

func (i incidentEdgeSet[N]) Len() int {
	return i.delegate.AdjacentNodes(i.node).Len()
}

func (i incidentEdgeSet[N]) ForEach(fn func(elem EndpointPair[N])) {
	if !i.delegate.IsDirected() {
		i.delegate.AdjacentNodes(i.node).ForEach(
			func(adjNode N) {
				fn(EndpointPairOf(i.node, adjNode))
			})
		return
	}

	i.delegate.Predecessors(i.node).ForEach(
		func(predecessor N) {
			fn(EndpointPairOf(predecessor, i.node))
		})
	i.delegate.Successors(i.node).ForEach(
		func(successor N) {
			if !i.delegate.Predecessors(i.node).Contains(successor) {
				fn(EndpointPairOf(i.node, successor))
			}
		})
}

func (i incidentEdgeSet[N]) String() string {
	return set.StringImpl[EndpointPair[N]](i)
}
