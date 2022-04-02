package container

type Graph[N comparable] interface {
	Nodes() Set[N]
	Edges() Set[EndpointPair[N]]
	IsDirected() bool
	AllowsSelfLoops() bool
	// TODO: Implement
	// NodeOrder() ElementOrder
	// IncidentEdgeOrder() ElementOrder
	AdjacentNodes() Set[N]
	Predecessors() Set[N]
	Successors() Set[N]
	IncidentEdges(node N) Set[EndpointPair[N]]
	Degree(node N) int
	InDegree(node N) int
	OutDegree(node N) int
	HasEdgeConnecting(nodeU N, nodeV N) bool
	HasEdgeConnectingEndpoints(endpointPair EndpointPair[N]) bool
	String() string
	// TODO: Is an Equals function needed to meet Guava's Graph::equals rules?
}

type MutableGraph[N comparable] interface {
	Graph[N]

	AddNode(n N) bool
	PutEdge(u N, v N) bool
	PutEdgeWithEndpoints(e EndpointPair[N]) bool
	RemoveNode(n N) bool
	RemoveEdge(u N, v N) bool
	RemoveEdgeWithEndpoints(e EndpointPair[N]) bool
}

// TODO: Implement
// type ElementOrder struct {}
