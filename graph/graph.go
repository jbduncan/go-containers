package graph

import (
	"slices"
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
			nodes:              set.Of[N](),
			nodeToPredecessors: make(map[N]set.MutableSet[N]),
			nodeToSuccessors:   make(map[N]set.MutableSet[N]),
			allowsSelfLoops:    b.allowsSelfLoops,
			numEdges:           0,
		}
	}

	return &undirectedGraph[N]{
		nodes:               set.Of[N](),
		nodeToAdjacentNodes: make(map[N]set.MutableSet[N]),
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
	nodes               set.MutableSet[N]
	nodeToAdjacentNodes map[N]set.MutableSet[N]
	allowsSelfLoops     bool
	numEdges            int
}

func (u *undirectedGraph[N]) IsDirected() bool {
	return false
}

func (u *undirectedGraph[N]) AllowsSelfLoops() bool {
	return u.allowsSelfLoops
}

func (u *undirectedGraph[N]) Nodes() set.Set[N] {
	return set.Unmodifiable(u.nodes)
}

func (u *undirectedGraph[N]) Edges() set.Set[EndpointPair[N]] {
	return edgeSet[N]{
		delegate: u,
		len:      func() int { return u.numEdges },
	}
}

func (u *undirectedGraph[N]) AdjacentNodes(node N) set.Set[N] {
	return neighborSet[N]{
		node:            node,
		nodeToNeighbors: u.nodeToAdjacentNodes,
	}
}

func (u *undirectedGraph[N]) Predecessors(node N) set.Set[N] {
	return u.AdjacentNodes(node)
}

func (u *undirectedGraph[N]) Successors(node N) set.Set[N] {
	return u.AdjacentNodes(node)
}

func (u *undirectedGraph[N]) IncidentEdges(node N) set.Set[EndpointPair[N]] {
	return incidentEdgeSet[N]{
		node:     node,
		delegate: u,
	}
}

func (u *undirectedGraph[N]) Degree(node N) int {
	selfLoop := u.AdjacentNodes(node).Contains(node)
	selfLoopCorrection := 0
	if selfLoop {
		selfLoopCorrection = 1
	}

	return u.AdjacentNodes(node).Len() + selfLoopCorrection
}

func (u *undirectedGraph[N]) InDegree(node N) int {
	return u.Degree(node)
}

func (u *undirectedGraph[N]) OutDegree(node N) int {
	return u.Degree(node)
}

func (u *undirectedGraph[N]) HasEdgeConnecting(source, target N) bool {
	return hasEdgeConnecting[N](u, source, target)
}

func (u *undirectedGraph[N]) HasEdgeConnectingEndpoints(endpointPair EndpointPair[N]) bool {
	return hasEdgeConnectingEndpoints[N](u, endpointPair)
}

func (u *undirectedGraph[N]) String() string {
	return stringOf[N](u)
}

func (u *undirectedGraph[N]) AddNode(node N) bool {
	return u.nodes.Add(node)
}

func (u *undirectedGraph[N]) PutEdge(source, target N) bool {
	checkSelfLoop[N](u, source, target)

	u.AddNode(source)
	u.AddNode(target)

	put := putConnection(u.nodeToAdjacentNodes, source, target)
	putConnection(u.nodeToAdjacentNodes, target, source)

	if put {
		u.numEdges++
		return true
	}

	return false
}

func (u *undirectedGraph[N]) RemoveNode(node N) bool {
	if !u.Nodes().Contains(node) {
		return false
	}

	u.nodes.Remove(node)
	u.numEdges -= u.AdjacentNodes(node).Len()

	for adjNode := range u.AdjacentNodes(node).All() {
		removeConnection(u.nodeToAdjacentNodes, adjNode, node)
	}

	delete(u.nodeToAdjacentNodes, node)

	return true
}

func (u *undirectedGraph[N]) RemoveEdge(source, target N) bool {
	removed := removeConnection(u.nodeToAdjacentNodes, source, target)
	removeConnection(u.nodeToAdjacentNodes, target, source)

	u.numEdges--

	return removed
}

type directedGraph[N comparable] struct {
	nodes              set.MutableSet[N]
	nodeToPredecessors map[N]set.MutableSet[N]
	nodeToSuccessors   map[N]set.MutableSet[N]
	allowsSelfLoops    bool
	numEdges           int
}

func (d *directedGraph[N]) Nodes() set.Set[N] {
	return set.Unmodifiable(d.nodes)
}

func (d *directedGraph[N]) Edges() set.Set[EndpointPair[N]] {
	return edgeSet[N]{
		delegate: d,
		len:      func() int { return d.numEdges },
	}
}

func (d *directedGraph[N]) IsDirected() bool {
	return true
}

func (d *directedGraph[N]) AllowsSelfLoops() bool {
	return d.allowsSelfLoops
}

func (d *directedGraph[N]) AdjacentNodes(node N) set.Set[N] {
	return directedGraphAdjacentNodeSet[N]{
		node:     node,
		delegate: d,
	}
}

func (d *directedGraph[N]) Predecessors(node N) set.Set[N] {
	return neighborSet[N]{
		node:            node,
		nodeToNeighbors: d.nodeToPredecessors,
	}
}

func (d *directedGraph[N]) Successors(node N) set.Set[N] {
	return neighborSet[N]{
		node:            node,
		nodeToNeighbors: d.nodeToSuccessors,
	}
}

func (d *directedGraph[N]) IncidentEdges(node N) set.Set[EndpointPair[N]] {
	return incidentEdgeSet[N]{
		node:     node,
		delegate: d,
	}
}

func (d *directedGraph[N]) Degree(node N) int {
	return d.InDegree(node) + d.OutDegree(node)
}

func (d *directedGraph[N]) InDegree(node N) int {
	return d.Predecessors(node).Len()
}

func (d *directedGraph[N]) OutDegree(node N) int {
	return d.Successors(node).Len()
}

func (d *directedGraph[N]) HasEdgeConnecting(source N, target N) bool {
	return hasEdgeConnecting[N](d, source, target)
}

func (d *directedGraph[N]) HasEdgeConnectingEndpoints(endpointPair EndpointPair[N]) bool {
	return hasEdgeConnectingEndpoints[N](d, endpointPair)
}

func (d *directedGraph[N]) String() string {
	return stringOf[N](d)
}

func (d *directedGraph[N]) AddNode(node N) bool {
	return d.nodes.Add(node)
}

func (d *directedGraph[N]) PutEdge(source, target N) bool {
	checkSelfLoop[N](d, source, target)

	d.nodes.Add(source)
	d.nodes.Add(target)

	put := putConnection(d.nodeToPredecessors, target, source)
	putConnection(d.nodeToSuccessors, source, target)
	if put {
		d.numEdges++
		return true
	}

	return false
}

func (d *directedGraph[N]) RemoveNode(node N) bool {
	if !d.Nodes().Contains(node) {
		return false
	}

	d.nodes.Remove(node)
	d.numEdges -= d.AdjacentNodes(node).Len()

	successors := d.Successors(node)
	predecessors := copyOf(d.Predecessors(node))

	for successor := range successors.All() {
		removeConnection(d.nodeToPredecessors, successor, node)
	}
	delete(d.nodeToPredecessors, node)

	for predecessor := range predecessors.All() {
		removeConnection(d.nodeToSuccessors, predecessor, node)
	}
	delete(d.nodeToSuccessors, node)

	return true
}

func copyOf[T comparable](s set.Set[T]) set.Set[T] {
	return set.Of[T](slices.Collect(s.All())...)
}

func (d *directedGraph[N]) RemoveEdge(source, target N) bool {
	removed := removeConnection(d.nodeToPredecessors, target, source)
	removeConnection(d.nodeToSuccessors, source, target)

	d.numEdges--

	return removed
}

func hasEdgeConnecting[N comparable](g Graph[N], source N, target N) bool {
	return g.Successors(source).Contains(target)
}

func hasEdgeConnectingEndpoints[N comparable](g Graph[N], endpointPair EndpointPair[N]) bool {
	return hasEdgeConnecting(g, endpointPair.Source(), endpointPair.Target())
}

func stringOf[N comparable](g Graph[N]) string {
	return "isDirected: " +
		strconv.FormatBool(g.IsDirected()) +
		", allowsSelfLoops: " +
		strconv.FormatBool(g.AllowsSelfLoops()) +
		", nodes: " +
		g.Nodes().String() +
		", edges: " +
		g.Edges().String()
}

func checkSelfLoop[N comparable](g Graph[N], source, target N) {
	if !g.AllowsSelfLoops() && source == target {
		panic("self-loops are disallowed")
	}
}

func putConnection[N comparable](nodeToNeighbors map[N]set.MutableSet[N], from, to N) bool {
	neighbors, ok := nodeToNeighbors[from]
	if !ok {
		neighbors = set.Of[N]()
		nodeToNeighbors[from] = neighbors
	}
	return neighbors.Add(to)
}

func removeConnection[N comparable](nodeToNeighbors map[N]set.MutableSet[N], from, to N) bool {
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
