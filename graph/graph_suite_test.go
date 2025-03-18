package graph_test

import (
	"testing"

	"github.com/jbduncan/go-containers/graph"
	"github.com/jbduncan/go-containers/graph/graphtest"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func addNodeToMutableGraph(g graph.Graph[int], node int) graph.Graph[int] {
	g.(graph.MutableGraph[int]).AddNode(node)
	return g
}

func putEdgeOnMutableGraph(
	g graph.Graph[int],
	source int,
	target int,
) graph.Graph[int] {
	g.(graph.MutableGraph[int]).PutEdge(source, target)
	return g
}

func TestUndirectedGraph(t *testing.T) {
	graphtest.Graph(
		t,
		func() graph.Graph[int] {
			return graph.Undirected[int]().Build()
		},
		addNodeToMutableGraph,
		putEdgeOnMutableGraph,
		graphtest.Mutable,
		graphtest.Undirected,
		graphtest.DisallowsSelfLoops,
	)
}

func TestUndirectedAllowsSelfLoopsGraph(t *testing.T) {
	graphtest.Graph(
		t,
		func() graph.Graph[int] {
			return graph.Undirected[int]().AllowsSelfLoops(true).Build()
		},
		addNodeToMutableGraph,
		putEdgeOnMutableGraph,
		graphtest.Mutable,
		graphtest.Undirected,
		graphtest.AllowsSelfLoops,
	)
}

func TestDirectedGraph(t *testing.T) {
	graphtest.Graph(
		t,
		func() graph.Graph[int] {
			return graph.Directed[int]().Build()
		},
		addNodeToMutableGraph,
		putEdgeOnMutableGraph,
		graphtest.Mutable,
		graphtest.Directed,
		graphtest.DisallowsSelfLoops,
	)
}

func TestDirectedAllowsSelfLoopsGraph(t *testing.T) {
	graphtest.Graph(
		t,
		func() graph.Graph[int] {
			return graph.Directed[int]().AllowsSelfLoops(true).Build()
		},
		addNodeToMutableGraph,
		putEdgeOnMutableGraph,
		graphtest.Mutable,
		graphtest.Directed,
		graphtest.AllowsSelfLoops,
	)
}

func TestGraph(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Graph Suite")
}
