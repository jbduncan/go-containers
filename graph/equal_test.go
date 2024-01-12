package graph_test

import (
	"fmt"
	"testing"

	"github.com/jbduncan/go-containers/graph"
	"github.com/jbduncan/go-containers/set"
	. "github.com/onsi/gomega"
)

func TestEqual(t *testing.T) {
	t.Run("graph a: [link]; graph b: [link]; equal", func(t *testing.T) {
		g := NewWithT(t)
		a := undirectedGraphOf("link")

		g.Expect(graph.Equal(a, a)).To(BeTrue())
	})

	t.Run("graph a: nil; graph b: nil; equal", func(t *testing.T) {
		g := NewWithT(t)

		g.Expect(graph.Equal[string](nil, nil)).To(BeTrue())
	})

	t.Run("graph a: [link]; graph b: nil; not equal", func(t *testing.T) {
		g := NewWithT(t)

		a := undirectedGraphOf("link")
		var b graph.Graph[string]
		g.Expect(graph.Equal(a, b)).To(BeFalse())
	})

	t.Run("graph a: undirected; graph b: directed; not equal", func(t *testing.T) {
		g := NewWithT(t)
		a := undirectedGraphOf("link")
		b := directedGraphOf("link")

		g.Expect(graph.Equal(a, b)).To(BeFalse())
	})

	t.Run("graph a: allows self-loops; graph b: disallows self-loops; not equal", func(t *testing.T) {
		g := NewWithT(t)
		a := undirectedAllowsSelfLoopsGraphOf("link")
		b := undirectedGraphOf("link")

		g.Expect(graph.Equal(a, b)).To(BeFalse())
	})

	t.Run("graph a: [link]; graph b: [zelda]; not equal", func(t *testing.T) {
		g := NewWithT(t)
		a := undirectedGraphOf("link")
		b := undirectedGraphOf("zelda")

		g.Expect(graph.Equal(a, b)).To(BeFalse())
	})

	t.Run("graph a: [[link, zelda]]; graph b: [[link], [zelda]]; not equal", func(t *testing.T) {
		g := NewWithT(t)
		a := undirectedGraphOf(
			graph.NewUnorderedEndpointPair("link", "zelda"))
		b := undirectedGraphOf("link", "zelda")

		g.Expect(graph.Equal(a, b)).To(BeFalse())
	})
}

func undirectedGraphOf(nodesAndEndpointPairs ...any) graph.Graph[string] {
	result := graph.Undirected[string]().Build()
	for _, elem := range nodesAndEndpointPairs {
		switch e := elem.(type) {
		case string:
			result.AddNode(e)
		case graph.EndpointPair[string]:
			result.PutEdge(e.NodeU(), e.NodeV())
		default:
			panic(fmt.Sprintf("Unexpected elem: %v", elem))
		}
	}
	return result
}

func undirectedAllowsSelfLoopsGraphOf(nodesAndEndpointPairs ...any) graph.Graph[string] {
	result := graph.Undirected[string]().AllowsSelfLoops(true).Build()
	for _, elem := range nodesAndEndpointPairs {
		switch e := elem.(type) {
		case string:
			result.AddNode(e)
		case graph.EndpointPair[string]:
			result.PutEdge(e.NodeU(), e.NodeV())
		default:
			panic(fmt.Sprintf("Unexpected elem: %v", elem))
		}
	}
	return result
}

func directedGraphOf(nodes ...string) graph.Graph[string] {
	return minimalDirectedGraph{
		nodes: nodes,
	}
}

type minimalDirectedGraph struct {
	nodes []string
}

func (d minimalDirectedGraph) Nodes() set.Set[string] {
	return set.Of(d.nodes...)
}

func (d minimalDirectedGraph) Edges() set.Set[graph.EndpointPair[string]] {
	return set.Of[graph.EndpointPair[string]]()
}

func (d minimalDirectedGraph) IsDirected() bool {
	return true
}

func (d minimalDirectedGraph) AllowsSelfLoops() bool {
	return true
}

func (d minimalDirectedGraph) AdjacentNodes(_ string) set.Set[string] {
	panic("this is a minimal graph, so this method is purposefully not implemented")
}

func (d minimalDirectedGraph) Predecessors(_ string) set.Set[string] {
	panic("this is a minimal graph, so this method is purposefully not implemented")
}

func (d minimalDirectedGraph) Successors(_ string) set.Set[string] {
	panic("this is a minimal graph, so this method is purposefully not implemented")
}

func (d minimalDirectedGraph) IncidentEdges(_ string) set.Set[graph.EndpointPair[string]] {
	panic("this is a minimal graph, so this method is purposefully not implemented")
}

func (d minimalDirectedGraph) Degree(_ string) int {
	panic("this is a minimal graph, so this method is purposefully not implemented")
}

func (d minimalDirectedGraph) InDegree(_ string) int {
	panic("this is a minimal graph, so this method is purposefully not implemented")
}

func (d minimalDirectedGraph) OutDegree(_ string) int {
	panic("this is a minimal graph, so this method is purposefully not implemented")
}

func (d minimalDirectedGraph) HasEdgeConnecting(_ string, _ string) bool {
	panic("this is a minimal graph, so this method is purposefully not implemented")
}

func (d minimalDirectedGraph) HasEdgeConnectingEndpoints(_ graph.EndpointPair[string]) bool {
	panic("this is a minimal graph, so this method is purposefully not implemented")
}

func (d minimalDirectedGraph) String() string {
	panic("this is a minimal graph, so this method is purposefully not implemented")
}
