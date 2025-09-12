package graph

import (
	"iter"
	"slices"
	"strconv"

	"github.com/jbduncan/go-containers/set"
)

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

func (b Builder[N]) Build() *Graph[N] {
	if b.directed {
		return &Graph[N]{
			directed:        true,
			allowsSelfLoops: b.allowsSelfLoops,
			nodes:           set.Of[N](),
			connections: directedConnections[N]{
				nodeToPredecessors: make(map[N]set.Set[N]),
				nodeToSuccessors:   make(map[N]set.Set[N]),
			},
			numEdges: 0,
		}
	}

	return &Graph[N]{
		directed:        false,
		allowsSelfLoops: b.allowsSelfLoops,
		nodes:           set.Of[N](),
		connections: undirectedConnections[N]{
			nodeToAdjacentNodes: make(map[N]set.Set[N]),
		},
		numEdges: 0,
	}
}

// SetView is a read-only set view; a generic, unordered collection of unique
// elements that only allows read operations. It is returned by a few of the
// methods on graph.Graph.
type SetView[T comparable] interface {
	// Contains returns true if this set contains the given element, otherwise
	// it returns false.
	Contains(element T) bool

	// Len returns the number of elements in this set.
	Len() int

	// All returns an iter.Seq that returns each and every element in this set.
	//
	// The iteration order is undefined; it may even change from one call to
	// the next.
	All() iter.Seq[T]

	// String returns a string representation of all the elements in this set.
	//
	// The format of this string is a single "[" followed by a comma-separated
	// list (", ") of this set's elements in the same order as All (which is
	// undefined and may change from one call to the next), followed by a
	// single "]".
	//
	// This method satisfies fmt.Stringer.
	String() string
}

type Graph[N comparable] struct {
	connections     connections[N]
	nodes           set.Set[N]
	numEdges        int
	directed        bool
	allowsSelfLoops bool
}

func (g *Graph[N]) IsDirected() bool {
	return g.directed
}

func (g *Graph[N]) AllowsSelfLoops() bool {
	return g.allowsSelfLoops
}

func (g *Graph[N]) Nodes() SetView[N] {
	return set.Unmodifiable[N](g.nodes)
}

func (g *Graph[N]) Edges() SetView[EndpointPair[N]] {
	return edgeSet[N]{
		delegate: g,
		len:      func() int { return g.numEdges },
	}
}

func (g *Graph[N]) AdjacentNodes(node N) SetView[N] {
	if g.directed {
		return directedGraphAdjacentNodeSet[N]{
			node:     node,
			delegate: g,
		}
	}

	return g.Predecessors(node)
}

func (g *Graph[N]) Predecessors(node N) SetView[N] {
	return g.connections.Predecessors(node)
}

func (g *Graph[N]) Successors(node N) SetView[N] {
	return g.connections.Successors(node)
}

func (g *Graph[N]) IncidentEdges(node N) SetView[EndpointPair[N]] {
	return incidentEdgeSet[N]{
		node:     node,
		delegate: g,
	}
}

func (g *Graph[N]) Degree(node N) int {
	if g.directed {
		return g.InDegree(node) + g.OutDegree(node)
	}

	return g.undirectedGraphDegree(node)
}

func (g *Graph[N]) InDegree(node N) int {
	if g.directed {
		return g.Predecessors(node).Len()
	}

	return g.undirectedGraphDegree(node)
}

func (g *Graph[N]) OutDegree(node N) int {
	if g.directed {
		return g.Successors(node).Len()
	}

	return g.undirectedGraphDegree(node)
}

func (g *Graph[N]) undirectedGraphDegree(node N) int {
	selfLoop := g.AdjacentNodes(node).Contains(node)
	selfLoopCorrection := 0
	if selfLoop {
		selfLoopCorrection = 1
	}
	return g.AdjacentNodes(node).Len() + selfLoopCorrection
}

func (g *Graph[N]) HasEdgeConnecting(source N, target N) bool {
	return g.Successors(source).Contains(target)
}

func (g *Graph[N]) HasEdgeConnectingEndpoints(endpointPair EndpointPair[N]) bool {
	return g.HasEdgeConnecting(endpointPair.Source(), endpointPair.Target())
}

func (g *Graph[N]) String() string {
	return "isDirected: " +
		strconv.FormatBool(g.IsDirected()) +
		", allowsSelfLoops: " +
		strconv.FormatBool(g.AllowsSelfLoops()) +
		", nodes: " +
		g.Nodes().String() +
		", edges: " +
		g.Edges().String()
}

func (g *Graph[N]) AddNode(node N) bool {
	return g.nodes.Add(node)
}

func (g *Graph[N]) PutEdge(source N, target N) bool {
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

func (g *Graph[N]) RemoveNode(node N) bool {
	if !g.Nodes().Contains(node) {
		return false
	}

	g.nodes.Remove(node)
	g.numEdges -= g.AdjacentNodes(node).Len()

	g.connections.RemoveNode(node)
	return true
}

func (g *Graph[N]) RemoveEdge(source N, target N) bool {
	g.numEdges--
	return g.connections.RemoveEdge(source, target)
}

type connections[N comparable] interface {
	Predecessors(node N) SetView[N]
	Successors(node N) SetView[N]
	PutEdge(source N, target N) bool
	RemoveNode(node N)
	RemoveEdge(source N, target N) bool
}

type undirectedConnections[N comparable] struct {
	nodeToAdjacentNodes map[N]set.Set[N]
}

func (u undirectedConnections[N]) adjacentNodes(node N) SetView[N] {
	return neighborSet[N]{
		node:            node,
		nodeToNeighbors: u.nodeToAdjacentNodes,
	}
}

func (u undirectedConnections[N]) Predecessors(node N) SetView[N] {
	return u.adjacentNodes(node)
}

func (u undirectedConnections[N]) Successors(node N) SetView[N] {
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
	nodeToPredecessors map[N]set.Set[N]
	nodeToSuccessors   map[N]set.Set[N]
}

func (d directedConnections[N]) Predecessors(node N) SetView[N] {
	return neighborSet[N]{
		node:            node,
		nodeToNeighbors: d.nodeToPredecessors,
	}
}

func (d directedConnections[N]) Successors(node N) SetView[N] {
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

func copyOf[T comparable](s SetView[T]) []T {
	return slices.Collect(s.All())
}

func putConnection[N comparable](nodeToNeighbors map[N]set.Set[N], from, to N) bool {
	neighbors, ok := nodeToNeighbors[from]
	if !ok {
		neighbors = set.Of[N]()
		nodeToNeighbors[from] = neighbors
	}
	return neighbors.Add(to)
}

func removeConnection[N comparable](nodeToNeighbors map[N]set.Set[N], from, to N) bool {
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
