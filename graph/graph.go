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
		return &directedGraph[N]{
			nodes:              set.NewMutable[N](),
			nodeToPredecessors: map[N]set.MutableSet[N]{},
			nodeToSuccessors:   map[N]set.MutableSet[N]{},
		}
	}

	return &undirectedGraph[N]{
		nodeToAdjacentNodes: map[N]set.MutableSet[N]{},
		allowsSelfLoops:     b.allowsSelfLoops,
		numEdges:            0,
	}
}

var (
	_ Graph[int]        = (*undirectedGraph[int])(nil)
	_ MutableGraph[int] = (*undirectedGraph[int])(nil)
	_ Graph[int]        = (*directedGraph[int])(nil)
	_ MutableGraph[int] = (*directedGraph[int])(nil)
)

type undirectedGraph[N comparable] struct {
	nodeToAdjacentNodes map[N]set.MutableSet[N]
	allowsSelfLoops     bool
	numEdges            int
}

func (g *undirectedGraph[N]) IsDirected() bool {
	return false
}

func (g *undirectedGraph[N]) AllowsSelfLoops() bool {
	return g.allowsSelfLoops
}

func (g *undirectedGraph[N]) Nodes() set.Set[N] {
	return keySet[N]{
		delegate: g.nodeToAdjacentNodes,
	}
}

func (g *undirectedGraph[N]) Edges() set.Set[EndpointPair[N]] {
	return edgeSet[N]{
		delegate: g,
	}
}

func (g *undirectedGraph[N]) AdjacentNodes(node N) set.Set[N] {
	return adjacentNodeSet[N]{
		node:                node,
		nodeToAdjacentNodes: g.nodeToAdjacentNodes,
	}
}

func (g *undirectedGraph[N]) Predecessors(node N) set.Set[N] {
	return g.AdjacentNodes(node)
}

func (g *undirectedGraph[N]) Successors(node N) set.Set[N] {
	return g.AdjacentNodes(node)
}

func (g *undirectedGraph[N]) IncidentEdges(node N) set.Set[EndpointPair[N]] {
	return incidentEdgeSet[N]{
		node:          node,
		adjacencyList: g.nodeToAdjacentNodes,
	}
}

func (g *undirectedGraph[N]) Degree(node N) int {
	return g.AdjacentNodes(node).Len()
}

func (g *undirectedGraph[N]) InDegree(node N) int {
	return g.Degree(node)
}

func (g *undirectedGraph[N]) OutDegree(node N) int {
	return g.Degree(node)
}

func (g *undirectedGraph[N]) HasEdgeConnecting(source, target N) bool {
	return g.AdjacentNodes(source).Contains(target)
}

func (g *undirectedGraph[N]) HasEdgeConnectingEndpoints(endpointPair EndpointPair[N]) bool {
	return g.HasEdgeConnecting(endpointPair.Source(), endpointPair.Target())
}

func (g *undirectedGraph[N]) String() string {
	return "isDirected: false, allowsSelfLoops: " +
		strconv.FormatBool(g.allowsSelfLoops) +
		", nodes: " +
		g.Nodes().String() +
		", edges: " +
		g.Edges().String()
}

func (g *undirectedGraph[N]) AddNode(node N) bool {
	if _, ok := g.nodeToAdjacentNodes[node]; ok {
		return false
	}

	g.nodeToAdjacentNodes[node] = set.NewMutable[N]()
	return true
}

func (g *undirectedGraph[N]) PutEdge(source, target N) bool {
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

func (g *undirectedGraph[N]) putEdge(source, target N) bool {
	added := false
	adjacentNodes, ok := g.nodeToAdjacentNodes[source]
	if !ok {
		adjacentNodes = set.NewMutable[N]()
		g.nodeToAdjacentNodes[source] = adjacentNodes
		added = true
	}
	if adjacentNodes.Add(target) {
		added = true
	}
	return added
}

func (g *undirectedGraph[N]) RemoveNode(node N) bool {
	adjacentNodes, ok := g.nodeToAdjacentNodes[node]
	if !ok {
		return false
	}

	delete(g.nodeToAdjacentNodes, node)

	for _, adjacentNodes := range g.nodeToAdjacentNodes {
		adjacentNodes.Remove(node)
	}

	g.numEdges -= adjacentNodes.Len()

	return true
}

func (g *undirectedGraph[N]) RemoveEdge(source, target N) bool {
	removedUToV := g.removeEdge(source, target)
	g.removeEdge(target, source)

	g.numEdges--

	return removedUToV
}

func (g *undirectedGraph[N]) removeEdge(from, to N) bool {
	adjacentNodes, ok := g.nodeToAdjacentNodes[from]
	if !ok {
		return false
	}

	return adjacentNodes.Remove(to)
}

type directedGraph[N comparable] struct {
	nodes              set.MutableSet[N]
	nodeToPredecessors map[N]set.MutableSet[N]
	nodeToSuccessors   map[N]set.MutableSet[N]
}

func (d *directedGraph[N]) Nodes() set.Set[N] {
	return d.nodes
}

func (d *directedGraph[N]) Edges() set.Set[EndpointPair[N]] {
	return set.Of[EndpointPair[N]]()
}

func (d *directedGraph[N]) IsDirected() bool {
	return true
}

func (d *directedGraph[N]) AllowsSelfLoops() bool {
	return false
}

//nolint:revive
func (d *directedGraph[N]) AdjacentNodes(node N) set.Set[N] {
	return set.Of[N]()
}

func (d *directedGraph[N]) Predecessors(node N) set.Set[N] {
	return adjacentNodeSet[N]{
		node:                node,
		nodeToAdjacentNodes: d.nodeToPredecessors,
	}
}

func (d *directedGraph[N]) Successors(node N) set.Set[N] {
	return adjacentNodeSet[N]{
		node:                node,
		nodeToAdjacentNodes: d.nodeToSuccessors,
	}
}

//nolint:revive
func (d *directedGraph[N]) IncidentEdges(node N) set.Set[EndpointPair[N]] {
	return set.Of[EndpointPair[N]]()
}

//nolint:revive
func (d *directedGraph[N]) Degree(node N) int {
	return 0
}

//nolint:revive
func (d *directedGraph[N]) InDegree(node N) int {
	return 0
}

//nolint:revive
func (d *directedGraph[N]) OutDegree(node N) int {
	return 0
}

//nolint:revive
func (d *directedGraph[N]) HasEdgeConnecting(source N, target N) bool {
	panic("implement me")
}

//nolint:revive
func (d *directedGraph[N]) HasEdgeConnectingEndpoints(endpointPair EndpointPair[N]) bool {
	panic("implement me")
}

func (d *directedGraph[N]) String() string {
	panic("implement me")
}

func (d *directedGraph[N]) AddNode(node N) bool {
	d.nodes.Add(node)

	return false
}

func (d *directedGraph[N]) PutEdge(source, target N) bool {
	predecessors, ok := d.nodeToPredecessors[target]
	if !ok {
		predecessors = set.NewMutable[N]()
		d.nodeToPredecessors[target] = predecessors
	}
	predecessors.Add(source)

	successors, ok := d.nodeToSuccessors[source]
	if !ok {
		successors = set.NewMutable[N]()
		d.nodeToSuccessors[source] = successors
	}
	successors.Add(target)

	return false
}

//nolint:revive
func (d *directedGraph[N]) RemoveNode(node N) bool {
	panic("implement me")
}

//nolint:revive
func (d *directedGraph[N]) RemoveEdge(source, target N) bool {
	panic("implement me")
}
