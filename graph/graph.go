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

func (g *graph[N]) IsDirected() bool {
	return false
}

func (g *graph[N]) AllowsSelfLoops() bool {
	return g.allowsSelfLoops
}

func (g *graph[N]) Nodes() set.Set[N] {
	return keySet[N]{
		delegate: g.adjacencyList,
	}
}

func (g *graph[N]) Edges() set.Set[EndpointPair[N]] {
	return edgeSet[N]{
		delegate: g,
	}
}

func (g *graph[N]) AdjacentNodes(node N) set.Set[N] {
	adjacentNodes, ok := g.adjacencyList[node]
	if !ok {
		// TODO: Go back to panicking, as this set is not a view and
		//       there is no sane way of testing that it's a view.
		//       Furthermore, panicking will allow programmers to
		//       flush out bugs faster.
		return set.Unmodifiable[N](set.NewMutable[N]())
	}

	return set.Unmodifiable(adjacentNodes)
}

func (g *graph[N]) Predecessors(node N) set.Set[N] {
	return g.AdjacentNodes(node)
}

func (g *graph[N]) Successors(node N) set.Set[N] {
	return g.AdjacentNodes(node)
}

func (g *graph[N]) IncidentEdges(node N) set.Set[EndpointPair[N]] {
	adjacentNodes, ok := g.adjacencyList[node]
	if !ok {
		// TODO: Go back to panicking, as this set is not a view and
		//       there is no sane way of testing that it's a view.
		//       Furthermore, panicking will allow programmers to
		//       flush out bugs faster.
		return set.Unmodifiable[EndpointPair[N]](set.NewMutable[EndpointPair[N]]())
	}

	return incidentEdgeSet[N]{
		node,
		adjacentNodes,
	}
}

func (g *graph[N]) Degree(node N) int {
	adjacentNodes, ok := g.adjacencyList[node]
	if !ok {
		return 0
	}

	return adjacentNodes.Len()
}

func (g *graph[N]) InDegree(node N) int {
	return g.Degree(node)
}

func (g *graph[N]) OutDegree(node N) int {
	return g.Degree(node)
}

func (g *graph[N]) HasEdgeConnecting(nodeU, nodeV N) bool {
	adjacentNodes, ok := g.adjacencyList[nodeU]
	return ok && adjacentNodes.Contains(nodeV)
}

func (g *graph[N]) HasEdgeConnectingEndpoints(endpointPair EndpointPair[N]) bool {
	if endpointPair.IsOrdered() {
		return false
	}

	return g.HasEdgeConnecting(endpointPair.NodeU(), endpointPair.NodeV())
}

func (g *graph[N]) AddNode(node N) bool {
	if _, ok := g.adjacencyList[node]; ok {
		return false
	}

	g.adjacencyList[node] = set.NewMutable[N]()
	return true
}

func (g *graph[N]) PutEdge(nodeU, nodeV N) bool {
	if !g.AllowsSelfLoops() && nodeU == nodeV {
		panic("self-loops are disallowed")
	}

	addedUToV := g.putEdge(nodeU, nodeV)
	g.putEdge(nodeV, nodeU)

	if addedUToV {
		g.numEdges++
	}

	// TODO: return booleans at all the right times
	return false
}

func (g *graph[N]) putEdge(nodeU, nodeV N) bool {
	added := false
	adjacentNodes, ok := g.adjacencyList[nodeU]
	if !ok {
		adjacentNodes = set.NewMutable[N]()
		g.adjacencyList[nodeU] = adjacentNodes
		added = true
	}
	if adjacentNodes.Add(nodeV) {
		added = true
	}
	return added
}

func (g *graph[N]) RemoveNode(node N) bool {
	adjacentNodes, ok := g.adjacencyList[node]
	if !ok {
		return false
	}

	delete(g.adjacencyList, node)

	for _, adjacentNodes := range g.adjacencyList {
		if !adjacentNodes.Remove(node) {
			panic(
				fmt.Sprintf(
					"Unexpected: adjacent node %v was not removed",
					node))
		}
	}

	g.numEdges -= adjacentNodes.Len()

	return true
}

func (g *graph[N]) RemoveEdge(nodeU, nodeV N) bool {
	removedUToV := g.removeEdge(nodeU, nodeV)
	removedVToU := g.removeEdge(nodeV, nodeU)

	if removedUToV != removedVToU {
		panic(
			fmt.Sprintf(
				"Unexpected: removedUToV (%t) != removedVToU (%t)",
				removedUToV, removedVToU))
	}

	g.numEdges--

	return removedUToV
}

func (g *graph[N]) removeEdge(from, to N) bool {
	adjacentNodes, ok := g.adjacencyList[from]
	if !ok {
		return false
	}

	return adjacentNodes.Remove(to)
}
