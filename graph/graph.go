package graph

import (
	"fmt"

	"github.com/jbduncan/go-containers/set"
)

// TODO: Docs
// TODO: Document with an example use
type Graph[N comparable] interface {
	Nodes() set.Set[N]
	Edges() set.Set[EndpointPair[N]]
	IsDirected() bool
	AllowsSelfLoops() bool
	// NodeOrder() ElementOrder
	// IncidentEdgeOrder() ElementOrder
	// TODO: Document that passing in an absent node panics, and to use Nodes().Contains() to tell when a
	//       given node is in the graph or not.
	AdjacentNodes(node N) set.Set[N]
	// TODO: Document that passing in an absent node panics, and to use Nodes().Contains() to tell when a
	//       given node is in the graph or not.
	Predecessors(node N) set.Set[N]
	// TODO: Document that passing in an absent node panics, and to use Nodes().Contains() to tell when a
	//       given node is in the graph or not.
	Successors(node N) set.Set[N]
	// TODO: Document that passing in an absent node panics, and to use Nodes().Contains() to tell when a
	//       given node is in the graph or not.
	IncidentEdges(node N) set.Set[EndpointPair[N]]
	// TODO: Document that passing in an absent node panics, and to use Nodes().Contains() to tell when a
	//       given node is in the graph or not.
	Degree(node N) int
	// TODO: Document that passing in an absent node panics, and to use Nodes().Contains() to tell when a
	//       given node is in the graph or not.
	InDegree(node N) int
	// TODO: Document that passing in an absent node panics, and to use Nodes().Contains() to tell when a
	//       given node is in the graph or not.
	OutDegree(node N) int
	HasEdgeConnecting(nodeU N, nodeV N) bool
	HasEdgeConnectingEndpoints(endpointPair EndpointPair[N]) bool
	// TODO: Implement Graph.String()
	// String() string
	// TODO: Is an Equal method needed to meet Guava's Graph::equals rules?
	//  If so, make the Equal method, discourage == from being used (documenting that its use is undefined), and
	//  optionally, if we decide to remove this interface, make the graph implementations have an incomparable field
	//  to force == to be unusable at compile time (see https://github.com/tailscale/tailscale/blob/main/types/structs/structs.go).
	// Equal(other Graph[N]) bool
}

// TODO: Docs
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
	return &graph[N]{
		adjacencyList:   map[N]set.MutableSet[N]{},
		allowsSelfLoops: b.allowsSelfLoops,
		numEdges:        0,
	}
}

// TODO: If the Graph and MutableGraph interfaces are ever eliminated, move them and these
//       compile-time type assertions to a test package.

var (
	_ Graph[int]        = (*graph[int])(nil)
	_ MutableGraph[int] = (*graph[int])(nil)
)

type graph[N comparable] struct {
	adjacencyList   map[N]set.MutableSet[N]
	allowsSelfLoops bool
	numEdges        int
}

func (m *graph[N]) IsDirected() bool {
	return false
}

func (m *graph[N]) AllowsSelfLoops() bool {
	return m.allowsSelfLoops
}

func (m *graph[N]) Nodes() set.Set[N] {
	return keySet[N]{
		delegate: m.adjacencyList,
	}
}

func (m *graph[N]) Edges() set.Set[EndpointPair[N]] {
	return edgeSet[N]{
		delegate: m,
	}
}

func (m *graph[N]) AdjacentNodes(node N) set.Set[N] {
	adjacentNodes, ok := m.adjacencyList[node]
	if !ok {
		// TODO: Go back to panicking, as this set is not a view and
		//       there is no sane way of testing that it's a view.
		//       Furthermore, panicking will allow programmers to
		//       flush out bugs faster.
		return set.Unmodifiable[N](set.New[N]())
	}

	return set.Unmodifiable(adjacentNodes)
}

func (m *graph[N]) Predecessors(node N) set.Set[N] {
	return m.AdjacentNodes(node)
}

func (m *graph[N]) Successors(node N) set.Set[N] {
	return m.AdjacentNodes(node)
}

func (m *graph[N]) IncidentEdges(node N) set.Set[EndpointPair[N]] {
	adjacentNodes, ok := m.adjacencyList[node]
	if !ok {
		// TODO: Go back to panicking, as this set is not a view and
		//       there is no sane way of testing that it's a view.
		//       Furthermore, panicking will allow programmers to
		//       flush out bugs faster.
		return set.Unmodifiable[EndpointPair[N]](set.New[EndpointPair[N]]())
	}

	return incidentEdgeSet[N]{
		node,
		adjacentNodes,
	}
}

func (m *graph[N]) Degree(node N) int {
	adjacentNodes, ok := m.adjacencyList[node]
	if !ok {
		return 0
	}

	return adjacentNodes.Len()
}

func (m *graph[N]) InDegree(node N) int {
	return m.Degree(node)
}

func (m *graph[N]) OutDegree(node N) int {
	return m.Degree(node)
}

func (m *graph[N]) HasEdgeConnecting(nodeU N, nodeV N) bool {
	adjacentNodes, ok := m.adjacencyList[nodeU]
	return ok && adjacentNodes.Contains(nodeV)
}

func (m *graph[N]) HasEdgeConnectingEndpoints(endpointPair EndpointPair[N]) bool {
	if endpointPair.IsOrdered() {
		return false
	}

	return m.HasEdgeConnecting(endpointPair.NodeU(), endpointPair.NodeV())
}

func (m *graph[N]) AddNode(node N) bool {
	if _, ok := m.adjacencyList[node]; ok {
		return false
	}

	m.adjacencyList[node] = set.New[N]()
	return true
}

func (m *graph[N]) PutEdge(nodeU N, nodeV N) bool {
	if !m.AllowsSelfLoops() && nodeU == nodeV {
		panic("self-loops are disallowed")
	}

	addedUToV := m.putEdge(nodeU, nodeV)
	m.putEdge(nodeV, nodeU)

	if addedUToV {
		m.numEdges++
	}

	// TODO: return booleans at all the right times
	return false
}

func (m *graph[N]) putEdge(nodeU N, nodeV N) bool {
	added := false
	adjacentNodes, ok := m.adjacencyList[nodeU]
	if !ok {
		adjacentNodes = set.New[N]()
		m.adjacencyList[nodeU] = adjacentNodes
		added = true
	}
	if adjacentNodes.Add(nodeV) {
		added = true
	}
	return added
}

func (m *graph[N]) RemoveNode(node N) bool {
	adjacentNodes, ok := m.adjacencyList[node]
	if !ok {
		return false
	}

	delete(m.adjacencyList, node)

	for _, adjacentNodes := range m.adjacencyList {
		if !adjacentNodes.Remove(node) {
			panic(
				fmt.Sprintf(
					"Unexpected: adjacent node %v was not removed",
					node))
		}
	}

	m.numEdges -= adjacentNodes.Len()

	return true
}

func (m *graph[N]) RemoveEdge(nodeU N, nodeV N) bool {
	removedUToV := m.removeEdge(nodeU, nodeV)
	removedVToU := m.removeEdge(nodeV, nodeU)

	if removedUToV != removedVToU {
		panic(
			fmt.Sprintf(
				"Unexpected: removedUToV (%t) != removedVToU (%t)",
				removedUToV, removedVToU))
	}

	m.numEdges--

	return removedUToV
}

func (m *graph[N]) removeEdge(from N, to N) bool {
	adjacentNodes, ok := m.adjacencyList[from]
	if !ok {
		return false
	}

	return adjacentNodes.Remove(to)
}
