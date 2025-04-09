package graph

import (
	"slices"
	"strconv"

	"github.com/jbduncan/go-containers/set"
)

type Graph[N comparable] interface {
	IsDirected() bool
	AllowsSelfLoops() bool
	Nodes() set.Set[N]
	Edges() set.Set[EndpointPair[N]]
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
		return &graph[N]{
			directed:        true,
			allowsSelfLoops: b.allowsSelfLoops,
			nodes:           set.Of[N](),
			connections: directedConnections[N]{
				nodeToPredecessors: make(map[N]*set.MapSet[N]),
				nodeToSuccessors:   make(map[N]*set.MapSet[N]),
			},
			numEdges: 0,
		}
	}

	return &graph[N]{
		directed:        false,
		allowsSelfLoops: b.allowsSelfLoops,
		nodes:           set.Of[N](),
		connections: undirectedConnections[N]{
			nodeToAdjacentNodes: make(map[N]*set.MapSet[N]),
		},
		numEdges: 0,
	}
}

var (
	_ Graph[int]        = (*graph[int])(nil)
	_ MutableGraph[int] = (*graph[int])(nil)
)

type graph[N comparable] struct {
	directed        bool
	allowsSelfLoops bool
	nodes           *set.MapSet[N]
	connections     connections[N]
	numEdges        int
}

func (g *graph[N]) IsDirected() bool {
	return g.directed
}

func (g *graph[N]) AllowsSelfLoops() bool {
	return g.allowsSelfLoops
}

func (g *graph[N]) Nodes() set.Set[N] {
	return set.Unmodifiable[N](g.nodes)
}

func (g *graph[N]) Edges() set.Set[EndpointPair[N]] {
	return edgeSet[N]{
		delegate: g,
		len:      func() int { return g.numEdges },
	}
}

func (g *graph[N]) AdjacentNodes(node N) set.Set[N] {
	if g.directed {
		return directedGraphAdjacentNodeSet[N]{
			node:     node,
			delegate: g,
		}
	}

	return g.Predecessors(node)
}

func (g *graph[N]) Predecessors(node N) set.Set[N] {
	return g.connections.Predecessors(node)
}

func (g *graph[N]) Successors(node N) set.Set[N] {
	return g.connections.Successors(node)
}

func (g *graph[N]) IncidentEdges(node N) set.Set[EndpointPair[N]] {
	return incidentEdgeSet[N]{
		node:     node,
		delegate: g,
	}
}

func (g *graph[N]) Degree(node N) int {
	if g.directed {
		return g.InDegree(node) + g.OutDegree(node)
	}

	selfLoop := g.AdjacentNodes(node).Contains(node)
	selfLoopCorrection := 0
	if selfLoop {
		selfLoopCorrection = 1
	}
	return g.AdjacentNodes(node).Len() + selfLoopCorrection
}

func (g *graph[N]) InDegree(node N) int {
	if g.directed {
		return g.Predecessors(node).Len()
	}

	return g.Degree(node)
}

func (g *graph[N]) OutDegree(node N) int {
	if g.directed {
		return g.Successors(node).Len()
	}

	return g.Degree(node)
}

func (g *graph[N]) HasEdgeConnecting(source N, target N) bool {
	return g.Successors(source).Contains(target)
}

func (g *graph[N]) HasEdgeConnectingEndpoints(endpointPair EndpointPair[N]) bool {
	return g.HasEdgeConnecting(endpointPair.Source(), endpointPair.Target())
}

func (g *graph[N]) String() string {
	return "isDirected: " +
		strconv.FormatBool(g.IsDirected()) +
		", allowsSelfLoops: " +
		strconv.FormatBool(g.AllowsSelfLoops()) +
		", nodes: " +
		g.Nodes().String() +
		", edges: " +
		g.Edges().String()
}

func (g *graph[N]) AddNode(node N) bool {
	return g.nodes.Add(node)
}

func (g *graph[N]) PutEdge(source N, target N) bool {
	if !g.AllowsSelfLoops() && source == target {
		panic("self-loops are disallowed")
	}

	g.AddNode(source)
	g.AddNode(target)

	put := g.connections.PutEdge(source, target)
	if put {
		g.numEdges++
		return true
	}

	return false
}

func (g *graph[N]) RemoveNode(node N) bool {
	if !g.Nodes().Contains(node) {
		return false
	}

	g.nodes.Remove(node)
	g.numEdges -= g.AdjacentNodes(node).Len()

	g.connections.RemoveNode(node)
	return true
}

func (g *graph[N]) RemoveEdge(source N, target N) bool {
	g.numEdges--
	return g.connections.RemoveEdge(source, target)
}

type connections[N comparable] interface {
	Predecessors(node N) set.Set[N]
	Successors(node N) set.Set[N]
	PutEdge(source N, target N) bool
	RemoveNode(node N)
	RemoveEdge(source N, target N) bool
}

type undirectedConnections[N comparable] struct {
	nodeToAdjacentNodes map[N]*set.MapSet[N]
}

func (u undirectedConnections[N]) adjacentNodes(node N) set.Set[N] {
	return neighborSet[N]{
		node:            node,
		nodeToNeighbors: u.nodeToAdjacentNodes,
	}
}

func (u undirectedConnections[N]) Predecessors(node N) set.Set[N] {
	return u.adjacentNodes(node)
}

func (u undirectedConnections[N]) Successors(node N) set.Set[N] {
	return u.adjacentNodes(node)
}

func (u undirectedConnections[N]) PutEdge(source N, target N) bool {
	put := putConnection(u.nodeToAdjacentNodes, source, target)
	putConnection(u.nodeToAdjacentNodes, target, source)
	return put
}

func (u undirectedConnections[N]) RemoveNode(node N) {
	for adjNode := range u.adjacentNodes(node).All() {
		removeConnection(u.nodeToAdjacentNodes, adjNode, node)
	}

	delete(u.nodeToAdjacentNodes, node)
}

func (u undirectedConnections[N]) RemoveEdge(source N, target N) bool {
	removed := removeConnection(u.nodeToAdjacentNodes, source, target)
	removeConnection(u.nodeToAdjacentNodes, target, source)
	return removed
}

type directedConnections[N comparable] struct {
	nodeToPredecessors map[N]*set.MapSet[N]
	nodeToSuccessors   map[N]*set.MapSet[N]
}

func (d directedConnections[N]) Predecessors(node N) set.Set[N] {
	return neighborSet[N]{
		node:            node,
		nodeToNeighbors: d.nodeToPredecessors,
	}
}

func (d directedConnections[N]) Successors(node N) set.Set[N] {
	return neighborSet[N]{
		node:            node,
		nodeToNeighbors: d.nodeToSuccessors,
	}
}

func (d directedConnections[N]) PutEdge(source N, target N) bool {
	put := putConnection(d.nodeToPredecessors, target, source)
	putConnection(d.nodeToSuccessors, source, target)
	return put
}

func (d directedConnections[N]) RemoveNode(node N) {
	successors := d.Successors(node)
	predecessors := copyOf(d.Predecessors(node))

	for successor := range successors.All() {
		removeConnection(d.nodeToPredecessors, successor, node)
	}
	delete(d.nodeToPredecessors, node)

	for _, predecessor := range predecessors {
		removeConnection(d.nodeToSuccessors, predecessor, node)
	}
	delete(d.nodeToSuccessors, node)
}

func (d directedConnections[N]) RemoveEdge(source N, target N) bool {
	removed := removeConnection(d.nodeToPredecessors, target, source)
	removeConnection(d.nodeToSuccessors, source, target)
	return removed
}

func copyOf[T comparable](s set.Set[T]) []T {
	return slices.Collect(s.All())
}

func putConnection[N comparable](nodeToNeighbors map[N]*set.MapSet[N], from, to N) bool {
	neighbors, ok := nodeToNeighbors[from]
	if !ok {
		neighbors = set.Of[N]()
		nodeToNeighbors[from] = neighbors
	}
	return neighbors.Add(to)
}

func removeConnection[N comparable](nodeToNeighbors map[N]*set.MapSet[N], from, to N) bool {
	neighbors, ok := nodeToNeighbors[from]
	if !ok {
		return false
	}

	removed := neighbors.Remove(to)
	if neighbors.Len() == 0 {
		delete(nodeToNeighbors, from)
	}
	return removed
}
