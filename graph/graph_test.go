package graph_test

import (
	"testing"

	"github.com/jbduncan/go-containers/graph"
	"github.com/jbduncan/go-containers/graph/graphtest"
)

func TestUndirectedGraph(t *testing.T) {
	t.Parallel()

	graphtest.TestMutable(
		t,
		func() graphtest.MutableGraph[int] {
			return graph.Undirected[int]().Build()
		},
		graphtest.Undirected,
		graphtest.DisallowsSelfLoops,
	)
}

func TestUndirectedAllowsSelfLoopsGraph(t *testing.T) {
	t.Parallel()

	graphtest.TestMutable(
		t,
		func() graphtest.MutableGraph[int] {
			return graph.Undirected[int]().AllowsSelfLoops(true).Build()
		},
		graphtest.Undirected,
		graphtest.AllowsSelfLoops,
	)
}

func TestDirectedGraph(t *testing.T) {
	t.Parallel()

	graphtest.TestMutable(
		t,
		func() graphtest.MutableGraph[int] {
			return graph.Directed[int]().Build()
		},
		graphtest.Directed,
		graphtest.DisallowsSelfLoops,
	)
}

func TestDirectedAllowsSelfLoopsGraph(t *testing.T) {
	t.Parallel()

	graphtest.TestMutable(
		t,
		func() graphtest.MutableGraph[int] {
			return graph.Directed[int]().AllowsSelfLoops(true).Build()
		},
		graphtest.Directed,
		graphtest.AllowsSelfLoops,
	)
}
