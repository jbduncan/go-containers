package graph

import (
	"iter"

	"github.com/jbduncan/go-containers/set"
)

type edgeSet[N comparable] struct {
	delegate *Mutable[N]
	len      func() int
}

func (e edgeSet[N]) Contains(element EndpointPair[N]) bool {
	return e.delegate.Nodes().Contains(element.Source()) &&
		e.delegate.Successors(element.Source()).Contains(element.Target())
}

func (e edgeSet[N]) Len() int {
	return e.len()
}

func (e edgeSet[N]) All() iter.Seq[EndpointPair[N]] {
	return func(yield func(EndpointPair[N]) bool) {
		seen := set.Of[EndpointPair[N]]()

		for source := range e.delegate.Nodes().All() {
			for target := range e.delegate.Successors(source).All() {
				edge := EndpointPairOf(source, target)
				if e.edgeSeen(edge, seen) {
					continue
				}

				seen.Add(edge)
				if !yield(edge) {
					return
				}
			}
		}
	}
}

func (e edgeSet[N]) edgeSeen(
	edge EndpointPair[N],
	seen set.Set[EndpointPair[N]],
) bool {
	if seen.Contains(edge) {
		return true
	}
	return !e.delegate.IsDirected() && seen.Contains(reverseOf(edge))
}

func reverseOf[N comparable](edge EndpointPair[N]) EndpointPair[N] {
	return EndpointPairOf(edge.Target(), edge.Source())
}

func (e edgeSet[N]) String() string {
	return set.StringImpl[EndpointPair[N]](e)
}
