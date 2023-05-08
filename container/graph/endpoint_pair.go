package graph

import "fmt"

// TODO: Docs

// TODO: Document that an unordered endpoint pair and its reverse are equal to each other,
//       and repeat again in the docs for EndpointPair.Equal().

func NewUnorderedEndpointPair[N comparable](nodeU N, nodeV N) EndpointPair[N] {
	return EndpointPair[N]{
		nodeU:     nodeU,
		nodeV:     nodeV,
		isOrdered: false,
	}
}

func NewOrderedEndpointPair[N comparable](source N, target N) EndpointPair[N] {
	return EndpointPair[N]{
		nodeU:     source,
		nodeV:     target,
		isOrdered: true,
	}
}

type EndpointPair[N comparable] struct {
	nodeU     N
	nodeV     N
	isOrdered bool
}

func (e EndpointPair[N]) Source() N {
	if !e.isOrdered {
		panic(notAvailableOnUndirected)
	}

	return e.nodeU
}

func (e EndpointPair[N]) Target() N {
	if !e.isOrdered {
		panic(notAvailableOnUndirected)
	}

	return e.nodeV
}

const notAvailableOnUndirected = "cannot call Source()/Target() on an EndpointPair from an " +
	"undirected graph; consider calling AdjacentNode(node) if you already have a node, or " +
	"NodeU()/NodeV() if you don't"

func (e EndpointPair[N]) NodeU() N {
	return e.nodeU
}

func (e EndpointPair[N]) NodeV() N {
	return e.nodeV
}

func (e EndpointPair[N]) AdjacentNode(node N) N {
	if node == e.NodeU() {
		return e.NodeV()
	}
	if node == e.NodeV() {
		return e.NodeU()
	}

	panic(fmt.Sprintf("EndpointPair %s does not contain node %v", e.String(), node))
}

func (e EndpointPair[N]) IsOrdered() bool {
	return e.isOrdered
}

// TODO: Document that == is discouraged.

func (e EndpointPair[N]) Equal(other EndpointPair[N]) bool {
	if e.IsOrdered() && other.IsOrdered() {
		return e.Source() == other.Source() &&
			e.Target() == other.Target()
	}

	if !e.IsOrdered() && !other.IsOrdered() {
		return e.equalNodes(other) ||
			e.equalNodesInReverse(other)
	}

	return false
}

func (e EndpointPair[N]) equalNodes(other EndpointPair[N]) bool {
	return e.NodeU() == other.NodeU() && e.NodeV() == other.NodeV()
}

func (e EndpointPair[N]) equalNodesInReverse(other EndpointPair[N]) bool {
	return e.NodeU() == other.NodeV() && e.NodeV() == other.NodeU()
}

func (e EndpointPair[N]) String() string {
	if e.isOrdered {
		return fmt.Sprintf("<%v -> %v>", e.Source(), e.Target())
	}
	return fmt.Sprintf("[%v, %v]", e.NodeU(), e.NodeV())
}
