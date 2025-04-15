package graph

import (
	"iter"

	"github.com/jbduncan/go-containers/set"
)

type directedGraphAdjacentNodeSet[N comparable] struct {
	node     N
	delegate *Mutable[N]
}

func (p directedGraphAdjacentNodeSet[N]) Contains(element N) bool {
	return p.union().Contains(element)
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

func (p directedGraphAdjacentNodeSet[N]) All() iter.Seq[N] {
	return p.union().All()
}

func (p directedGraphAdjacentNodeSet[N]) String() string {
	return p.union().String()
}

func (p directedGraphAdjacentNodeSet[N]) union() SetView[N] {
	return set.Union[N](
		p.delegate.Predecessors(p.node),
		p.delegate.Successors(p.node),
	)
}
