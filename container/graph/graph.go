package graph

import (
	"fmt"
	"go-containers/container/set"
)

const fmtNodeNotElementOfGraph = "node %v not an element of this graph"

type Graph[N comparable] interface {
	Nodes() set.Set[N]
	// Edges() set.Set[EndpointPair[N]]
	// IsDirected() bool
	// AllowsSelfLoops() bool
	// NodeOrder() ElementOrder
	// IncidentEdgeOrder() ElementOrder
	AdjacentNodes(n N) (set.Set[N], error)
	Predecessors(n N) (set.Set[N], error)
	Successors(n N) (set.Set[N], error)
	IncidentEdges(n N) (set.Set[EndpointPair[N]], error)
	// TODO: Implement Degree next.
	// Degree(n N) (int, error)
	// InDegree(n N) (int, error)
	// OutDegree(n N) (int, error)
	// HasEdgeConnecting(u N, v N) bool
	// HasEdgeConnectingEndpoints(endpointPair EndpointPair[N]) bool
	// String() string
	// TODO: Is an Equals function needed to meet Guava's Graph::equals rules?
}

type MutableGraph[N comparable] interface {
	Graph[N]

	AddNode(n N) bool
	PutEdge(u N, v N) bool
	// PutEdgeWithEndpoints(e EndpointPair[N]) bool
	// RemoveNode(n N) bool
	// RemoveEdge(u N, v N) bool
	// RemoveEdgeWithEndpoints(e EndpointPair[N]) bool
}

func Undirected[N comparable](opts ...Option[N]) Builder[N] {
	return Builder[N]{}
}

type Option[N comparable] func(b Builder[N]) (Builder[N], error)

type Builder[N comparable] struct{}

func (b Builder[N]) Build() MutableGraph[N] {
	return &mutableGraph[N]{
		adjacencyList: map[N]set.MutableSet[N]{},
	}
}

type mutableGraph[N comparable] struct {
	adjacencyList map[N]set.MutableSet[N]
}

func (m *mutableGraph[N]) Nodes() set.Set[N] {
	return wrapKeys(m.adjacencyList)
}

func wrapKeys[N comparable](delegate map[N]set.MutableSet[N]) set.Set[N] {
	return keySet[N]{
		delegate: delegate,
	}
}

type keySet[N comparable] struct {
	delegate map[N]set.MutableSet[N]
}

func (k keySet[N]) Contains(elem N) bool {
	//TODO implement me
	panic("implement me")
}

func (k keySet[N]) Len() int {
	//TODO implement me
	panic("implement me")
}

func (k keySet[N]) ForEach(fn func(elem N)) {
	for key := range k.delegate {
		fn(key)
	}
}

func (k keySet[N]) String() string {
	//TODO implement me
	panic("implement me")
}

func (m *mutableGraph[N]) AdjacentNodes(n N) (set.Set[N], error) {
	adjacentNodes, ok := m.adjacencyList[n]
	if !ok {
		return nil, fmt.Errorf(fmtNodeNotElementOfGraph, n)
	}
	return set.Unmodifiable(adjacentNodes), nil
}

func (m *mutableGraph[N]) Predecessors(n N) (set.Set[N], error) {
	if _, ok := m.adjacencyList[n]; !ok {
		return nil, fmt.Errorf(fmtNodeNotElementOfGraph, n)
	}
	// TODO: Non-empty case(s)
	return set.Unmodifiable(set.New[N]()), nil
}

func (m *mutableGraph[N]) Successors(n N) (set.Set[N], error) {
	if _, ok := m.adjacencyList[n]; !ok {
		return nil, fmt.Errorf(fmtNodeNotElementOfGraph, n)
	}
	// TODO: Non-empty case(s)
	return set.Unmodifiable(set.New[N]()), nil
}

func (m *mutableGraph[N]) IncidentEdges(n N) (set.Set[EndpointPair[N]], error) {
	if _, ok := m.adjacencyList[n]; !ok {
		return nil, fmt.Errorf(fmtNodeNotElementOfGraph, n)
	}
	// TODO: Non-empty case(s)
	return set.Unmodifiable(set.New[EndpointPair[N]]()), nil
}

func (m *mutableGraph[N]) AddNode(n N) bool {
	if _, ok := m.adjacencyList[n]; ok {
		return false
	}
	m.adjacencyList[n] = set.New[N]()
	return true
}

func (m *mutableGraph[N]) PutEdge(u N, v N) bool {
	adjacentNodes, ok := m.adjacencyList[u]
	if !ok {
		adjacentNodes = set.New[N]()
		m.adjacencyList[u] = adjacentNodes
	}
	adjacentNodes.Add(v)

	adjacentNodes, ok = m.adjacencyList[v]
	if !ok {
		adjacentNodes = set.New[N]()
		m.adjacencyList[v] = adjacentNodes
	}
	adjacentNodes.Add(u)

	//TODO
	return false
}

// type ElementOrder struct {}
