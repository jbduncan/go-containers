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

// TODO: Rename to `emptyGraph` and add a note to the docs about how this
//       function should always return a newly initialized, empty graph.

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
	graphEdgesName         = "Graph.Edges"
	graphIncidentEdgesName = "Graph.IncidentEdges"
)

func (tt tester) test() {
	tt.testEmptyGraph()

	tt.testGraphWithOneNode()

	tt.testGraphWithTwoNodes()

	tt.testGraphWithOneEdge()

	tt.testGraphWithSameEdgePutTwice()

	tt.testGraphWithTwoEdgesWithSameSourceNode()

	tt.testGraphWithTwoEdgesWithSameTargetNode()

	tt.testEmptyMutableGraph()
}

func (tt tester) testEmptyGraph() {
	tt.t.Run("empty graph", func(t *testing.T) {
		t.Run("has no nodes", func(t *testing.T) {
			testNodeSet(t, graphNodesName, tt.graphBuilder().Nodes())
		})

		t.Run("has no edges", func(t *testing.T) {
			tt.testEdges(t, tt.graphBuilder())
		})

		t.Run("has no predecessors for an absent node", func(t *testing.T) {
			testNodeSet(
				t,
				graphPredecessorsName,
				tt.graphBuilder().Predecessors(nodeNotInGraph),
			)
		})

		t.Run("has no successors for an absent node", func(t *testing.T) {
			testNodeSet(
				t,
				graphSuccessorsName,
				tt.graphBuilder().Successors(nodeNotInGraph),
			)
		})

		t.Run("has no adjacent nodes for an absent node", func(t *testing.T) {
			testNodeSet(
				t,
				graphAdjacentNodesName,
				tt.graphBuilder().AdjacentNodes(nodeNotInGraph),
			)
		})

		t.Run("has a degree of 0 for an absent node", func(t *testing.T) {
			testDegree(t, tt.graphBuilder(), nodeNotInGraph, 0)
		})

		t.Run("has an in-degree of 0 for an absent node", func(t *testing.T) {
			testInDegree(t, tt.graphBuilder(), nodeNotInGraph, 0)
		})

		t.Run("has an out-degree of 0 for an absent node", func(t *testing.T) {
			testOutDegree(t, tt.graphBuilder(), nodeNotInGraph, 0)
		})

		t.Run("has an unmodifiable nodes set view", func(t *testing.T) {
			g := tt.graphBuilder()
			nodes := g.Nodes()

			if _, mutable := nodes.(set.MutableSet[int]); mutable {
				t.Fatalf(
					"%s: got a set.MutableSet: %v, want just a set.Set",
					graphNodesName,
					nodes,
				)
			}

			_ = addNode(g, node1)

			testNodeSet(t, graphNodesName, nodes, node1)
		})

		t.Run(
			"has an unmodifiable adjacent nodes set view",
			func(t *testing.T) {
				g := tt.graphBuilder()
				adjacentNodes := g.AdjacentNodes(node1)

				testSetIsMutable(t, adjacentNodes, graphAdjacentNodesName)

				g = putEdge(g, node1, node2)
				_ = putEdge(g, node3, node1)

				testNodeSet(
					t,
					graphAdjacentNodesName,
					adjacentNodes,
					node2,
					node3,
				)
			},
		)

		t.Run("has an unmodifiable predecessors set view", func(t *testing.T) {
			g := tt.graphBuilder()
			predecessors := g.Predecessors(node1)

			testSetIsMutable(t, predecessors, graphPredecessorsName)

			_ = putEdge(g, node2, node1)

			testNodeSet(t, graphPredecessorsName, predecessors, node2)
		})

		t.Run("has an unmodifiable successors set view", func(t *testing.T) {
			g := tt.graphBuilder()
			successors := g.Successors(node1)

			testSetIsMutable(t, successors, graphSuccessorsName)

			_ = putEdge(g, node1, node2)

			testNodeSet(t, graphSuccessorsName, successors, node2)
		})

		t.Run("has an unmodifiable edges set view", func(t *testing.T) {
			g := tt.graphBuilder()
			edges := g.Edges()

			testSetIsMutable(t, edges, graphEdgesName)

			_ = putEdge(g, node1, node2)

			tt.testEdges(t, g, graph.EndpointPairOf(node1, node2))
		})

		t.Run(
			"has an unmodifiable incident edges set view",
			func(t *testing.T) {
				g := tt.graphBuilder()
				edges := g.IncidentEdges(node1)

				testSetIsMutable(t, edges, graphIncidentEdgesName)

				_ = putEdge(g, node1, node2)

				tt.testIncidentEdges(
					t,
					g,
					node1,
					graph.EndpointPairOf(node1, node2),
				)
			},
		)
	})
}

func (tt tester) testGraphWithOneNode() {
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
	tt.t.Run("graph with same edge put twice", func(t *testing.T) {
		t.Run("has only one edge", func(t *testing.T) {
			g := tt.graphBuilder()
			g = putEdge(g, node1, node2)

			tt.testEdges(t, g, graph.EndpointPairOf(node1, node2))
		})
	})
}

func (tt tester) testGraphWithTwoEdgesWithSameSourceNode() {
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
		},
	)
}

func (tt tester) testGraphWithTwoEdgesWithSameTargetNode() {
	tt.t.Run(
		"graph with two edges with the same target node",
		func(t *testing.T) {
			g := func() graph.Graph[int] {
				g := tt.graphBuilder()
				g = putEdge(g, node1, node2)
				g = putEdge(g, node3, node2)
				return g
			}

			t.Run(
				"has a common node with an in-degree of 2",
				func(t *testing.T) {
					testInDegree(t, g(), node2, 2)
				},
			)

			t.Run(
				"has a common node with two predecessors",
				func(t *testing.T) {
					testNodeSet(
						t,
						graphPredecessorsName,
						g().Predecessors(node2),
						node1,
						node3,
					)
				},
			)

			t.Run("has a common with two incident edges", func(t *testing.T) {
				tt.testIncidentEdges(
					t,
					g(),
					node2,
					graph.EndpointPairOf(node1, node2),
					graph.EndpointPairOf(node3, node2),
				)
			})
		},
	)
}

func (tt tester) testEmptyMutableGraph() {
	tt.t.Run("empty mutable graph", func(t *testing.T) {
		emptyMutableGraph := func() graph.MutableGraph[int] {
			tt.t.Helper()

			g := tt.graphBuilder()

			mutG, ok := g.(graph.MutableGraph[int])
			if !ok {
				tt.t.Fatalf(
					"graph was expected to implement graph.MutableGraph, " +
						"but it did not")
				return nil // Make the compiler happy
			}
			return mutG
		}

		t.Run("adding a new node returns true", func(t *testing.T) {
			if got := emptyMutableGraph().AddNode(node1); !got {
				t.Fatalf("MutableGraph.AddNode: got false, want true")
			}
		})

		t.Run("adding an existing node returns false", func(t *testing.T) {
			g := emptyMutableGraph()
			g.AddNode(node1)

			if got := g.AddNode(node1); got {
				t.Fatalf("MutableGraph.AddNode: got true, want false")
			}
		})

		t.Run("removing an existing node", func(t *testing.T) {
			setup := func() (g graph.MutableGraph[int], removed bool) {
				g = emptyMutableGraph()
				g.PutEdge(node1, node2)
				g.PutEdge(node3, node1)
				g.PutEdge(node2, node3)

				return g, g.RemoveNode(node1)
			}

			t.Run("returns true", func(t *testing.T) {
				_, removed := setup()

				if got := removed; !got {
					t.Fatalf("MutableGraph.RemoveNode: got false, want true")
				}
			})

			t.Run("leaves the other nodes alone", func(t *testing.T) {
				g, _ := setup()

				testNodeSet(t, graphNodesName, g.Nodes(), node2, node3)
			})

			t.Run("detaches it from its adjacent nodes", func(t *testing.T) {
				g, _ := setup()

				testNodeSet(
					t,
					graphAdjacentNodesName,
					g.AdjacentNodes(node2),
					node3,
				)
				testNodeSet(
					t,
					graphAdjacentNodesName,
					g.AdjacentNodes(node3),
					node2,
				)
			})

			t.Run("removes the connected edges", func(t *testing.T) {
				g, _ := setup()

				tt.testEdges(t, g, graph.EndpointPairOf(node2, node3))
			})
		})

		t.Run("removing an absent node", func(t *testing.T) {
			setup := func() (g graph.MutableGraph[int], removed bool) {
				g = emptyMutableGraph()
				g.AddNode(node1)

				return g, g.RemoveNode(nodeNotInGraph)
			}

			t.Run("returns false", func(t *testing.T) {
				_, removed := setup()

				if got := removed; got {
					t.Fatalf("MutableGraph.RemoveNode: got true, want false")
				}
			})

			t.Run("leaves all the nodes alone", func(t *testing.T) {
				g, _ := setup()

				testNodeSet(t, graphNodesName, g.Nodes(), node1)
			})
		})
	})

	// TODO: continue from graph_test.go, line 550, "when putting a new edge".
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
		setName:       graphEdgesName,
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
		setName:       graphIncidentEdgesName,
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

func testSetIsMutable[T comparable](
	t *testing.T,
	s set.Set[T],
	setName string,
) {
	t.Helper()

	if _, mutable := s.(set.MutableSet[T]); mutable {
		t.Fatalf(
			"%s: got a set.MutableSet: %v, want just a set.Set",
			setName,
			s,
		)
	}
}
