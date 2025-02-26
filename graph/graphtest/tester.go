package graphtest

import (
	"slices"
	"testing"

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

//go:generate mise x -- stringer -type=Mutability
type Mutability int

const (
	Mutable Mutability = iota
	Immutable
)

//go:generate mise x -- stringer -type=DirectionMode
type DirectionMode int

const (
	Directed DirectionMode = iota
	Undirected
)

//go:generate mise x -- stringer -type=SelfLoopsMode
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
	validate(t, mutability, directionMode, selfLoopsMode)

	newTester(
		t,
		graphBuilder,
		mutability,
		directionMode,
		selfLoopsMode,
	).test()
}

func validate(
	t TestingT,
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
	graphAdjacentNodesName = "Graph.AdjacentNodes"
	graphPredecessorsName  = "Graph.Predecessors"
	graphSuccessorsName    = "Graph.Successors"
)

func (tt tester) test() {
	tt.t.Helper()

	tt.testEmptyGraph()

	tt.testGraphWithOneNode()

	tt.testGraphWithTwoNodes()

	tt.testGraphWithOneEdge()

	tt.testGraphWithSameEdgePutTwice()

	tt.testGraphWithTwoEdgesWithSameSourceNode()
}

func (tt tester) testEmptyGraph() {
	tt.t.Helper()

	tt.t.Run("empty graph", func(t *testing.T) {
		t.Run("has no nodes", func(t *testing.T) {
			testNodeSet(t, graphNodesName, tt.graphBuilder().Nodes())
		})

		t.Run("has no edges", func(t *testing.T) {
			tt.testEdges(t, tt.graphBuilder())
		})
	})
}

func (tt tester) testGraphWithOneNode() {
	tt.t.Helper()

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
			testDegree(t, g(), node1, 0)
		})

		t.Run("the node has an in-degree of 0", func(t *testing.T) {
			testInDegree(t, g(), node1, 0)
		})

		t.Run("the node has an out-degree of 0", func(t *testing.T) {
			testOutDegree(t, g(), node1, 0)
		})
	})
}

func (tt tester) testGraphWithTwoNodes() {
	tt.t.Helper()

	tt.t.Run("graph with two nodes", func(t *testing.T) {
		t.Run("has both nodes", func(t *testing.T) {
			g := tt.graphBuilder()
			g = addNode(g, node1)
			g = addNode(g, node2)

			testNodeSet(t, graphNodesName, g.Nodes(), node1, node2)
		})
	})
}

func (tt tester) testGraphWithOneEdge() {
	tt.t.Helper()

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
			testDegree(t, g(), node1, 1)
		})

		t.Run("the target node has a degree of 1", func(t *testing.T) {
			testDegree(t, g(), node2, 1)
		})

		t.Run("the target node has an in-degree of 1", func(t *testing.T) {
			testInDegree(t, g(), node2, 1)
		})

		t.Run("the source node has an out-degree of 1", func(t *testing.T) {
			testOutDegree(t, g(), node1, 1)
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
				testHasEdgeConnecting(t, g(), node1, node2)
			},
		)

		t.Run(
			"sees the first node as being connected to no other node",
			func(t *testing.T) {
				testHasNoEdgeConnecting(t, g(), node1, nodeNotInGraph)
				testHasNoEdgeConnecting(t, g(), nodeNotInGraph, node1)
			},
		)

		t.Run(
			"sees the second node as being connected to no other node",
			func(t *testing.T) {
				testHasNoEdgeConnecting(t, g(), node2, nodeNotInGraph)
			},
		)
	})
}

func (tt tester) testGraphWithSameEdgePutTwice() {
	tt.t.Helper()

	tt.t.Run("graph with same edge put twice", func(t *testing.T) {
		t.Run("has only one edge", func(t *testing.T) {
			g := tt.graphBuilder()
			g = putEdge(g, node1, node2)

			tt.testEdges(t, g, graph.EndpointPairOf(node1, node2))
		})
	})
}

func (tt tester) testGraphWithTwoEdgesWithSameSourceNode() {
	tt.t.Helper()

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
				testDegree(t, g(), node1, 2)
			})

			t.Run("has a common node with two successors", func(t *testing.T) {
				testNodeSet(
					t,
					graphSuccessorsName,
					g().Successors(node1),
					node2,
					node3,
				)
			})

			t.Run(
				"has a common node with two unique adjacent nodes",
				func(t *testing.T) {
					testNodeSet(
						t,
						graphAdjacentNodesName,
						g().AdjacentNodes(node1),
						node2,
						node3,
					)
				},
			)

			t.Run("has a common with two edges", func(t *testing.T) {
				tt.testEdges(
					t,
					g(),
					graph.EndpointPairOf(node1, node2),
					graph.EndpointPairOf(node1, node3),
				)
			})

			t.Run("has a common with two incident edges", func(t *testing.T) {
				tt.testIncidentEdges(
					t,
					g(),
					node1,
					graph.EndpointPairOf(node1, node2),
					graph.EndpointPairOf(node1, node3),
				)
			})

			t.Run("has a common with an out-degree of 2", func(t *testing.T) {
				testOutDegree(t, g(), node1, 2)
			})

			// TODO: Add more tests, starting again from graph_test.go, "when
			//       putting two connected edges with the same target node".
		},
	)
}

func addNode(g graph.Graph[int], node int) graph.Graph[int] {
	if m, ok := g.(graph.MutableGraph[int]); ok {
		m.AddNode(node)
	}

	return g
}

func putEdge(g graph.Graph[int], source int, target int) graph.Graph[int] {
	if m, ok := g.(graph.MutableGraph[int]); ok {
		m.PutEdge(source, target)
	}

	return g
}

func complement(nodes []int) []int {
	all := []int{node1, node2, node3, nodeNotInGraph}
	return slices.DeleteFunc(all, func(value int) bool {
		return slices.Contains(nodes, value)
	})
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
	t.Helper()

	edgeSetTester{
		t:             t,
		setName:       "Graph.Edges",
		edges:         g.Edges(),
		directionMode: tt.directionMode,
		expectedEdges: expectedEdges,
	}.test()
}

func (tt tester) testIncidentEdges(
	t *testing.T,
	g graph.Graph[int],
	node int,
	expectedEdges ...graph.EndpointPair[int],
) {
	t.Helper()

	edgeSetTester{
		t:             t,
		setName:       "Graph.IncidentEdges",
		edges:         g.IncidentEdges(node),
		directionMode: tt.directionMode,
		expectedEdges: expectedEdges,
	}.test()
}

func testDegree(
	t *testing.T,
	g graph.Graph[int],
	node int,
	expectedDegree int,
) {
	t.Helper()

	if got, want := g.Degree(node), expectedDegree; got != want {
		t.Errorf("Graph.Degree: got %d, want %d", got, want)
	}
}

func testInDegree(
	t *testing.T,
	g graph.Graph[int],
	node int,
	expectedDegree int,
) {
	t.Helper()

	if got, want := g.InDegree(node), expectedDegree; got != want {
		t.Errorf("Graph.InDegree: got %d, want %d", got, want)
	}
}

func testOutDegree(
	t *testing.T,
	g graph.Graph[int],
	node int,
	expectedDegree int,
) {
	t.Helper()

	if got, want := g.OutDegree(node), expectedDegree; got != want {
		t.Errorf("Graph.OutDegree: got %d, want %d", got, want)
	}
}
