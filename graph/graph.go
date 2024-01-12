package graph

import (
	"strconv"

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
	AdjacentNodes(node N) set.Set[N]
	Predecessors(node N) set.Set[N]
	Successors(node N) set.Set[N]
	IncidentEdges(node N) set.Set[EndpointPair[N]]
	Degree(node N) int
	InDegree(node N) int
	OutDegree(node N) int
	HasEdgeConnecting(nodeU N, nodeV N) bool
	HasEdgeConnectingEndpoints(endpointPair EndpointPair[N]) bool
	String() string
	// TODO: Make a graph.Equal method similar to set.Equal, discourage == from being used (documenting that its use is
	//  undefined), and optionally, if we decide to make the builder return a concrete type, make the graph
	//  implementations have an incomparable field to force == to be unusable at compile time (see
	//  https://github.com/tailscale/tailscale/blob/main/types/structs/structs.go).
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
	return adjacentNodeSet[N]{
		node:          node,
		adjacencyList: g.adjacencyList,
	}
}

func (g *graph[N]) Predecessors(node N) set.Set[N] {
	return g.AdjacentNodes(node)
}

func (g *graph[N]) Successors(node N) set.Set[N] {
	return g.AdjacentNodes(node)
}

func (g *graph[N]) IncidentEdges(node N) set.Set[EndpointPair[N]] {
	return incidentEdgeSet[N]{
		node:          node,
		adjacencyList: g.adjacencyList,
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
	if adjacentNodes, ok := g.adjacencyList[nodeU]; ok {
		return adjacentNodes.Contains(nodeV)
	}

	return false
}

func (g *graph[N]) HasEdgeConnectingEndpoints(endpointPair EndpointPair[N]) bool {
	if endpointPair.IsOrdered() {
		return false
	}

	return g.HasEdgeConnecting(endpointPair.NodeU(), endpointPair.NodeV())
}

func (g *graph[N]) String() string {
	return "isDirected: false, allowsSelfLoops: " +
		strconv.FormatBool(g.allowsSelfLoops) +
		", nodes: [], edges: []"
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
		return true
	}

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
		adjacentNodes.Remove(node)
	}

	g.numEdges -= adjacentNodes.Len()

	return true
}

func (g *graph[N]) RemoveEdge(nodeU, nodeV N) bool {
	removedUToV := g.removeEdge(nodeU, nodeV)
	g.removeEdge(nodeV, nodeU)

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
