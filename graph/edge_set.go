package graph

import (
	"iter"

	"github.com/jbduncan/go-containers/set"
)

var _ set.Set[EndpointPair[int]] = (*edgeSet[int])(nil)

type edgeSet[N comparable] struct {
	delegate Graph[N]
	len      func() int
}

func (e edgeSet[N]) Contains(elem EndpointPair[N]) bool {
	return e.delegate.Nodes().Contains(elem.Source()) &&
		e.delegate.Successors(elem.Source()).Contains(elem.Target())
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
	if !e.delegate.IsDirected() {
		reverse := EndpointPairOf(edge.Target(), edge.Source())
		if seen.Contains(reverse) {
			return true
		}
	}
	return false
}

func (e edgeSet[N]) String() string {
	return set.StringImpl[EndpointPair[N]](e)
}
