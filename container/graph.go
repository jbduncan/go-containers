package container

import "fmt"

type Graph[N comparable] interface {
	Nodes() Set[N]
	Edges() Set[EndpointPair[N]]
	IsDirected() bool
	AllowsSelfLoops() bool
	// TODO: Implement
	// NodeOrder() ElementOrder
	// IncidentEdgeOrder() ElementOrder
	AdjacentNodes() Set[N]
	Predecessors() Set[N]
	Successors() Set[N]
	IncidentEdges(node N) Set[EndpointPair[N]]
	Degree(node N) int
	InDegree(node N) int
	OutDegree(node N) int
	HasEdgeConnecting(nodeU N, nodeV N) bool
	HasEdgeConnectingEndpoints(endpointPair EndpointPair[N]) bool
	String() string
	// TODO: Is an Equals function needed to meet Guava's Graph::equals rules?
}

type MutableGraph[N comparable] interface {
	Graph[N]

	AddNode(n N) bool
	PutEdge(u N, v N) bool
	PutEdgeWithEndpoints(e EndpointPair[N]) bool
	RemoveNode(n N) bool
	RemoveEdge(u N, v N) bool
	RemoveEdgeWithEndpoints(e EndpointPair[N]) bool
}

// TODO: Implement
// type ElementOrder struct {}

func NewUnorderedEndpointPair[N comparable](nodeU N, nodeV N) EndpointPair[N] {
	return EndpointPair[N]{
		nodeU:     nodeU,
		nodeV:     nodeV,
		isOrdered: false,
	}
}

func NewOrderedEndpointPair[N comparable](source N, target N) EndpointPair[N] {
	return EndpointPair[N]{
		nodeU:     source,
		nodeV:     target,
		isOrdered: true,
	}
}

type EndpointPair[N comparable] struct {
	nodeU     N
	nodeV     N
	isOrdered bool
}

func (e EndpointPair[N]) Source() N {
	if e.isOrdered {
		return e.nodeU
	}
	panic("cannot call Source()/Target() on an EndpointPair from an undirected graph; consider calling " +
		"AdjacentNode(node) if you already have a node, or NodeU()/NodeV() if you don't")
}

func (e EndpointPair[N]) Target() N {
	if e.isOrdered {
		return e.nodeV
	}
	panic("cannot call Source()/Target() on an EndpointPair from an undirected graph; consider calling " +
		"AdjacentNode(node) if you already have a node, or NodeU()/NodeV() if you don't")
}

func (e EndpointPair[N]) NodeU() N {
	return e.nodeU
}

func (e EndpointPair[N]) NodeV() N {
	return e.nodeV
}

func (e EndpointPair[N]) AdjacentNode(n N) N {
	if n == e.nodeU {
		return e.nodeV
	}
	if n == e.nodeV {
		return e.nodeU
	}
	panic(fmt.Sprintf("EndpointPair %v does not contain node %v", e, n))
}

func (e EndpointPair[N]) IsOrdered() bool {
	return e.isOrdered
}

// TODO: EndpointPair: make Equals method and discourage == from being used (documenting that its use is undefined).
//       See this link:
//       https://github.com/google/guava/blob/4d323b2b117a5906ab16074c8c88b4ff162b1b82/guava/src/com/google/common/graph/EndpointPair.java#L131-L145

func (e EndpointPair[N]) String() string {
	if e.isOrdered {
		return fmt.Sprintf("<%v -> %v>", e.Source(), e.Target())
	}
	return fmt.Sprintf("[%v, %v]", e.NodeU(), e.NodeV())
}
