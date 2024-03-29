package graph_test

import (
	"fmt"
	"testing"

	"github.com/jbduncan/go-containers/graph"
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
		a := undirectedGraphOf(edge("link", "zelda"))
		b := undirectedGraphOf("link", "zelda")

		g.Expect(graph.Equal(a, b)).To(BeFalse())
	})
}

func undirectedGraphOf(nodesAndEndpointPairs ...any) graph.Graph[string] {
	result := graph.Undirected[string]().Build()
	for _, elem := range nodesAndEndpointPairs {
		switch el := elem.(type) {
		case string:
			result.AddNode(el)
		case graph.EndpointPair[string]:
			result.PutEdge(el.Source(), el.Target())
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
			result.PutEdge(e.Source(), e.Target())
		default:
			panic(fmt.Sprintf("Unexpected elem: %v", elem))
		}
	}
	return result
}

func directedGraphOf(nodes ...string) graph.Graph[string] {
	result := graph.Directed[string]().Build()
	for _, elem := range nodes {
		result.AddNode(elem)
	}
	return result
}

func edge(source string, target string) graph.EndpointPair[string] {
	return graph.EndpointPairOf(source, target)
}
