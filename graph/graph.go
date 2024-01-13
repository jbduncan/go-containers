package graph

import (
	"strconv"

	"github.com/jbduncan/go-containers/set"
)

type Graph[N comparable] interface {
	Nodes() set.Set[N]
	Edges() set.Set[EndpointPair[N]]
	IsDirected() bool
	AllowsSelfLoops() bool
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
}

type MutableGraph[N comparable] interface {
	Graph[N]

	AddNode(node N) bool
	PutEdge(nodeU N, nodeV N) bool
	RemoveNode(node N) bool
	RemoveEdge(nodeU N, nodeV N) bool
}

func Undirected[N comparable]() Builder[N] {
	return Builder[N]{}
}

func Directed[N comparable]() Builder[N] {
	return Builder[N]{}
}

type Builder[N comparable] struct {
	allowsSelfLoops bool
}

func (b Builder[N]) AllowsSelfLoops(allowsSelfLoops bool) Builder[N] {
	b.allowsSelfLoops = allowsSelfLoops
	return b
}

func (b Builder[N]) Build() MutableGraph[N] {
	return &graph[N]{
		adjacencyList:   map[N]set.MutableSet[N]{},
		allowsSelfLoops: b.allowsSelfLoops,
		numEdges:        0,
	}
}

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
		", nodes: " +
		g.Nodes().String() +
		", edges: " +
		g.Edges().String()
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
