package graph

import (
	"iter"

	"github.com/jbduncan/go-containers/set"
)

type incidentEdgeSet[N comparable] struct {
	node     N
	delegate *Mutable[N]
}

func (i incidentEdgeSet[N]) Contains(element EndpointPair[N]) bool {
	return i.containsForwards(element) || i.containsBackwards(element)
}

func (i incidentEdgeSet[N]) containsForwards(element EndpointPair[N]) bool {
	return i.node == element.Source() &&
		i.delegate.Successors(i.node).Contains(element.Target())
}

func (i incidentEdgeSet[N]) containsBackwards(element EndpointPair[N]) bool {
	return i.node == element.Target() &&
		i.delegate.Predecessors(i.node).Contains(element.Source())
}

func (i incidentEdgeSet[N]) Len() int {
	return i.delegate.AdjacentNodes(i.node).Len()
}

func (i incidentEdgeSet[N]) All() iter.Seq[EndpointPair[N]] {
	if !i.delegate.IsDirected() {
		return i.allDirected()
	}

	return i.allUndirected()
}

func (i incidentEdgeSet[N]) allDirected() iter.Seq[EndpointPair[N]] {
	return func(yield func(EndpointPair[N]) bool) {
		for adjNode := range i.delegate.AdjacentNodes(i.node).All() {
			if !yield(EndpointPairOf(i.node, adjNode)) {
				return
			}
		}
	}
}

func (i incidentEdgeSet[N]) allUndirected() iter.Seq[EndpointPair[N]] {
	return func(yield func(EndpointPair[N]) bool) {
		for predecessor := range i.delegate.Predecessors(i.node).All() {
			if !yield(EndpointPairOf(predecessor, i.node)) {
				return
			}
		}
		for successor := range i.delegate.Successors(i.node).All() {
			if i.delegate.Predecessors(i.node).Contains(successor) {
				continue
			}
			if !yield(EndpointPairOf(i.node, successor)) {
				return
			}
		}
	}
}

func (i incidentEdgeSet[N]) String() string {
	return set.StringImpl[EndpointPair[N]](i)
}
