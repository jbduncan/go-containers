package graph

import "github.com/jbduncan/go-containers/set"

type directedGraphAdjacentNodeSet[N comparable] struct {
	node     N
	delegate *directedGraph[N]
}

func (p directedGraphAdjacentNodeSet[N]) Contains(elem N) bool {
	return p.union().Contains(elem)
}

func (p directedGraphAdjacentNodeSet[N]) Len() int {
	selfLoop := p.Contains(p.node)
	selfLoopCorrection := 0
	if selfLoop {
		selfLoopCorrection = 1
	}

	return p.delegate.InDegree(p.node) +
		p.delegate.OutDegree(p.node) -
		selfLoopCorrection
}

func (p directedGraphAdjacentNodeSet[N]) ForEach(fn func(node N)) {
	p.union().ForEach(fn)
}

func (p directedGraphAdjacentNodeSet[N]) String() string {
	return p.union().String()
}

func (p directedGraphAdjacentNodeSet[N]) union() set.Set[N] {
	return set.Union[N](
		p.delegate.Predecessors(p.node),
		p.delegate.Successors(p.node),
	)
}
