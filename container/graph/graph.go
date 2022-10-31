package graph

import (
	"go-containers/container/set"
)

type Graph[N comparable] interface {
	Nodes() set.Set[N]
	Edges() set.Set[EndpointPair[N]]
	IsDirected() bool
	AllowsSelfLoops() bool
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
	HasEdgeConnecting(nodeU N, nodeV N) bool
	HasEdgeConnectingEndpoints(endpointPair EndpointPair[N]) bool
	// String() string
	// TODO: Is an Equal function needed to meet Guava's Graph::equals rules?
	// Equal(other Graph[N]) bool
}

type MutableGraph[N comparable] interface {
	Graph[N]

	AddNode(node N) bool
	// TODO: Document that PutEdge will panic if Graph.AllowsSelfLoops() is true and nodeU and nodeV are equal
	//       according to ==, and to check that nodeU != nodeV beforehand.
	PutEdge(nodeU N, nodeV N) bool
	// PutEdgeWithEndpoints(e EndpointPair[N]) bool
	RemoveNode(node N) bool
	RemoveEdge(nodeU N, nodeV N) bool
	// RemoveEdgeWithEndpoints(e EndpointPair[N]) bool
}

func Undirected[N comparable]() Builder[N] {
	return Builder[N]{}
}

type Builder[N comparable] struct {
	allowsSelfLoops bool
}

func (b Builder[N]) AllowsSelfLoops(allowsSelfLoops bool) Builder[N] {
	b.allowsSelfLoops = allowsSelfLoops
	return b
}

// TODO: Consider returning a public version of the concrete type, rather
//       than the MutableGraph interface, to allow new methods to be
//       added without breaking backwards compatibility:
//       - https://github.com/golang/go/wiki/CodeReviewComments#interfaces

func (b Builder[N]) Build() MutableGraph[N] {
	return &mutableGraph[N]{
		adjacencyList:   map[N]set.MutableSet[N]{},
		allowsSelfLoops: b.allowsSelfLoops,
	}
}

type mutableGraph[N comparable] struct {
	adjacencyList   map[N]set.MutableSet[N]
	allowsSelfLoops bool
}

func (m *mutableGraph[N]) IsDirected() bool {
	return false
}

func (m *mutableGraph[N]) AllowsSelfLoops() bool {
	return m.allowsSelfLoops
}

func (m *mutableGraph[N]) Nodes() set.Set[N] {
	return keySet[N]{
		delegate: m.adjacencyList,
	}
}

func (m *mutableGraph[N]) Edges() set.Set[EndpointPair[N]] {
	// TODO: flesh out
	return set.Unmodifiable(set.New[EndpointPair[N]]())
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

func (m *mutableGraph[N]) Degree(node N) int {
	adjacentNodes, ok := m.adjacencyList[node]
	if !ok {
		return 0
	}

	return adjacentNodes.Len()
}

func (m *mutableGraph[N]) InDegree(node N) int {
	return m.Degree(node)
}

func (m *mutableGraph[N]) OutDegree(node N) int {
	return m.Degree(node)
}

func (m *mutableGraph[N]) HasEdgeConnecting(nodeU N, nodeV N) bool {
	adjacentNodes, ok := m.adjacencyList[nodeU]
	if !ok {
		return false
	}

	if !adjacentNodes.Contains(nodeV) {
		return false
	}

	adjacentNodes = m.adjacencyList[nodeV]

	return adjacentNodes.Contains(nodeU)
}

func (m *mutableGraph[N]) HasEdgeConnectingEndpoints(endpointPair EndpointPair[N]) bool {
	if endpointPair.IsOrdered() {
		return false
	}

	return m.HasEdgeConnecting(endpointPair.NodeU(), endpointPair.NodeV())
}

func (m *mutableGraph[N]) AddNode(node N) bool {
	if _, ok := m.adjacencyList[node]; ok {
		return false
	}

	m.adjacencyList[node] = set.New[N]()
	return true
}

func (m *mutableGraph[N]) PutEdge(nodeU N, nodeV N) bool {
	if !m.AllowsSelfLoops() && nodeU == nodeV {
		panic("self-loops are disallowed")
	}

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

	//TODO: return booleans at all the right times
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
	adjacentNodes, ok := m.adjacencyList[nodeU]
	if !ok {
		return false
	}
	adjacentNodes.Remove(nodeV)

	adjacentNodes, ok = m.adjacencyList[nodeV]
	if !ok {
		return false
	}
	adjacentNodes.Remove(nodeU)

	return true
}
