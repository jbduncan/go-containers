package graph

import (
	"fmt"
	"strings"

	"github.com/jbduncan/go-containers/set"
)

var _ set.Set[EndpointPair[int]] = (*incidentEdgeSet[int])(nil)

type incidentEdgeSet[N comparable] struct {
	node          N
	adjacentNodes set.MutableSet[N]
}

func (i incidentEdgeSet[N]) Contains(elem EndpointPair[N]) bool {
	if !elem.IsOrdered() {
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
	var builder strings.Builder

	builder.WriteRune('[')
	index := 0
	i.ForEach(func(elem EndpointPair[N]) {
		if index > 0 {
			builder.WriteString(", ")
		}

		builder.WriteString(fmt.Sprintf("%v", elem))
		index++
	})

	builder.WriteRune(']')
	return builder.String()
}
