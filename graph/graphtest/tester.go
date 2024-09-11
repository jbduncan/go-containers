package graphtest

import (
	"cmp"
	gocmp "github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jbduncan/go-containers/set"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/jbduncan/go-containers/graph"
)

// TestingT is an interface for the parts of *testing.T that graphtest.Graph
// needs to run. Whenever you see an argument of this type, pass in an instance
// of *testing.T or your unit testing framework's equivalent.
type TestingT interface {
	Helper()
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
	tt := newTester(
		t,
		graphBuilder,
		mutableOrImmutable,
		directedOrUndirected,
		allowsOrDisallowsSelfLoops,
	)

	tt.emptyGraphHasNoNodes()
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

func (tt tester) emptyGraphHasNoNodes() {
	tt.t.Helper()

	const graphNodesName = "Graph.Nodes"
	const graphEdgesName = "Graph.Edges"
	const graphAdjacentNodesName = "Graph.AdjacentNodes"
	const graphPredecessorsName = "Graph.Predecessors"
	const graphSuccessorsName = "Graph.Successors"
	const graphIncidentEdgesName = "Graph.IncidentEdges"

	tt.t.Run("empty graph", func(t *testing.T) {
		t.Run("has no nodes", func(t *testing.T) {
			g := tt.graphBuilder()
			testSet(t, graphNodesName, g.Nodes())
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
			testSet(t, graphNodesName, g().Nodes(), node1)
		})

		t.Run("the node has no adjacent nodes", func(t *testing.T) {
			testSet(t, graphAdjacentNodesName, g().AdjacentNodes(node1))
		})

		t.Run("the node has no predecessors", func(t *testing.T) {
			testSet(t, graphPredecessorsName, g().Predecessors(node1))
		})

		t.Run("the node has no successors", func(t *testing.T) {
			testSet(t, graphSuccessorsName, g().Successors(node1))
		})

		t.Run("the node has no incident edges", func(t *testing.T) {
			testEmptyEdges(t, graphIncidentEdgesName, g().IncidentEdges(node1))
		})

		t.Run("the node has a degree of 0", func(t *testing.T) {
			g := g()
			degree := g.Degree(node1)
			if degree != 0 {
				t.Errorf("graph.Degree: got degree of %d, want 0", degree)
			}
		})

		t.Run("the node has an in-degree of 0", func(t *testing.T) {
			g := g()
			inDegree := g.InDegree(node1)
			if inDegree != 0 {
				t.Errorf(
					"graph.InDegree: got in-degree of %d, want 0",
					inDegree,
				)
			}
		})

		t.Run("the node has an out-degree of 0", func(t *testing.T) {
			g := g()
			outDegree := g.OutDegree(node1)
			if outDegree != 0 {
				t.Errorf(
					"graph.OutDegree: got out-degree of %d, want 0",
					outDegree,
				)
			}
		})
	})
}

func addNode(g graph.Graph[int], node int) graph.Graph[int] {
	if gAsMutable, ok := g.(graph.MutableGraph[int]); ok {
		gAsMutable.AddNode(node)
	}

	return g
}

func testSet(
	t *testing.T,
	setName string,
	s set.Set[int],
	expectedValues ...int,
) {
	t.Helper()

	t.Run("Set.Len", func(t *testing.T) {
		setLen := s.Len()
		if setLen != len(expectedValues) {
			t.Errorf("%s: got Set.Len of %d, want %d", setName, setLen, len(expectedValues))
		}
	})

	t.Run("Set.All", func(t *testing.T) {
		all := slices.Collect(s.All())
		diff := orderAgnosticDiff(all, expectedValues)
		if diff != "" {
			t.Errorf("%s: Set.All mismatch (-want +got):\n%s", setName, diff)
		}
	})

	// Set.Contains()
	t.Run("Set.Contains", func(t *testing.T) {
		for _, value := range []int{node1, node2, node3} {
			if slices.Contains(expectedValues, value) {
				if !s.Contains(value) {
					t.Errorf("%s: got Set.Contains(%d) == false, want true", setName, value)
				}
			} else {
				if s.Contains(value) {
					t.Errorf("%s: got Set.Contains(%d) == true, want false", setName, value)
				}
			}
		}
		if s.Contains(nodeNotInGraph) {
			t.Errorf("%s: got Set.Contains(%d) == true, want false", setName, nodeNotInGraph)
		}
	})

	t.Run("Set.String", func(t *testing.T) {
		str := s.String()
		trimmed, prefixFound := strings.CutPrefix(str, "[")
		if !prefixFound {
			t.Fatalf(
				`%s: got Set.String of %s, want to have prefix "["`,
				setName,
				str,
			)
		}
		trimmed, suffixFound := strings.CutSuffix(trimmed, "]")
		if !suffixFound {
			t.Fatalf(
				`%s: got Set.String of %s, want to have suffix "]"`,
				setName,
				str,
			)
		}

		var expectedValueStrs []string
		for _, v := range expectedValues {
			expectedValueStrs = append(expectedValueStrs, strconv.Itoa(v))
		}

		actualValueStrs := strings.SplitN(trimmed, ", ", len(expectedValues))
		diff := orderAgnosticDiff(actualValueStrs, expectedValueStrs)
		if diff != "" {
			t.Fatalf(
				"%s: got Set.String of %s, want elements to be %v in any order: (-want +got):\n%s",
				setName,
				str,
				expectedValueStrs,
				diff,
			)
		}
	})
}

func testEmptyEdges(
	t *testing.T,
	setName string,
	edges set.Set[graph.EndpointPair[int]],
) {
	t.Run("Set.Len", func(t *testing.T) {
		setLen := edges.Len()
		if setLen != 0 {
			t.Errorf("%s: got Set.Len of %d, want 0", setName, setLen)
		}
	})

	t.Run("Set.All", func(t *testing.T) {
		all := slices.Collect(edges.All())
		if len(all) != 0 {
			t.Errorf("%s: got Set.All len of %d, want 0", setName, len(all))
		}
	})

	t.Run("Set.Contains", func(t *testing.T) {
		if edges.Contains(
			graph.EndpointPairOf(nodeNotInGraph, nodeNotInGraph),
		) {
			t.Errorf(
				"%s: got Set.Contains(graph.EndpointPairOf(nodeNotInGraph, nodeNotInGraph)) == true, want false",
				setName,
			)
		}
	})

	t.Run("Set.String", func(t *testing.T) {
		setString := edges.String()
		if setString != "[]" {
			t.Fatalf(`%s: got Set.String of %q, want "[]"`, setName, setString)
		}
	})
}

func orderAgnosticDiff[T cmp.Ordered](actual []T, expected []T) string {
	return gocmp.Diff(
		actual,
		expected,
		cmpopts.SortSlices(func(a, b T) bool {
			return a < b
		}),
	)
}
