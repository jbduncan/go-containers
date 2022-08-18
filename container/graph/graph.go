package graph

import (
	"go-containers/container/set"
)

type Graph[N comparable] interface {
	Nodes() set.Set[N]
	Edges() set.Set[EndpointPair[N]]
	IsDirected() bool
	// AllowsSelfLoops() bool
	// NodeOrder() ElementOrder
	// IncidentEdgeOrder() ElementOrder
	// TODO: Document that passing in an absent node returns an empty set, and to use Nodes().Contains() to tell when a
	//       given node is in the graph or not.
	AdjacentNodes(node N) set.Set[N]
	// TODO: Document that passing in an absent node returns an empty set, and to use Nodes().Contains() to tell when a
	//       given node is in the graph or not.
	Predecessors(node N) set.Set[N]
	// TODO: Document that passing in an absent node returns an empty set, and to use Nodes().Contains() to tell when a
	//       given node is in the graph or not.
	Successors(node N) set.Set[N]
	// TODO: Document that passing in an absent node returns an empty set, and to use Nodes().Contains() to tell when a
	//       given node is in the graph or not.
	IncidentEdges(node N) set.Set[EndpointPair[N]]
	// TODO: Document that passing in an absent node returns zero, and to use Nodes().Contains() to tell when a
	//       given node is in the graph or not.
	Degree(node N) int
	// TODO: Document that passing in an absent node returns zero, and to use Nodes().Contains() to tell when a
	//       given node is in the graph or not.
	InDegree(node N) int
	// TODO: Document that passing in an absent node returns zero, and to use Nodes().Contains() to tell when a
	//       given node is in the graph or not.
	OutDegree(node N) int
	// HasEdgeConnecting(nodeU N, nodeV N) bool
	// HasEdgeConnectingEndpoints(endpointPair EndpointPair[N]) bool
	// String() string
	// TODO: Is an Equal function needed to meet Guava's Graph::equals rules?
	// Equal(other Graph[N]) bool
}

type MutableGraph[N comparable] interface {
	Graph[N]

	AddNode(node N) bool
	PutEdge(nodeU N, nodeV N) bool
	// PutEdgeWithEndpoints(e EndpointPair[N]) bool
	RemoveNode(node N) bool
	RemoveEdge(nodeU N, nodeV N) bool
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

func (m *mutableGraph[N]) IsDirected() bool {
	return false
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

func (m *mutableGraph[N]) AdjacentNodes(node N) set.Set[N] {
	adjacentNodes, ok := m.adjacencyList[node]
	if !ok {
		return set.Unmodifiable(set.New[N]())
	}

	return set.Unmodifiable(adjacentNodes)
}

func (m *mutableGraph[N]) Predecessors(node N) set.Set[N] {
	return m.AdjacentNodes(node)
}

func (m *mutableGraph[N]) Successors(node N) set.Set[N] {
	return m.AdjacentNodes(node)
}

func (m *mutableGraph[N]) IncidentEdges(node N) set.Set[EndpointPair[N]] {
	adjacentNodes, ok := m.adjacencyList[node]
	if !ok {
		return set.Unmodifiable(set.New[EndpointPair[N]]())
	}

	return incidentEdgeSet[N]{
		node,
		adjacentNodes,
	}
}

type incidentEdgeSet[N comparable] struct {
	node          N
	adjacentNodes set.MutableSet[N]
}

func (i incidentEdgeSet[N]) Contains(elem EndpointPair[N]) bool {
	//TODO implement me
	panic("implement me")
}

func (i incidentEdgeSet[N]) Len() int {
	//TODO implement me
	panic("implement me")
}

func (i incidentEdgeSet[N]) ForEach(fn func(elem EndpointPair[N])) {
	i.adjacentNodes.ForEach(func(adjNode N) {
		fn(NewUnorderedEndpointPair(i.node, adjNode))
	})
}

func (i incidentEdgeSet[N]) String() string {
	//TODO implement me
	panic("implement me")
}

func (m *mutableGraph[N]) Degree(node N) int {
	adjacentNodes, ok := m.adjacencyList[node]
	if !ok {
		return 0
	}

	return adjacentNodes.Len()
}

func (m *mutableGraph[N]) InDegree(node N) int {
	if _, ok := m.adjacencyList[node]; !ok {
		return 0
	}

	return 0
}

func (m *mutableGraph[N]) OutDegree(node N) int {
	if _, ok := m.adjacencyList[node]; !ok {
		return 0
	}

	return 0
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

func (m *mutableGraph[N]) RemoveEdge(nodeU N, nodeV N) bool {
	if _, ok := m.adjacencyList[nodeU]; !ok {
		return false
	}
	if _, ok := m.adjacencyList[nodeV]; !ok {
		return false
	}

	return true
}

// type ElementOrder struct {}
