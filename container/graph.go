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
	IncidentEdges(n N) Set[EndpointPair[N]]
	Degree(n N) uint
	InDegree(n N) uint
	OutDegree(n N) uint
	HasEdgeConnecting(u N, v N) bool
	HasEdgeConnectingEndpoints(e EndpointPair[N]) bool
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

func NewUnorderedEndpointPair[N comparable](s N, t N) EndpointPair[N] {
	return EndpointPair[N]{
		u: s,
		v: t,
	}
}

type EndpointPair[N comparable] struct {
	u N
	v N
}

func (e EndpointPair[N]) NodeU() N {
	return e.u
}

func (e EndpointPair[N]) NodeV() N {
	return e.v
}

func (e EndpointPair[N]) AdjacentNode(n N) (N, error) {
	if n == e.u {
		return e.v, nil
	}
	if n == e.v {
		return e.u, nil
	}
	return *new(N), fmt.Errorf("EndpointPair %v does not contain node %#v", e.String(), n)
}

func (e EndpointPair[N]) IsOrdered() bool {
	return false
}

func (e EndpointPair[N]) String() string {
	return fmt.Sprintf("[%#v, %#v]", e.NodeU(), e.NodeV())
}

func NewOrderedEndpointPair[N comparable](s N, t N) OrderedEndpointPair[N] {
	return OrderedEndpointPair[N]{
		u: s,
		v: t,
	}
}

type OrderedEndpointPair[N comparable] struct {
	u N
	v N
}

func (o OrderedEndpointPair[N]) Source() N {
	return o.NodeU()
}

func (o OrderedEndpointPair[N]) Target() N {
	return o.NodeV()
}

func (o OrderedEndpointPair[N]) NodeU() N {
	return o.u
}

func (o OrderedEndpointPair[N]) NodeV() N {
	return o.v
}

func (o OrderedEndpointPair[N]) AdjacentNode(n N) (N, error) {
	if n == o.u {
		return o.v, nil
	}
	if n == o.v {
		return o.u, nil
	}
	return *new(N), fmt.Errorf("OrderedEndpointPair %v does not contain node %#v", o.String(), n)
}

func (o OrderedEndpointPair[N]) IsOrdered() bool {
	return true
}

func (o OrderedEndpointPair[N]) String() string {
	return fmt.Sprintf("<%#v -> %#v>", o.Source(), o.Target())
}
