package graphtest

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"testing"

	gocmp "github.com/google/go-cmp/cmp"
	"github.com/jbduncan/go-containers/graph"
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

// Graph produces a suite of test cases for testing implementations of the graph.Graph and
// graph.MutableGraph interfaces. Graph instances created for testing are to have int nodes.
//
// Test cases that should be handled similarly in any graph implementation are included in this
// function; for example, testing that Nodes method returns the set of the nodes in the
// graph. Details of specific implementations of the graph.Graph and graph.MutableGraph
// interfaces are not tested.
func Graph(
	t TestingT,
	graphBuilder func() graph.Graph[int],
	mutableOrImmutable Mutability,
	directedOrUndirected DirectionMode,
	allowsOrDisallowsSelfLoops SelfLoopsMode,
) {
	if mutableOrImmutable != Mutable && mutableOrImmutable != Immutable {
		t.Fatalf(
			"mutableOrImmutable expected to be Mutable or Immutable but was %v",
			mutableOrImmutable,
		)
	}
	if directedOrUndirected != Directed && directedOrUndirected != Undirected {
		t.Fatalf(
			"directedOrUndirected expected to be Directed or Undirected but was %v",
			directedOrUndirected,
		)
	}
	if allowsOrDisallowsSelfLoops != AllowsSelfLoops &&
		allowsOrDisallowsSelfLoops != DisallowsSelfLoops {
		t.Fatalf(
			"allowsOrDisallowsSelfLoops expected to be AllowsSelfLoops or DisallowsSelfLoops but was %v",
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
			testEmptyEdges(t, graphEdgesName, g.Edges())
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
			testEmptyEdges(t, graphIncidentEdgesName, g().IncidentEdges(node1))
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
			"has an incident edge connecting the first node to the second node",
			func(t *testing.T) {
				g := g()
				if tt.directedOrUndirected == Directed {
					testSingleDirectedEdge(
						t,
						graphIncidentEdgesName,
						g.IncidentEdges(node1),
						graph.EndpointPairOf(node1, node2),
					)
				} else {
					testSingleUndirectedEdge(
						t,
						graphIncidentEdgesName,
						g.IncidentEdges(node1),
						graph.EndpointPairOf(node1, node2),
					)
				}
			},
		)

		// TODO: refactor the common bits between test(NodeSet|EmptyEdges|SingleDirectedEdge|SingleUndirectedEdge)
		// TODO: continue from graph_test.go, line 218, "has just one edge"
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

func testNodeSet(
	t *testing.T,
	setName string,
	s set.Set[int],
	expectedValues ...int,
) {
	t.Helper()

	testSetLen(t, setName, s, len(expectedValues))
	testSetAll(t, setName, s, expectedValues)

	t.Run("Set.Contains", func(t *testing.T) {
		for _, value := range []int{node1, node2, node3} {
			if slices.Contains(expectedValues, value) {
				if !s.Contains(value) {
					t.Errorf(
						"%s: got Set.Contains(%d) == false, want true",
						setName,
						value,
					)
				}
			} else {
				if s.Contains(value) {
					t.Errorf(
						"%s: got Set.Contains(%d) == true, want false",
						setName,
						value,
					)
				}
			}
		}
		if s.Contains(nodeNotInGraph) {
			t.Errorf(
				"%s: got Set.Contains(%d) == true, want false",
				setName,
				nodeNotInGraph,
			)
		}
	})

	testSetString(t, setName, s, expectedValues)
}

func testEmptyEdges(
	t *testing.T,
	setName string,
	edges set.Set[graph.EndpointPair[int]],
) {
	t.Helper()

	testSetLen(t, setName, edges, 0)
	testSetAll(t, setName, edges, make([]graph.EndpointPair[int], 0))

	t.Run("Set.Contains", func(t *testing.T) {
		if edges.Contains(
			graph.EndpointPairOf(nodeNotInGraph, nodeNotInGraph),
		) {
			t.Errorf(
				"%s: got Set.Contains(%s) == true, want false",
				setName,
				graph.EndpointPairOf(nodeNotInGraph, nodeNotInGraph),
			)
		}
	})

	testSetString(t, setName, edges, make([]graph.EndpointPair[int], 0))
}

func testSingleDirectedEdge(
	t *testing.T,
	setName string,
	edges set.Set[graph.EndpointPair[int]],
	expectedEdge graph.EndpointPair[int],
) {
	t.Helper()

	testSetLen(t, setName, edges, 1)
	testSetAll(t, setName, edges, []graph.EndpointPair[int]{expectedEdge})

	t.Run("Set.Contains", func(t *testing.T) {
		if !edges.Contains(expectedEdge) {
			t.Errorf(
				"%s: got Set.Contains(%s) == false, want true",
				setName,
				expectedEdge,
			)
		}

		if expectedEdge != reverseOf(expectedEdge) {
			if edges.Contains(reverseOf(expectedEdge)) {
				t.Errorf(
					"%s: got Set.Contains(%s) == true, want false",
					setName,
					reverseOf(expectedEdge),
				)
			}
		}

		if edges.Contains(
			graph.EndpointPairOf(nodeNotInGraph, nodeNotInGraph),
		) {
			t.Errorf(
				"%s: got Set.Contains(%s) == true, want false",
				setName,
				graph.EndpointPairOf(nodeNotInGraph, nodeNotInGraph),
			)
		}
	})

	testSetString(t, setName, edges, []graph.EndpointPair[int]{expectedEdge})
}

func testSingleUndirectedEdge(
	t *testing.T,
	setName string,
	edges set.Set[graph.EndpointPair[int]],
	expectedEdge graph.EndpointPair[int],
) {
	t.Helper()

	testSetLen(t, setName, edges, 1)
	testSetAll(
		t,
		setName,
		edges,
		[]graph.EndpointPair[int]{expectedEdge},
		orderAgnosticEndpointPairComparer(),
	)

	t.Run("Set.Contains", func(t *testing.T) {
		if !edges.Contains(expectedEdge) {
			t.Errorf(
				"%s: got Set.Contains(%s) == false, want true",
				setName,
				expectedEdge,
			)
		}

		if expectedEdge != reverseOf(expectedEdge) {
			if !edges.Contains(reverseOf(expectedEdge)) {
				t.Errorf(
					"%s: got Set.Contains(%s) == false, want true",
					setName,
					reverseOf(expectedEdge),
				)
			}
		}

		if edges.Contains(
			graph.EndpointPairOf(nodeNotInGraph, nodeNotInGraph),
		) {
			t.Errorf(
				"%s: got Set.Contains(%s) == true, want false",
				setName,
				graph.EndpointPairOf(nodeNotInGraph, nodeNotInGraph),
			)
		}
	})

	testSetString(
		t,
		setName,
		edges,
		[]graph.EndpointPair[int]{expectedEdge},
		orderAgnosticEndpointPairComparer(),
	)
}

func testSetLen[T comparable](
	t *testing.T,
	setName string,
	s set.Set[T],
	expectedLen int,
) {
	t.Helper()

	t.Run("Set.Len", func(t *testing.T) {
		setLen := s.Len()
		if setLen != expectedLen {
			t.Errorf(
				"%s: got Set.Len of %d, want %d",
				setName,
				setLen,
				expectedLen,
			)
		}
	})
}

func testSetAll[T comparable](
	t *testing.T,
	setName string,
	s set.Set[T],
	expectedValues []T,
	extraOptions ...gocmp.Option,
) {
	t.Helper()

	t.Run("Set.All", func(t *testing.T) {
		all := slices.Collect(s.All())
		diff := orderAgnosticDiff(all, expectedValues, extraOptions...)
		if diff != "" {
			t.Errorf("%s: Set.All mismatch (-want +got):\n%s", setName, diff)
		}
	})
}

func testSetString[T comparable](
	t *testing.T,
	setName string,
	s set.Set[T],
	expectedValues []T,
	extraOptions ...gocmp.Option,
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

		var expectedValueStrs []string
		for _, v := range expectedValues {
			expectedValueStrs = append(expectedValueStrs, fmt.Sprintf("%v", v))
		}
		actualValueStrs := strings.SplitN(trimmed, ", ", len(expectedValues))

		diff := orderAgnosticDiff(actualValueStrs, expectedValueStrs, extraOptions...)
		if diff != "" {
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
	if actualDegree != expectedDegree {
		t.Errorf(
			"%s: got degree of %d, want %d",
			degreeName,
			actualDegree,
			expectedDegree,
		)
	}
}

func reverseOf(endpointPair graph.EndpointPair[int]) graph.EndpointPair[int] {
	return graph.EndpointPairOf(endpointPair.Target(), endpointPair.Source())
}

func orderAgnosticDiff[T comparable](
	got []T,
	want []T,
	extraOptions ...gocmp.Option,
) string {
	sliceComparer := gocmp.Comparer(func(a, b []T) bool {
		x := make(map[T]int)
		for _, value := range a {
			x[value]++
		}
		y := make(map[T]int)
		for _, value := range b {
			y[value]++
		}
		return maps.Equal(x, y)
	})
	allOptions := allOf(
		sliceComparer,
		extraOptions,
	)
	return gocmp.Diff(
		want,
		got,
		allOptions...,
	)
}

func orderAgnosticEndpointPairComparer() gocmp.Option {
	// TODO(jbduncan): consider reintroducing EndpointPair.Equal so this
	//                 comparer can be removed.
	return gocmp.Comparer(func(a, b graph.EndpointPair[int]) bool {
		return a == b || a == reverseOf(b)
	})
}

func allOf[T any](first T, rest []T) []T {
	return slices.Concat([]T{first}, rest)
}
