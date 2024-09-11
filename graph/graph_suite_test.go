package graph_test

import (
	"testing"

	"github.com/jbduncan/go-containers/graph"
	"github.com/jbduncan/go-containers/graph/graphtest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUndirectedGraph(t *testing.T) {
	graphtest.Graph(
		t,
		func() graph.Graph[int] {
			return graph.Undirected[int]().Build()
		},
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
		graphtest.Mutable,
		graphtest.Directed,
		graphtest.AllowsSelfLoops,
	)
}

func TestGraph(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Graph Suite")
}
