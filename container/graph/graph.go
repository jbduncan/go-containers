package graph

import (
	"errors"
	"fmt"

	"go-containers/container/set"
)

var errNodeNotElementOfGraph = errors.New("node not an element of this graph")

type Graph[N comparable] interface {
	Nodes() set.Set[N]
	Edges() set.Set[EndpointPair[N]]
	// IsDirected() bool
	// AllowsSelfLoops() bool
	// NodeOrder() ElementOrder
	// IncidentEdgeOrder() ElementOrder
	AdjacentNodes(node N) (set.Set[N], error)
	// MustAdjacentNodes(node N) set.Set[N]
	Predecessors(node N) (set.Set[N], error)
	// MustPredecessors(node N) set.Set[N]
	Successors(node N) (set.Set[N], error)
	// MustSuccessors(node N) set.Set[N]
	IncidentEdges(node N) (set.Set[EndpointPair[N]], error)
	// MustIncidentEdges(node N) set.Set[EndpointPair[N]]
	Degree(node N) (int, error)
	InDegree(node N) (int, error)
	OutDegree(node N) (int, error)
	// HasEdgeConnecting(nodeU N, nodeV N) bool
	// HasEdgeConnectingEndpoints(endpointPair EndpointPair[N]) bool
	// String() string
	// TODO: Is an Equals function needed to meet Guava's Graph::equals rules?
	// Equal(other Graph[N]) bool
}

type MutableGraph[N comparable] interface {
	Graph[N]

	AddNode(node N) bool
	PutEdge(nodeU N, nodeV N) bool
	// PutEdgeWithEndpoints(e EndpointPair[N]) bool
	RemoveNode(node N) bool
	// RemoveEdge(nodeU N, nodeV N) bool
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
	return keySet[N]{
		delegate: m.adjacencyList,
	}
}

func (m *mutableGraph[N]) Edges() set.Set[EndpointPair[N]] {
	return set.Unmodifiable(set.New[EndpointPair[N]]())
}

type keySet[N comparable] struct {
	delegate map[N]set.MutableSet[N]
}

func (k keySet[N]) Contains(elem N) bool {
	// TODO implement me
	panic("implement me")
}

func (k keySet[N]) Len() int {
	// TODO implement me
	panic("implement me")
}

func (k keySet[N]) ForEach(fn func(elem N)) {
	for key := range k.delegate {
		fn(key)
	}
}

func (k keySet[N]) String() string {
	// TODO implement me
	panic("implement me")
}

func (m *mutableGraph[N]) AdjacentNodes(node N) (set.Set[N], error) {
	adjacentNodes, ok := m.adjacencyList[node]
	if !ok {
		return nil, fmt.Errorf("%v: %w", node, errNodeNotElementOfGraph)
	}

	return set.Unmodifiable(adjacentNodes), nil
}

func (m *mutableGraph[N]) Predecessors(node N) (set.Set[N], error) {
	if _, ok := m.adjacencyList[node]; !ok {
		return nil, fmt.Errorf("%v: %w", node, errNodeNotElementOfGraph)
	}

	return set.Unmodifiable(set.New[N]()), nil
}

func (m *mutableGraph[N]) Successors(node N) (set.Set[N], error) {
	if _, ok := m.adjacencyList[node]; !ok {
		return nil, fmt.Errorf("%v: %w", node, errNodeNotElementOfGraph)
	}

	return set.Unmodifiable(set.New[N]()), nil
}

func (m *mutableGraph[N]) IncidentEdges(node N) (set.Set[EndpointPair[N]], error) {
	if _, ok := m.adjacencyList[node]; !ok {
		return nil, fmt.Errorf("%v: %w", node, errNodeNotElementOfGraph)
	}

	return set.Unmodifiable(set.New[EndpointPair[N]]()), nil
}

func (m *mutableGraph[N]) Degree(node N) (int, error) {
	if _, ok := m.adjacencyList[node]; !ok {
		return 0, fmt.Errorf("%v: %w", node, errNodeNotElementOfGraph)
	}

	return m.adjacencyList[node].Len(), nil
}

func (m *mutableGraph[N]) InDegree(node N) (int, error) {
	if _, ok := m.adjacencyList[node]; !ok {
		return 0, fmt.Errorf("%v: %w", node, errNodeNotElementOfGraph)
	}

	return 0, nil
}

func (m *mutableGraph[N]) OutDegree(node N) (int, error) {
	if _, ok := m.adjacencyList[node]; !ok {
		return 0, fmt.Errorf("%v: %w", node, errNodeNotElementOfGraph)
	}

	return 0, nil
}

func (m *mutableGraph[N]) AddNode(node N) bool {
	if _, ok := m.adjacencyList[node]; ok {
		return false
	}

	m.adjacencyList[node] = set.New[N]()
	return true
}

func (m *mutableGraph[N]) PutEdge(nodeU N, nodeV N) bool {
	adjacentNodes, ok := m.adjacencyList[nodeU]
	if !ok {
		adjacentNodes = set.New[N]()
		m.adjacencyList[nodeU] = adjacentNodes
	}
	adjacentNodes.Add(nodeV)

	adjacentNodes, ok = m.adjacencyList[nodeV]
	if !ok {
		adjacentNodes = set.New[N]()
		m.adjacencyList[nodeV] = adjacentNodes
	}
	adjacentNodes.Add(nodeU)

	//TODO
	return false
}

func (m *mutableGraph[N]) RemoveNode(node N) bool {
	_, ok := m.adjacencyList[node]
	if !ok {
		return false
	}

	delete(m.adjacencyList, node)

	for _, adjacentNodes := range m.adjacencyList {
		adjacentNodes.Remove(node)
	}

	return true
}

// type ElementOrder struct {}
