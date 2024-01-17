package graph

import "fmt"

func EndpointPairOf[N comparable](source N, target N) EndpointPair[N] {
	return EndpointPair[N]{
		source: source,
		target: target,
	}
}

type EndpointPair[N comparable] struct {
	source N
	target N
}

func (e EndpointPair[N]) Source() N {
	return e.source
}

func (e EndpointPair[N]) Target() N {
	return e.target
}

func (e EndpointPair[N]) AdjacentNode(node N) N {
	switch node {
	case e.Source():
		return e.Target()
	case e.Target():
		return e.Source()
	default:
		panic(fmt.Sprintf("EndpointPair %s does not contain node %v", e.String(), node))
	}
}

func (e EndpointPair[N]) String() string {
	return fmt.Sprintf("<%v -> %v>", e.Source(), e.Target())
}
