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
	HasEdgeConnecting(source N, target N) bool
	HasEdgeConnectingEndpoints(endpointPair EndpointPair[N]) bool
	String() string
}

type MutableGraph[N comparable] interface {
	Graph[N]

	AddNode(node N) bool
	PutEdge(source N, target N) bool
	RemoveNode(node N) bool
	RemoveEdge(source N, target N) bool
}

func Undirected[N comparable]() Builder[N] {
	return Builder[N]{
		directed:        false,
		allowsSelfLoops: false,
	}
}

func Directed[N comparable]() Builder[N] {
	return Builder[N]{
		directed:        true,
		allowsSelfLoops: false,
	}
}

type Builder[N comparable] struct {
	directed        bool
	allowsSelfLoops bool
}

func (b Builder[N]) AllowsSelfLoops(allowsSelfLoops bool) Builder[N] {
	b.allowsSelfLoops = allowsSelfLoops
	return b
}

func (b Builder[N]) Build() MutableGraph[N] {
	if b.directed {
		panic("unimplemented")
	}

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

func (g *graph[N]) HasEdgeConnecting(source, target N) bool {
	if adjacentNodes, ok := g.adjacencyList[source]; ok {
		return adjacentNodes.Contains(target)
	}

	return false
}

func (g *graph[N]) HasEdgeConnectingEndpoints(endpointPair EndpointPair[N]) bool {
	return g.HasEdgeConnecting(endpointPair.Source(), endpointPair.Target())
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

func (g *graph[N]) PutEdge(source, target N) bool {
	if !g.AllowsSelfLoops() && source == target {
		panic("self-loops are disallowed")
	}

	addedUToV := g.putEdge(source, target)
	g.putEdge(target, source)

	if addedUToV {
		g.numEdges++
		return true
	}

	return false
}

func (g *graph[N]) putEdge(source, target N) bool {
	added := false
	adjacentNodes, ok := g.adjacencyList[source]
	if !ok {
		adjacentNodes = set.NewMutable[N]()
		g.adjacencyList[source] = adjacentNodes
		added = true
	}
	if adjacentNodes.Add(target) {
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

func (g *graph[N]) RemoveEdge(source, target N) bool {
	removedUToV := g.removeEdge(source, target)
	g.removeEdge(target, source)

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
