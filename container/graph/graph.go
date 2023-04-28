package graph

import (
	"fmt"

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

// TODO: If the Graph and MutableGraph interfaces are ever eliminated, move them and these
//       compile-time type assertions to a test package.

var (
	_ Graph[int]        = (*mutableGraph[int])(nil)
	_ MutableGraph[int] = (*mutableGraph[int])(nil)
)

// TODO: Rename to `graph` for consistency with `set.set`.
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

// TODO: Add tests for all set.Set methods of mutableGraph.Nodes()

func (m *mutableGraph[N]) Nodes() set.Set[N] {
	return &keySet[N]{
		delegate: m.adjacencyList,
	}
}

// TODO: Add tests for all set.Set methods of mutableGraph.Edges()

func (m *mutableGraph[N]) Edges() set.Set[EndpointPair[N]] {
	// TODO: To get this code working, we used an O(E^2) algorithm. Let's
	//       ask on https://codereview.stackexchange.com/ to see if we can
	//       reduce the runtime to O(E). Worst case scenario, we can come up
	//       with a set that allows for custom equivalences, allowing the
	//       `contains` function to be replaced with an O(1) Set.Contains check
	//       that compares elements by EndpointPair.Equal. (Time to TDD a
	//       custom hash table into existence!)
	//       Note: If we go for a custom hash table, let's use FNV-1a (via
	//       fnv.New64a()), which is the hashing algorithm of choice in
	//       craftinginterpreters.com, Chapter 20. Note that the book chose it
	//       for brevity, not necessarily performance, but it seems a good
	//       start. If performance really matters to us in the future, we can
	//       read "The Art of Computer Programming, Volume 3", according to
	//       https://stackoverflow.com/a/34652.
	result := set.New[EndpointPair[N]]()
	m.Nodes().ForEach(func(u N) {
		// TODO: Replace .AdjacentNodes with .Successors when building
		//       a directed graph type.
		m.AdjacentNodes(u).ForEach(func(v N) {
			newEdge := NewUnorderedEndpointPair(u, v)
			if !contains(result, newEdge) {
				result.Add(newEdge)
			}
		})
	})
	return set.Unmodifiable(result)
}

func contains[T comparable](s set.MutableSet[EndpointPair[T]], value EndpointPair[T]) bool {
	result := false
	s.ForEach(func(existing EndpointPair[T]) {
		if value.Equal(existing) {
			result = true
		}
	})
	return result
}

// TODO: Add tests for all set.Set methods of mutableGraph.AdjacentNodes()
//       for both a present node and an absent node

func (m *mutableGraph[N]) AdjacentNodes(node N) set.Set[N] {
	adjacentNodes, ok := m.adjacencyList[node]
	if !ok {
		return set.Unmodifiable(set.New[N]())
	}

	return set.Unmodifiable(adjacentNodes)
}

// TODO: Add tests for all set.Set methods of mutableGraph.Predecessors()
//       for both a present node and an absent node

func (m *mutableGraph[N]) Predecessors(node N) set.Set[N] {
	return m.AdjacentNodes(node)
}

// TODO: Add tests for all set.Set methods of mutableGraph.Successors()
//       for both a present node and an absent node

func (m *mutableGraph[N]) Successors(node N) set.Set[N] {
	return m.AdjacentNodes(node)
}

// TODO: Add tests for all set.Set methods of mutableGraph.IncidentEdges()
//       for both a present node and an absent node

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
	return ok && adjacentNodes.Contains(nodeV)
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

	m.putEdge(nodeU, nodeV)
	m.putEdge(nodeV, nodeU)

	// TODO: return booleans at all the right times
	return false
}

func (m *mutableGraph[N]) putEdge(nodeU N, nodeV N) {
	adjacentNodes, ok := m.adjacencyList[nodeU]
	if !ok {
		adjacentNodes = set.New[N]()
		m.adjacencyList[nodeU] = adjacentNodes
	}
	adjacentNodes.Add(nodeV)
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
	removedUToV := m.removeEdge(nodeU, nodeV)
	removedVToU := m.removeEdge(nodeV, nodeU)

	if removedUToV != removedVToU {
		panic(
			fmt.Sprintf(
				"Unexpected: removedUToV (%t) != removedVToU (%t)",
				removedUToV, removedVToU))
	}

	return removedUToV
}

func (m *mutableGraph[N]) removeEdge(from N, to N) bool {
	adjacentNodes, ok := m.adjacencyList[from]
	if !ok {
		return false
	}

	// TODO: Simplify when MutableSet.Remove returns bools
	if adjacentNodes.Contains(to) {
		adjacentNodes.Remove(to)
		return true
	}
	return false
}
