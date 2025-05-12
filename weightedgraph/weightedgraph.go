package weightedgraph

func Undirected[N comparable]() Builder[N] {
	return Builder[N]{
		directed:        false,
		allowsSelfLoops: false,
	}
}

type Builder[N comparable] struct {
	directed        bool
	allowsSelfLoops bool
}

func (b Builder[N]) Build() *WeightedGraph[N] {
	return &WeightedGraph[N]{}
}

type WeightedGraph[N comparable] struct{}
