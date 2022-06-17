package graph

import "fmt"

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

// TODO: EndpointPair: make Equals method and discourage == from being used (documenting that its use is undefined).
//       See this link:
//       https://github.com/google/guava/blob/4d323b2b117a5906ab16074c8c88b4ff162b1b82/guava/src/com/google/common/graph/EndpointPair.java#L131-L145

func (e EndpointPair[N]) String() string {
	if e.isOrdered {
		return fmt.Sprintf("<%v -> %v>", e.Source(), e.Target())
	}
	return fmt.Sprintf("[%v, %v]", e.NodeU(), e.NodeV())
}
