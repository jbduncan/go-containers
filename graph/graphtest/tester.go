package graphtest

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jbduncan/go-containers/graph"
	"github.com/jbduncan/go-containers/internal/orderagnostic"
	"github.com/jbduncan/go-containers/internal/slicesx"
	"github.com/jbduncan/go-containers/set"
)

// TestingT is an interface for the parts of *testing.T that graphtest.Graph
// needs to run. Whenever you see an argument of this type, pass in an instance
// of *testing.T or your unit testing framework's equivalent.
type TestingT interface {
	Helper()
	Fatalf(format string, args ...any)
	Run(name string, f func(t *testing.T)) bool
}

const (
	node1          = 1
	node2          = 2
	node3          = 3
	nodeNotInGraph = 1_000
)

type Mutability int

const (
	Mutable Mutability = iota
	Immutable
)

type DirectionMode int

const (
	Directed DirectionMode = iota
	Undirected
)

type SelfLoopsMode int

const (
	AllowsSelfLoops SelfLoopsMode = iota
	DisallowsSelfLoops
)

// Graph produces a suite of test cases for testing implementations of the
// graph.Graph and graph.MutableGraph interfaces. Graph instances created for
// testing are to have int nodes.
//
// Test cases that should be handled similarly in any graph implementation are
// included in this function; for example, testing that Nodes method returns
// the set of the nodes in the graph. Details of specific implementations of
// the graph.Graph and graph.MutableGraph interfaces are not tested.
func Graph(
	t TestingT,
	graphBuilder func() graph.Graph[int],
	mutableOrImmutable Mutability,
	directedOrUndirected DirectionMode,
	allowsOrDisallowsSelfLoops SelfLoopsMode,
) {
	if mutableOrImmutable != Mutable && mutableOrImmutable != Immutable {
		t.Fatalf(
			"mutableOrImmutable expected to be Mutable or Immutable "+
				"but was %v",
			mutableOrImmutable,
		)
	}
	if directedOrUndirected != Directed && directedOrUndirected != Undirected {
		t.Fatalf(
			"directedOrUndirected expected to be Directed or Undirected "+
				"but was %v",
			directedOrUndirected,
		)
	}
	if allowsOrDisallowsSelfLoops != AllowsSelfLoops &&
		allowsOrDisallowsSelfLoops != DisallowsSelfLoops {
		t.Fatalf(
			"allowsOrDisallowsSelfLoops expected to be AllowsSelfLoops or "+
				"DisallowsSelfLoops but was %v",
			allowsOrDisallowsSelfLoops,
		)
	}

	newTester(
		t,
		graphBuilder,
		mutableOrImmutable,
		directedOrUndirected,
		allowsOrDisallowsSelfLoops,
	).test()
}

func newTester(
	t TestingT,
	graphBuilder func() graph.Graph[int],
	mutableOrImmutable Mutability,
	directedOrUndirected DirectionMode,
	allowsOrDisallowsSelfLoops SelfLoopsMode,
) *tester {
	return &tester{
		t:                          t,
		graphBuilder:               graphBuilder,
		mutableOrImmutable:         mutableOrImmutable,
		directedOrUndirected:       directedOrUndirected,
		allowsOrDisallowsSelfLoops: allowsOrDisallowsSelfLoops,
	}
}

type tester struct {
	t                          TestingT
	graphBuilder               func() graph.Graph[int]
	mutableOrImmutable         Mutability
	directedOrUndirected       DirectionMode
	allowsOrDisallowsSelfLoops SelfLoopsMode
}

func (tt tester) test() {
	tt.t.Helper()

	const graphNodesName = "Graph.Nodes"
	const graphEdgesName = "Graph.Edges"
	const graphAdjacentNodesName = "Graph.AdjacentNodes"
	const graphPredecessorsName = "Graph.Predecessors"
	const graphSuccessorsName = "Graph.Successors"
	const graphIncidentEdgesName = "Graph.IncidentEdges"
	const graphDegreeName = "Graph.Degree"
	const graphInDegreeName = "Graph.InDegree"
	const graphOutDegreeName = "Graph.OutDegree"

	tt.t.Run("empty graph", func(t *testing.T) {
		t.Run("has no nodes", func(t *testing.T) {
			g := tt.graphBuilder()
			testNodeSet(t, graphNodesName, g.Nodes())
		})

		t.Run("has no edges", func(t *testing.T) {
			g := tt.graphBuilder()
			tt.testEdgeSet(
				t,
				graphEdgesName,
				g.Edges(),
				g.IsDirected(),
				nil,
			)
		})
	})

	tt.t.Run("graph with one node", func(t *testing.T) {
		g := func() graph.Graph[int] {
			g := tt.graphBuilder()
			g = addNode(g, node1)
			return g
		}
		t.Run("has just that node", func(t *testing.T) {
			testNodeSet(t, graphNodesName, g().Nodes(), node1)
		})

		t.Run("the node has no adjacent nodes", func(t *testing.T) {
			testNodeSet(t, graphAdjacentNodesName, g().AdjacentNodes(node1))
		})

		t.Run("the node has no predecessors", func(t *testing.T) {
			testNodeSet(t, graphPredecessorsName, g().Predecessors(node1))
		})

		t.Run("the node has no successors", func(t *testing.T) {
			testNodeSet(t, graphSuccessorsName, g().Successors(node1))
		})

		t.Run("the node has no incident edges", func(t *testing.T) {
			g := g()
			tt.testEdgeSet(
				t,
				graphIncidentEdgesName,
				g.IncidentEdges(node1),
				g.IsDirected(),
				nil,
			)
		})

		t.Run("the node has a degree of 0", func(t *testing.T) {
			testDegree(t, graphDegreeName, g().Degree(node1), 0)
		})

		t.Run("the node has an in-degree of 0", func(t *testing.T) {
			testDegree(t, graphInDegreeName, g().InDegree(node1), 0)
		})

		t.Run("the node has an out-degree of 0", func(t *testing.T) {
			testDegree(t, graphOutDegreeName, g().OutDegree(node1), 0)
		})
	})

	tt.t.Run("graph with two nodes", func(t *testing.T) {
		t.Run("has both nodes", func(t *testing.T) {
			g := tt.graphBuilder()
			g = addNode(g, node1)
			g = addNode(g, node2)

			testNodeSet(t, graphNodesName, g.Nodes(), node1, node2)
		})
	})

	tt.t.Run("graph with one edge", func(t *testing.T) {
		g := func() graph.Graph[int] {
			g := tt.graphBuilder()
			g = putEdge(g, node1, node2)
			return g
		}

		t.Run(
			"the source node is adjacent to the target node",
			func(t *testing.T) {
				testNodeSet(
					t,
					graphAdjacentNodesName,
					g().AdjacentNodes(node1),
					node2,
				)
			},
		)

		t.Run(
			"the target node is adjacent to the source node",
			func(t *testing.T) {
				testNodeSet(
					t,
					graphAdjacentNodesName,
					g().AdjacentNodes(node2),
					node1,
				)
			},
		)

		t.Run(
			"the source node is the predecessor of the target node",
			func(t *testing.T) {
				testNodeSet(
					t,
					graphPredecessorsName,
					g().Predecessors(node2),
					node1,
				)
			},
		)

		t.Run(
			"the target node is the successor of the source node",
			func(t *testing.T) {
				testNodeSet(
					t,
					graphSuccessorsName,
					g().Successors(node1),
					node2,
				)
			},
		)

		t.Run("the source node has a degree of 1", func(t *testing.T) {
			testDegree(t, graphDegreeName, g().Degree(node1), 1)
		})

		t.Run("the target node has a degree of 1", func(t *testing.T) {
			testDegree(t, graphDegreeName, g().Degree(node2), 1)
		})

		t.Run("the target node has an in-degree of 1", func(t *testing.T) {
			testDegree(t, graphInDegreeName, g().InDegree(node2), 1)
		})

		t.Run("the source node has an out-degree of 1", func(t *testing.T) {
			testDegree(t, graphOutDegreeName, g().OutDegree(node1), 1)
		})

		t.Run(
			"has an incident edge connecting the first node to the "+
				"second node",
			func(t *testing.T) {
				g := g()
				tt.testEdgeSet(
					t,
					graphIncidentEdgesName,
					g.IncidentEdges(node1),
					g.IsDirected(),
					[]graph.EndpointPair[int]{
						graph.EndpointPairOf(node1, node2),
					},
				)
			},
		)

		t.Run(
			"has just one edge",
			func(t *testing.T) {
				g := g()
				tt.testEdgeSet(
					t,
					graphEdgesName,
					g.Edges(),
					g.IsDirected(),
					[]graph.EndpointPair[int]{
						graph.EndpointPairOf(node1, node2),
					},
				)
			},
		)

		t.Run(
			"sees the first node as being connected to the second",
			func(t *testing.T) {
				if got := g().HasEdgeConnecting(node1, node2); !got {
					t.Errorf("Graph.HasEdgeConnecting: got false, want true")
				}
				if got := g().HasEdgeConnectingEndpoints(
					graph.EndpointPairOf(node1, node2),
				); !got {
					t.Errorf(
						"Graph.HasEdgeConnectingEndpoints: " +
							"got false, want true",
					)
				}
			},
		)

		// TODO: continue from graph_test.go, line 233, "does not see the first node as being connected to any other node"
	})
}

func addNode(g graph.Graph[int], node int) graph.Graph[int] {
	if gAsMutable, ok := g.(graph.MutableGraph[int]); ok {
		gAsMutable.AddNode(node)
	}

	return g
}

func putEdge(g graph.Graph[int], source int, target int) graph.Graph[int] {
	if gAsMutable, ok := g.(graph.MutableGraph[int]); ok {
		gAsMutable.PutEdge(source, target)
	}

	return g
}

var allNodesToConsider = []int{node1, node2, node3, nodeNotInGraph}

func complement(nodes []int) []int {
	result := slices.Clone(allNodesToConsider)
	result = slices.DeleteFunc(result, func(value int) bool {
		return slices.Contains(nodes, value)
	})
	return result
}

func testNodeSet(
	t *testing.T,
	setName string,
	s set.Set[int],
	expectedValues ...int,
) {
	t.Helper()

	testSetLen(t, setName, s, len(expectedValues))
	testSetAll(t, setName, s, expectedValues)
	testSetContains(
		t,
		setName,
		s,
		expectedValues,
		complement(expectedValues),
	)
	testSetString(t, setName, s, expectedValues)
}

// TODO: refactor 'graphIsDirected' control flag away
//
//nolint:revive
func (tt tester) testEdgeSet(
	t *testing.T,
	setName string,
	edges set.Set[graph.EndpointPair[int]],
	graphIsDirected bool,
	expectedEdges []graph.EndpointPair[int],
) {
	t.Helper()

	var contains []graph.EndpointPair[int]
	var doesNotContain []graph.EndpointPair[int]
	if tt.directedOrUndirected == Directed {
		contains = expectedEdges
		doesNotContain = slicesx.AllOf(
			graph.EndpointPairOf(nodeNotInGraph, nodeNotInGraph),
			reversesOf(expectedEdges),
		)
	} else {
		// Even though there are only len(expectedEdges) edges in the graph,
		// test that the set contains both the edges and their reverses to
		// make sure that the edges are undirected.
		contains = slices.Concat(expectedEdges, reversesOf(expectedEdges))
		doesNotContain = []graph.EndpointPair[int]{
			graph.EndpointPairOf(nodeNotInGraph, nodeNotInGraph),
		}
	}

	testSetLen(t, setName, edges, len(expectedEdges))
	// TODO: Extract into helper function
	t.Run("Set.All", func(t *testing.T) {
		got, want := slices.Collect(edges.All()), expectedEdges
		if graphIsDirected {
			if diff := orderagnostic.Diff(got, want); diff != "" {
				t.Errorf(
					"%s: Set.All mismatch (-want +got):\n%s",
					setName,
					diff,
				)
			}
		} else {
			if diff := undirectedEndpointPairsDiff(got, want); diff != "" {
				t.Errorf(
					"%s: Set.All mismatch (-want +got):\n%s",
					setName,
					diff,
				)
			}
		}
	})
	testSetContains(t, setName, edges, contains, doesNotContain)
	// TODO: Extract into helper function
	t.Run("Set.String", func(t *testing.T) {
		str := edges.String()

		report := func() {
			t.Helper()

			var msg strings.Builder
			if len(expectedEdges) == 0 {
				msg.WriteString("%s: got Set.String of %q, want \"[]\"")
			} else if graphIsDirected {
				msg.WriteString("%s: got Set.String of %q, want to contain substrings:\n")
				for _, edge := range expectedEdges {
					msg.WriteString("    ")
					msg.WriteString(edge.String())
				}
			} else {
				msg.WriteString("%s: got Set.String of %q, want to contain substrings:\n")
				for _, edge := range expectedEdges {
					msg.WriteString("    ")
					msg.WriteString(edge.String())
					msg.WriteString(" or ")
					msg.WriteString(reverseOf(edge).String())
				}
			}
			t.Fatalf(msg.String(), setName, str)
		}

		parseEndpointPairString := func(s string) graph.EndpointPair[int] {
			t.Helper()

			// TODO: Extract into global variable
			endpointPairStringRegex := regexp.MustCompile(`<(\d+) -> (\d)+>`)
			matches := endpointPairStringRegex.FindStringSubmatch(s)
			if len(matches) != 3 {
				report()
			}
			source, err := strconv.Atoi(matches[1])
			if err != nil {
				report()
			}
			target, err := strconv.Atoi(matches[2])
			if err != nil {
				report()
			}
			return graph.EndpointPairOf(source, target)
		}

		parseString := func() []graph.EndpointPair[int] {
			t.Helper()

			trimmed, prefixFound := strings.CutPrefix(str, "[")
			if !prefixFound {
				t.Fatalf(
					`%s: got Set.String of %q, want to have prefix "["`,
					setName,
					str,
				)
			}
			trimmed, suffixFound := strings.CutSuffix(trimmed, "]")
			if !suffixFound {
				t.Fatalf(
					`%s: got Set.String of %q, want to have suffix "]"`,
					setName,
					str,
				)
			}

			elems := splitByComma(trimmed)
			result := make([]graph.EndpointPair[int], 0, len(elems))
			for _, elemStr := range elems {
				result = append(result, parseEndpointPairString(elemStr))
			}
			return result
		}

		got, want := expectedEdges, parseString()
		if graphIsDirected {
			if diff := orderagnostic.Diff(got, want); diff != "" {
				report()
			}
		} else {
			if diff := undirectedEndpointPairsDiff(got, want); diff != "" {
				report()
			}
		}
	})
}

func testSetLen[T comparable](
	t *testing.T,
	setName string,
	s set.Set[T],
	expectedLen int,
) {
	t.Helper()

	t.Run("Set.Len", func(t *testing.T) {
		if got, want := s.Len(), expectedLen; got != want {
			t.Errorf(
				"%s: got Set.Len of %d, want %d",
				setName,
				got,
				want,
			)
		}
	})
}

func testSetAll[T comparable](
	t *testing.T,
	setName string,
	s set.Set[T],
	expectedValues []T,
) {
	t.Helper()

	t.Run("Set.All", func(t *testing.T) {
		got, want := slices.Collect(s.All()), expectedValues
		if diff := orderagnostic.Diff(got, want); diff != "" {
			t.Errorf("%s: Set.All mismatch (-want +got):\n%s", setName, diff)
		}
	})
}

func testSetContains[T comparable](
	t *testing.T,
	setName string,
	s set.Set[T],
	contains []T,
	doesNotContain []T,
) {
	t.Helper()

	t.Run("Set.Contains", func(t *testing.T) {
		for _, value := range contains {
			if !s.Contains(value) {
				t.Errorf(
					"%s: got Set.Contains(%v) == false, want true",
					setName,
					value,
				)
			}
		}
		for _, value := range doesNotContain {
			if s.Contains(value) {
				t.Errorf(
					"%s: got Set.Contains(%v) == true, want false",
					setName,
					value,
				)
			}
		}
	})
}

func testSetString[T comparable](
	t *testing.T,
	setName string,
	s set.Set[T],
	expectedValues []T,
) {
	t.Helper()

	t.Run("Set.String", func(t *testing.T) {
		str := s.String()
		trimmed, prefixFound := strings.CutPrefix(str, "[")
		if !prefixFound {
			t.Fatalf(
				`%s: got Set.String of %q, want to have prefix "["`,
				setName,
				str,
			)
		}
		trimmed, suffixFound := strings.CutSuffix(trimmed, "]")
		if !suffixFound {
			t.Fatalf(
				`%s: got Set.String of %q, want to have suffix "]"`,
				setName,
				str,
			)
		}

		want := make([]string, 0, len(expectedValues))
		for _, v := range expectedValues {
			want = append(want, fmt.Sprintf("%v", v))
		}
		got := splitByComma(trimmed)

		if diff := orderagnostic.Diff(got, want); diff != "" {
			t.Fatalf(
				"%s: Set.String of %q: elements mismatch: (-want +got):\n%s",
				setName,
				str,
				diff,
			)
		}
	})
}

func testDegree(
	t *testing.T,
	degreeName string,
	actualDegree int,
	expectedDegree int,
) {
	if got, want := actualDegree, expectedDegree; got != want {
		t.Errorf(
			"%s: got degree of %d, want %d",
			degreeName,
			got,
			want,
		)
	}
}

func splitByComma(s string) []string {
	if len(s) == 0 {
		return make([]string, 0)
	}
	return strings.Split(s, ", ")
}

func reverseOf(endpointPair graph.EndpointPair[int]) graph.EndpointPair[int] {
	return graph.EndpointPairOf(endpointPair.Target(), endpointPair.Source())
}

func reversesOf(edges []graph.EndpointPair[int]) []graph.EndpointPair[int] {
	result := make([]graph.EndpointPair[int], 0, len(edges))
	for _, edge := range edges {
		result = append(result, reverseOf(edge))
	}
	return result
}

func undirectedEndpointPairsDiff(
	got []graph.EndpointPair[int],
	want []graph.EndpointPair[int],
) string {
	return cmp.Diff(
		want,
		got,
		cmp.Comparer(undirectedEndpointPairSlicesEqualInAnyOrder),
	)
}

func undirectedEndpointPairSlicesEqualInAnyOrder(
	a []graph.EndpointPair[int],
	b []graph.EndpointPair[int],
) bool {
	aCopy, bCopy := deepCopyAndNormalise(a), deepCopyAndNormalise(b)
	return orderagnostic.SlicesEqual(aCopy, bCopy)
}

func deepCopyAndNormalise(
	s []graph.EndpointPair[int],
) []graph.EndpointPair[int] {
	var result []graph.EndpointPair[int]
	for _, edge := range s {
		var newEdge graph.EndpointPair[int]
		if edge.Source() < edge.Target() {
			newEdge = graph.EndpointPairOf(edge.Source(), edge.Target())
		} else {
			newEdge = graph.EndpointPairOf(edge.Target(), edge.Source())
		}
		result = append(result, newEdge)
	}
	return result
}
