package graph

import "go-containers/container/set"

type Graph[N comparable] interface {
	Nodes() set.Set[N]
	// Edges() set.Set[EndpointPair[N]]
	// IsDirected() bool
	// AllowsSelfLoops() bool
	// NodeOrder() ElementOrder
	// IncidentEdgeOrder() ElementOrder
	// AdjacentNodes() set.Set[N]
	// Predecessors() set.Set[N]
	// Successors() set.Set[N]
	// IncidentEdges(node N) set.Set[EndpointPair[N]]
	// Degree(node N) int
	// InDegree(node N) int
	// OutDegree(node N) int
	// HasEdgeConnecting(nodeU N, nodeV N) bool
	// HasEdgeConnectingEndpoints(endpointPair EndpointPair[N]) bool
	// String() string
	// TODO: Is an Equals function needed to meet Guava's Graph::equals rules?
}

type MutableGraph[N comparable] interface {
	Graph[N]

	AddNode(n N) bool
	// PutEdge(u N, v N) bool
	// PutEdgeWithEndpoints(e EndpointPair[N]) bool
	// RemoveNode(n N) bool
	// RemoveEdge(u N, v N) bool
	// RemoveEdgeWithEndpoints(e EndpointPair[N]) bool
}

func Undirected[N comparable](opts ...Option[N]) Builder[N] {
	return Builder[N]{}
}

type Option[N comparable] func(graphBuilder Builder[N]) (Builder[N], error)

type Builder[N comparable] struct{}

func (b Builder[N]) Build() MutableGraph[N] {
	return &mutableGraph[N]{
		nodes: set.New[N](),
	}
}

type mutableGraph[N comparable] struct {
	nodes set.MutableSet[N]
}

func (m *mutableGraph[N]) Nodes() set.Set[N] {
	// TODO: Make sure is unmodifiable
	return m.nodes
}

func (m *mutableGraph[N]) AddNode(n N) bool {
	if m.nodes.Contains(n) {
		return false
	}
	m.nodes.Add(n)
	return true
}

// type ElementOrder struct {}
