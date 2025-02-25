package graphtest

import (
	"fmt"
	"slices"
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

//go:generate stringer -type=Mutability
type Mutability int

const (
	Mutable Mutability = iota
	Immutable
)

//go:generate stringer -type=DirectionMode
type DirectionMode int

const (
	Directed DirectionMode = iota
	Undirected
)

//go:generate stringer -type=SelfLoopsMode
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
	mutability Mutability,
	directionMode DirectionMode,
	selfLoopsMode SelfLoopsMode,
) {
	if mutability != Mutable && mutability != Immutable {
		t.Fatalf(
			"mutability expected to be Mutable or Immutable "+
				"but was %v",
			mutability,
		)
	}
	if directionMode != Directed && directionMode != Undirected {
		t.Fatalf(
			"directionMode expected to be Directed or Undirected "+
				"but was %v",
			directionMode,
		)
	}
	if selfLoopsMode != AllowsSelfLoops &&
		selfLoopsMode != DisallowsSelfLoops {
		t.Fatalf(
			"selfLoopsMode expected to be AllowsSelfLoops or "+
				"DisallowsSelfLoops but was %v",
			selfLoopsMode,
		)
	}

	newTester(
		t,
		graphBuilder,
		mutability,
		directionMode,
		selfLoopsMode,
	).test()
}

func newTester(
	t TestingT,
	graphBuilder func() graph.Graph[int],
	mutability Mutability,
	directionMode DirectionMode,
	selfLoopsMode SelfLoopsMode,
) *tester {
	return &tester{
		t:             t,
		graphBuilder:  graphBuilder,
		mutability:    mutability,
		directionMode: directionMode,
		selfLoopsMode: selfLoopsMode,
	}
}

type tester struct {
	t             TestingT
	graphBuilder  func() graph.Graph[int]
	mutability    Mutability
	directionMode DirectionMode
	selfLoopsMode SelfLoopsMode
}

const (
	graphNodesName         = "Graph.Nodes"
	graphEdgesName         = "Graph.Edges"
	graphAdjacentNodesName = "Graph.AdjacentNodes"
	graphPredecessorsName  = "Graph.Predecessors"
	graphSuccessorsName    = "Graph.Successors"
	graphIncidentEdgesName = "Graph.IncidentEdges"
	graphDegreeName        = "Graph.Degree"
	graphInDegreeName      = "Graph.InDegree"
	graphOutDegreeName     = "Graph.OutDegree"
)

func (tt tester) test() {
	tt.t.Helper()

	tt.t.Run("empty graph", func(t *testing.T) {
		t.Run("has no nodes", func(t *testing.T) {
			g := tt.graphBuilder()
			testNodeSet(t, graphNodesName, g.Nodes())
		})

		t.Run("has no edges", func(t *testing.T) {
			g := tt.graphBuilder()
			tt.testEdges(t, g)
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
			tt.testIncidentEdges(t, g(), node1)
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
				tt.testIncidentEdges(
					t,
					g(),
					node1,
					graph.EndpointPairOf(node1, node2),
				)
			},
		)

		t.Run(
			"has just one edge",
			func(t *testing.T) {
				tt.testEdges(t, g(), graph.EndpointPairOf(node1, node2))
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

		t.Run(
			"sees the first node as being connected to no other node",
			func(t *testing.T) {
				if got := g().HasEdgeConnecting(node1, nodeNotInGraph); got {
					t.Errorf("Graph.HasEdgeConnecting: got true, want false")
				}
				if got := g().HasEdgeConnecting(nodeNotInGraph, node1); got {
					t.Errorf("Graph.HasEdgeConnecting: got true, want false")
				}
				if got := g().HasEdgeConnectingEndpoints(
					graph.EndpointPairOf(node1, nodeNotInGraph),
				); got {
					t.Errorf("Graph.HasEdgeConnectingEndpoints: " +
						"got true, want false",
					)
				}
				if got := g().HasEdgeConnectingEndpoints(
					graph.EndpointPairOf(nodeNotInGraph, node1),
				); got {
					t.Errorf("Graph.HasEdgeConnectingEndpoints: " +
						"got true, want false",
					)
				}
			},
		)

		t.Run(
			"sees the second node as being connected to no other node",
			func(t *testing.T) {
				if got := g().HasEdgeConnecting(node2, nodeNotInGraph); got {
					t.Errorf("Graph.HasEdgeConnecting: got true, want false")
				}
				if got := g().HasEdgeConnectingEndpoints(
					graph.EndpointPairOf(node2, nodeNotInGraph),
				); got {
					t.Errorf("Graph.HasEdgeConnectingEndpoints: " +
						"got true, want false",
					)
				}
			},
		)
	})

	tt.t.Run("graph with same edge put twice", func(t *testing.T) {
		t.Run("has only one edge", func(t *testing.T) {
			g := tt.graphBuilder()
			g = putEdge(g, node1, node2)

			tt.testEdges(t, g, graph.EndpointPairOf(node1, node2))
		})
	})

	tt.t.Run(
		"graph with two edges with the same source node",
		func(t *testing.T) {
			g := func() graph.Graph[int] {
				g := tt.graphBuilder()
				g = putEdge(g, node1, node2)
				g = putEdge(g, node1, node3)
				return g
			}

			t.Run("has a common node with a degree of 2", func(t *testing.T) {
				testDegree(t, graphDegreeName, g().Degree(node1), 2)
			})

			t.Run("has a common node with two successors", func(t *testing.T) {
				testNodeSet(
					t,
					"Graph.Successors",
					g().Successors(node1),
					node2,
					node3,
				)
			})

			// TODO: Adapt from graph_test.go, line 280, test "reports the two
			//       unique nodes as adjacent to the common one"
		},
	)
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

func (tt tester) testEdges(
	t *testing.T,
	g graph.Graph[int],
	expectedEdges ...graph.EndpointPair[int],
) {
	tt.testEdgeSet(
		t,
		graphEdgesName,
		g.Edges(),
		expectedEdges,
	)
}

func (tt tester) testIncidentEdges(
	t *testing.T,
	g graph.Graph[int],
	node int,
	expectedEdges ...graph.EndpointPair[int],
) {
	tt.testEdgeSet(
		t,
		graphIncidentEdgesName,
		g.IncidentEdges(node),
		expectedEdges,
	)
}

func (tt tester) testEdgeSet(
	t *testing.T,
	setName string,
	edges set.Set[graph.EndpointPair[int]],
	expectedEdges []graph.EndpointPair[int],
) {
	t.Helper()

	var contains []graph.EndpointPair[int]
	var doesNotContain []graph.EndpointPair[int]
	if tt.directionMode == Directed {
		contains = expectedEdges
		doesNotContain = slicesx.AllOf(
			graph.EndpointPairOf(nodeNotInGraph, nodeNotInGraph),
			reversesOf(expectedEdges),
		)
	} else {
		contains = slices.Concat(expectedEdges, reversesOf(expectedEdges))
		doesNotContain = []graph.EndpointPair[int]{
			graph.EndpointPairOf(nodeNotInGraph, nodeNotInGraph),
		}
	}

	testSetLen(t, setName, edges, len(expectedEdges))
	tt.testEdgeSetAll(t, setName, edges, expectedEdges)
	testSetContains(t, setName, edges, contains, doesNotContain)
	newEdgeSetStringTester(
		t,
		setName,
		tt.directionMode,
		edges,
		expectedEdges,
	).Test()
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

func (tt tester) testEdgeSetAll(
	t *testing.T,
	setName string,
	edges set.Set[graph.EndpointPair[int]],
	expectedEdges []graph.EndpointPair[int],
) {
	t.Helper()

	t.Run("Set.All", func(t *testing.T) {
		got, want := slices.Collect(edges.All()), expectedEdges
		if tt.directionMode == Directed {
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
		cmp.Comparer(
			func(
				a []graph.EndpointPair[int],
				b []graph.EndpointPair[int],
			) bool {
				aCopy, bCopy := deepCopyAndNormalise(a), deepCopyAndNormalise(b)
				return orderagnostic.SlicesEqual(aCopy, bCopy)
			},
		),
	)
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
