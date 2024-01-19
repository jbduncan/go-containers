package graph

import (
	"github.com/jbduncan/go-containers/set"
)

var _ set.Set[EndpointPair[int]] = (*edgeSet[int])(nil)

type edgeSet[N comparable] struct {
	delegate *undirectedGraph[N]
}

func (e edgeSet[N]) Contains(elem EndpointPair[N]) bool {
	return e.delegate.Nodes().Contains(elem.Source()) &&
		e.delegate.AdjacentNodes(elem.Source()).Contains(elem.Target())
}

func (e edgeSet[N]) Len() int {
	return e.delegate.numEdges
}

func (e edgeSet[N]) ForEach(fn func(elem EndpointPair[N])) {
	result := set.NewMutable[EndpointPair[N]]()
	e.delegate.Nodes().ForEach(func(s N) {
		e.delegate.AdjacentNodes(s).ForEach(func(t N) {
			st := EndpointPairOf(s, t)
			tu := EndpointPairOf(t, s)
			if !result.Contains(st) && !result.Contains(tu) {
				result.Add(st)
				fn(st)
			}
		})
	})
}

func (e edgeSet[N]) String() string {
	return set.StringImpl[EndpointPair[N]](e)
}
