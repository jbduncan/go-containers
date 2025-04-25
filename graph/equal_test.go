package graph_test

import (
	"fmt"
	"testing"

	"github.com/jbduncan/go-containers/graph"
)

func TestEqual(t *testing.T) {
	t.Parallel()

	t.Run("graph a: [link]; graph b: [link]; equal", func(t *testing.T) {
		t.Parallel()
		a := undirectedGraphOf("link")

		if got := graph.Equal[string](a, a); !got {
			t.Errorf("graph.Equal: got false, want true")
		}
	})

	t.Run("graph a: nil; graph b: nil; equal", func(t *testing.T) {
		t.Parallel()
		if got := graph.Equal[string](nil, nil); !got {
			t.Errorf("graph.Equal: got false, want true")
		}
	})

	t.Run("graph a: [link]; graph b: nil; not equal", func(t *testing.T) {
		t.Parallel()
		a := undirectedGraphOf("link")
		if got := graph.Equal[string](a, nil); got {
			t.Errorf("graph.Equal: got true, want false")
		}
	})

	t.Run(
		"graph a: undirected; graph b: directed; not equal",
		func(t *testing.T) {
			t.Parallel()
			a := undirectedGraphOf("link")
			b := directedGraphOf("link")

			if got := graph.Equal[string](a, b); got {
				t.Errorf("graph.Equal: got true, want false")
			}
		},
	)

	t.Run(
		"graph a: allows self-loops; graph b: disallows self-loops; not equal",
		func(t *testing.T) {
			t.Parallel()
			a := graph.Undirected[string]().AllowsSelfLoops(true).Build()
			b := graph.Undirected[string]().Build()

			if got := graph.Equal[string](a, b); got {
				t.Errorf("graph.Equal: got true, want false")
			}
		},
	)

	t.Run("graph a: [link]; graph b: [zelda]; not equal", func(t *testing.T) {
		t.Parallel()
		a := undirectedGraphOf("link")
		b := undirectedGraphOf("zelda")

		if got := graph.Equal[string](a, b); got {
			t.Errorf("graph.Equal: got true, want false")
		}
	})

	t.Run(
		"graph a: [<link -> zelda>]; graph b: [link, zelda]; not equal",
		func(t *testing.T) {
			t.Parallel()
			a := undirectedGraphOf(edge("link", "zelda"))
			b := undirectedGraphOf("link", "zelda")

			if got := graph.Equal[string](a, b); got {
				t.Errorf("graph.Equal: got true, want false")
			}
		},
	)
}

func undirectedGraphOf(nodesAndEndpointPairs ...any) *graph.Graph[string] {
	result := graph.Undirected[string]().Build()
	for _, element := range nodesAndEndpointPairs {
		switch el := element.(type) {
		case string:
			result.AddNode(el)
		case graph.EndpointPair[string]:
			result.PutEdge(el.Source(), el.Target())
		default:
			panic(fmt.Sprintf("Unexpected element: %v", element))
		}
	}
	return result
}

func directedGraphOf(nodes ...string) *graph.Graph[string] {
	result := graph.Directed[string]().Build()
	for _, element := range nodes {
		result.AddNode(element)
	}
	return result
}

func edge(source string, target string) graph.EndpointPair[string] {
	return graph.EndpointPairOf(source, target)
}
