package graph

import (
	"github.com/jbduncan/go-containers/set"
)

type graphWithEdgeCount[N comparable] interface {
	Graph[N]
	edgeCount() int
}

var _ set.Set[EndpointPair[int]] = (*edgeSet[int])(nil)

type edgeSet[N comparable] struct {
	delegate graphWithEdgeCount[N]
}

func (e edgeSet[N]) Contains(elem EndpointPair[N]) bool {
	return e.delegate.Nodes().Contains(elem.Source()) &&
		e.delegate.Successors(elem.Source()).Contains(elem.Target())
}

func (e edgeSet[N]) Len() int {
	return e.delegate.edgeCount()
}

func (e edgeSet[N]) ForEach(fn func(elem EndpointPair[N])) {
	seen := set.NewMutable[EndpointPair[N]]()

	e.delegate.Nodes().ForEach(func(source N) {
		e.delegate.Successors(source).ForEach(func(target N) {
			edge := EndpointPairOf(source, target)
			if e.delegate.IsDirected() {
				if seen.Contains(edge) {
					return
				}
			} else {
				reverse := EndpointPairOf(target, source)
				if seen.Contains(edge) || seen.Contains(reverse) {
					return
				}
			}

			seen.Add(edge)
			fn(edge)
		})
	})
}

func (e edgeSet[N]) String() string {
	return set.StringImpl[EndpointPair[N]](e)
}
