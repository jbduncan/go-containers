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
			allowsSelfLoops:    b.allowsSelfLoops,
			numEdges:           0,
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

func (u *undirectedGraph[N]) IsDirected() bool {
	return false
}

func (u *undirectedGraph[N]) AllowsSelfLoops() bool {
	return u.allowsSelfLoops
}

func (u *undirectedGraph[N]) Nodes() set.Set[N] {
	return keySet[N]{
		delegate: u.nodeToAdjacentNodes,
	}
}

func (u *undirectedGraph[N]) Edges() set.Set[EndpointPair[N]] {
	return edgeSet[N]{
		delegate:  u,
		edgeCount: func() int { return u.numEdges },
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
	return u.AdjacentNodes(node).Len()
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
	if _, ok := u.nodeToAdjacentNodes[node]; ok {
		return false
	}

	u.nodeToAdjacentNodes[node] = set.NewMutable[N]()
	return true
}

func (u *undirectedGraph[N]) PutEdge(source, target N) bool {
	if !u.AllowsSelfLoops() && source == target {
		panic("self-loops are disallowed")
	}

	added := u.putEdge(source, target)
	u.putEdge(target, source)

	if added {
		u.numEdges++
		return true
	}

	return false
}

func (u *undirectedGraph[N]) putEdge(source, target N) bool {
	added := false
	adjacentNodes, ok := u.nodeToAdjacentNodes[source]
	if !ok {
		adjacentNodes = set.NewMutable[N]()
		u.nodeToAdjacentNodes[source] = adjacentNodes
		added = true
	}
	if adjacentNodes.Add(target) {
		added = true
	}
	return added
}

func (u *undirectedGraph[N]) RemoveNode(node N) bool {
	adjacentNodes, ok := u.nodeToAdjacentNodes[node]
	if !ok {
		return false
	}

	delete(u.nodeToAdjacentNodes, node)

	for _, adjacentNodes := range u.nodeToAdjacentNodes {
		adjacentNodes.Remove(node)
	}

	u.numEdges -= adjacentNodes.Len()

	return true
}

func (u *undirectedGraph[N]) RemoveEdge(source, target N) bool {
	removedUToV := u.removeEdge(source, target)
	u.removeEdge(target, source)

	u.numEdges--

	return removedUToV
}

func (u *undirectedGraph[N]) removeEdge(from, to N) bool {
	if adjacentNodes, ok := u.nodeToAdjacentNodes[from]; ok {
		return adjacentNodes.Remove(to)
	}

	return false
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
		delegate:  d,
		edgeCount: func() int { return d.numEdges },
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
	d.nodes.Add(node)

	return false
}

func (d *directedGraph[N]) PutEdge(source, target N) bool {
	d.nodes.Add(source)
	d.nodes.Add(target)

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

	d.numEdges++

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
