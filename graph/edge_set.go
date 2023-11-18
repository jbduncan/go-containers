package graph

import (
	"fmt"
	"strings"

	"github.com/jbduncan/go-containers/set"
)

var _ set.Set[EndpointPair[int]] = (*edgeSet[int])(nil)

type edgeSet[N comparable] struct {
	delegate *graph[N]
}

func (e edgeSet[N]) Contains(elem EndpointPair[N]) bool {
	return e.delegate.IsDirected() == elem.IsOrdered() &&
		e.delegate.Nodes().Contains(elem.NodeU()) &&
		// TODO: Change to successors when accounting for directed graphs
		e.delegate.AdjacentNodes(elem.NodeU()).Contains(elem.NodeV())
}

func (e edgeSet[N]) Len() int {
	return e.delegate.numEdges
}

func (e edgeSet[N]) ForEach(fn func(elem EndpointPair[N])) {
	result := set.New[EndpointPair[N]]()
	e.delegate.Nodes().ForEach(func(u N) {
		// TODO: Replace .AdjacentNodes with .Successors when building
		//       a directed graph type.
		e.delegate.AdjacentNodes(u).ForEach(func(v N) {
			uv := NewUnorderedEndpointPair(u, v)
			vu := NewUnorderedEndpointPair(v, u)
			if !result.Contains(uv) && !result.Contains(vu) {
				result.Add(uv)
				fn(uv)
			}
		})
	})
}

func (e edgeSet[N]) String() string {
	var builder strings.Builder

	builder.WriteRune('[')
	index := 0
	e.ForEach(func(elem EndpointPair[N]) {
		if index > 0 {
			builder.WriteString(", ")
		}

		builder.WriteString(fmt.Sprintf("%v", elem))
		index++
	})

	builder.WriteRune(']')
	return builder.String()
}
