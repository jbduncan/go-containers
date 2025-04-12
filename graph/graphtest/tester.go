package graphtest

import (
	"slices"
	"strconv"
	"testing"

	"github.com/jbduncan/go-containers/graph"
	internalsettest "github.com/jbduncan/go-containers/internal/settest"
	"github.com/jbduncan/go-containers/set"
)

type Graph[N comparable] interface {
	IsDirected() bool
	AllowsSelfLoops() bool
	Nodes() set.Set[N]
	Edges() set.Set[graph.EndpointPair[N]]
	AdjacentNodes(node N) set.Set[N]
	Predecessors(node N) set.Set[N]
	Successors(node N) set.Set[N]
	IncidentEdges(node N) set.Set[graph.EndpointPair[N]]
	Degree(node N) int
	InDegree(node N) int
	OutDegree(node N) int
	HasEdgeConnecting(source N, target N) bool
	HasEdgeConnectingEndpoints(endpointPair graph.EndpointPair[N]) bool
	String() string
}

type MutableGraph[N comparable] interface {
	Graph[N]

	AddNode(node N) bool
	PutEdge(source N, target N) bool
	RemoveNode(node N) bool
	RemoveEdge(source N, target N) bool
}

const (
	node1          = 1
	node2          = 2
	node3          = 3
	nodeNotInGraph = 1_000
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

// TestReadOnly runs a suite of test cases for read-only Graph implementations.
// Graph instances created for testing are to have int nodes.
//
// Test cases that should be handled similarly in any graph implementation are
// included in this function; for example, testing that the Nodes method
// returns the set of the nodes in the graph. Details of specific
// implementations of the Graph interface are not tested.
//
// Parameter `emptyGraph` should always return a newly-initialized empty graph
// with no nodes and no edges. Otherwise, the behaviour of this function is
// undefined.
func TestReadOnly(
	t *testing.T,
	emptyGraph func() Graph[int],
	addNode func(g Graph[int], node int) Graph[int],
	putEdge func(g Graph[int], source int, target int) Graph[int],
	directionMode DirectionMode,
	selfLoopsMode SelfLoopsMode,
) {
	validate(t, directionMode, selfLoopsMode)

	newTester(
		t,
		emptyGraph,
		addNode,
		putEdge,
		false,
		directionMode == Directed,
		selfLoopsMode == AllowsSelfLoops,
	).test()
}

// TestMutable runs a suite of test cases for MutableGraph implementations.
// MutableGraph instances created for testing are to have int nodes.
//
// Test cases that should be handled similarly in any graph implementation are
// included in this function; for example, testing that the Nodes method
// returns the set of the nodes in the graph. Details of specific
// implementations of the MutableGraph interface are not tested.
//
// Parameter `emptyGraph` should always return a newly-initialized empty graph
// with no nodes and no edges. Otherwise, the behaviour of this function is
// undefined.
func TestMutable(
	t *testing.T,
	emptyGraph func() MutableGraph[int],
	directionMode DirectionMode,
	selfLoopsMode SelfLoopsMode,
) {
	validate(t, directionMode, selfLoopsMode)

	newTester(
		t,
		func() Graph[int] {
			return emptyGraph()
		},
		func(g Graph[int], node int) Graph[int] {
			g.(MutableGraph[int]).AddNode(node)
			return g
		},
		func(g Graph[int], source int, target int) Graph[int] {
			g.(MutableGraph[int]).PutEdge(source, target)
			return g
		},
		true,
		directionMode == Directed,
		selfLoopsMode == AllowsSelfLoops,
	).test()
}

func validate(
	t *testing.T,
	directionMode DirectionMode,
	selfLoopsMode SelfLoopsMode,
) {
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
	t *testing.T,
	emptyGraph func() Graph[int],
	addNode func(g Graph[int], node int) Graph[int],
	putEdge func(g Graph[int], source int, target int) Graph[int],
	mutable bool,
	directed bool,
	allowsSelfLoops bool,
) *tester {
	return &tester{
		t:               t,
		emptyGraph:      emptyGraph,
		addNode:         addNode,
		putEdge:         putEdge,
		mutable:         mutable,
		directed:        directed,
		allowsSelfLoops: allowsSelfLoops,
	}
}

type tester struct {
	t               *testing.T
	emptyGraph      func() Graph[int]
	addNode         func(g Graph[int], node int) Graph[int]
	putEdge         func(g Graph[int], source int, target int) Graph[int]
	mutable         bool
	directed        bool
	allowsSelfLoops bool
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

	if tt.mutable {
		tt.testMutableGraph()
	}

	if tt.directed {
		tt.testDirectedGraph()
	} else {
		tt.testUndirectedGraph()
	}

	if tt.allowsSelfLoops {
		tt.testSelfLoopingGraph()
	} else {
		tt.testSelfLoopDisallowingGraph()
	}

	if tt.mutable && tt.allowsSelfLoops {
		tt.testMutableSelfLoopingGraph()
	}

	if !tt.directed && tt.allowsSelfLoops {
		tt.testUndirectedSelfLoopingGraph()
	}

	if !tt.directed && !tt.allowsSelfLoops {
		tt.testUndirectedSelfLoopDisallowingGraph()
	}

	if tt.directed && tt.allowsSelfLoops {
		tt.testDirectedSelfLoopingGraph()
	}

	if tt.directed && !tt.allowsSelfLoops {
		tt.testDirectedSelfLoopDisallowingGraph()
	}
}

func (tt tester) testEmptyGraph() {
	tt.t.Run("empty graph", func(t *testing.T) {
		t.Run("has no nodes", func(t *testing.T) {
			testNodes(t, tt.emptyGraph())
		})

		t.Run("has no edges", func(t *testing.T) {
			tt.testEdges(t, tt.emptyGraph())
		})

		t.Run("has no predecessors for an absent node", func(t *testing.T) {
			testPredecessors(t, tt.emptyGraph(), nodeNotInGraph)
		})

		t.Run("has no successors for an absent node", func(t *testing.T) {
			testSuccessors(t, tt.emptyGraph(), nodeNotInGraph)
		})

		t.Run("has no adjacent nodes for an absent node", func(t *testing.T) {
			testAdjacentNodes(t, tt.emptyGraph(), nodeNotInGraph)
		})

		t.Run("has a degree of 0 for an absent node", func(t *testing.T) {
			testDegree(t, tt.emptyGraph(), nodeNotInGraph, 0)
		})

		t.Run("has an in-degree of 0 for an absent node", func(t *testing.T) {
			testInDegree(t, tt.emptyGraph(), nodeNotInGraph, 0)
		})

		t.Run("has an out-degree of 0 for an absent node", func(t *testing.T) {
			testOutDegree(t, tt.emptyGraph(), nodeNotInGraph, 0)
		})

		t.Run("has an unmodifiable nodes set view", func(t *testing.T) {
			g := tt.emptyGraph()
			nodes := g.Nodes()

			internalsettest.IsMutable(t, graphNodesName, nodes)

			_ = tt.addNode(g, node1)

			testNodeSet(t, graphNodesName, nodes, node1)
		})

		t.Run(
			"has an unmodifiable adjacent nodes set view",
			func(t *testing.T) {
				g := tt.emptyGraph()
				adjacentNodes := g.AdjacentNodes(node1)

				internalsettest.IsMutable(t, graphAdjacentNodesName, adjacentNodes)

				g = tt.putEdge(g, node1, node2)
				_ = tt.putEdge(g, node3, node1)

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
			g := tt.emptyGraph()
			predecessors := g.Predecessors(node1)

			internalsettest.IsMutable(t, graphPredecessorsName, predecessors)

			_ = tt.putEdge(g, node2, node1)

			testNodeSet(t, graphPredecessorsName, predecessors, node2)
		})

		t.Run("has an unmodifiable successors set view", func(t *testing.T) {
			g := tt.emptyGraph()
			successors := g.Successors(node1)

			internalsettest.IsMutable(t, graphSuccessorsName, successors)

			_ = tt.putEdge(g, node1, node2)

			testNodeSet(t, graphSuccessorsName, successors, node2)
		})

		t.Run("has an unmodifiable edges set view", func(t *testing.T) {
			g := tt.emptyGraph()
			edges := g.Edges()

			internalsettest.IsMutable(t, graphEdgesName, edges)

			_ = tt.putEdge(g, node1, node2)

			tt.testEdges(t, g, graph.EndpointPairOf(node1, node2))
		})

		t.Run(
			"has an unmodifiable incident edges set view",
			func(t *testing.T) {
				g := tt.emptyGraph()
				edges := g.IncidentEdges(node1)

				internalsettest.IsMutable(t, graphIncidentEdgesName, edges)

				_ = tt.putEdge(g, node1, node2)

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
		g := func() Graph[int] {
			g := tt.emptyGraph()
			g = tt.addNode(g, node1)
			return g
		}

		t.Run("has just that node", func(t *testing.T) {
			testNodes(t, g(), node1)
		})

		t.Run("the node has no adjacent nodes", func(t *testing.T) {
			testAdjacentNodes(t, g(), node1)
		})

		t.Run("the node has no predecessors", func(t *testing.T) {
			testPredecessors(t, g(), node1)
		})

		t.Run("the node has no successors", func(t *testing.T) {
			testSuccessors(t, g(), node1)
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
			g := tt.emptyGraph()
			g = tt.addNode(g, node1)
			g = tt.addNode(g, node2)

			testNodes(t, g, node1, node2)
		})
	})
}

func (tt tester) testGraphWithOneEdge() {
	tt.t.Run("graph with one edge", func(t *testing.T) {
		g := func() Graph[int] {
			g := tt.emptyGraph()
			g = tt.putEdge(g, node1, node2)
			return g
		}

		t.Run(
			"the source node is adjacent to the target node",
			func(t *testing.T) {
				testAdjacentNodes(t, g(), node1, node2)
			},
		)

		t.Run(
			"the target node is adjacent to the source node",
			func(t *testing.T) {
				testAdjacentNodes(t, g(), node2, node1)
			},
		)

		t.Run(
			"the source node is the predecessor of the target node",
			func(t *testing.T) {
				testPredecessors(t, g(), node2, node1)
			},
		)

		t.Run(
			"the target node is the successor of the source node",
			func(t *testing.T) {
				testSuccessors(t, g(), node1, node2)
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
			"connects the first node to the second",
			func(t *testing.T) {
				testHasEdgeConnecting(t, g(), node1, node2)
			},
		)

		t.Run(
			"connects the first node to no other node",
			func(t *testing.T) {
				testHasNoEdgeConnecting(t, g(), node1, nodeNotInGraph)
				testHasNoEdgeConnecting(t, g(), nodeNotInGraph, node1)
			},
		)

		t.Run(
			"connects the second node to no other node",
			func(t *testing.T) {
				testHasNoEdgeConnecting(t, g(), node2, nodeNotInGraph)
			},
		)
	})
}

func (tt tester) testGraphWithSameEdgePutTwice() {
	tt.t.Run("graph with same edge put twice", func(t *testing.T) {
		t.Run("has only one edge", func(t *testing.T) {
			g := tt.emptyGraph()
			g = tt.putEdge(g, node1, node2)

			tt.testEdges(t, g, graph.EndpointPairOf(node1, node2))
		})
	})
}

func (tt tester) testGraphWithTwoEdgesWithSameSourceNode() {
	tt.t.Run(
		"graph with two edges with the same source node",
		func(t *testing.T) {
			g := func() Graph[int] {
				g := tt.emptyGraph()
				g = tt.putEdge(g, node1, node2)
				g = tt.putEdge(g, node1, node3)
				return g
			}

			t.Run("has a common node with a degree of 2", func(t *testing.T) {
				testDegree(t, g(), node1, 2)
			})

			t.Run("has a common node with two successors", func(t *testing.T) {
				testSuccessors(t, g(), node1, node2, node3)
			})

			t.Run(
				"has a common node with two unique adjacent nodes",
				func(t *testing.T) {
					testAdjacentNodes(t, g(), node1, node2, node3)
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
			g := func() Graph[int] {
				g := tt.emptyGraph()
				g = tt.putEdge(g, node1, node2)
				g = tt.putEdge(g, node3, node2)
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
					testPredecessors(t, g(), node2, node1, node3)
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

func (tt tester) emptyMutableGraph() MutableGraph[int] {
	tt.t.Helper()

	g := tt.emptyGraph()

	mutG, ok := g.(MutableGraph[int])
	if !ok {
		tt.t.Fatalf(
			"graph was expected to implement graphtest.MutableGraph, " +
				"but it did not")
		return nil // Make the compiler happy
	}
	return mutG
}

func (tt tester) testMutableGraph() {
	tt.t.Run("mutable graph", func(t *testing.T) {
		tt.testMutableGraphAddingNewNode(t)

		tt.testMutableGraphAddingExistingNode(t)

		tt.testMutableGraphRemovingExistingNode(t)

		tt.testMutableGraphRemovingAbsentNode(t)

		tt.testMutableGraphPuttingNewEdge(t)

		tt.testMutableGraphPuttingExistingEdge(t)

		tt.testMutableGraphPuttingTwoAntiParallelEdges(t)

		tt.testMutableGraphRemovingExistingEdge(t)

		tt.testMutableGraphRemovingAbsentEdgeWithExistingSource(t)

		tt.testMutableGraphRemovingAbsentEdgeWithExistingTarget(t)

		tt.testMutableGraphRemovingAbsentEdgeWithTwoExistingNodes(t)
	})
}

func (tt tester) testMutableGraphAddingNewNode(t *testing.T) {
	t.Run("adding a new node returns true", func(t *testing.T) {
		if got := tt.emptyMutableGraph().AddNode(node1); !got {
			t.Fatalf("MutableGraph.AddNode: got false, want true")
		}
	})
}

func (tt tester) testMutableGraphAddingExistingNode(t *testing.T) {
	t.Run("adding an existing node returns false", func(t *testing.T) {
		g := tt.emptyMutableGraph()
		g.AddNode(node1)

		if got := g.AddNode(node1); got {
			t.Fatalf("MutableGraph.AddNode: got true, want false")
		}
	})
}

func (tt tester) testMutableGraphRemovingExistingNode(t *testing.T) {
	t.Run("removing an existing node", func(t *testing.T) {
		setup := func() (g MutableGraph[int], removed bool) {
			g = tt.emptyMutableGraph()
			g.PutEdge(node1, node2)
			g.PutEdge(node3, node1)
			g.PutEdge(node2, node3)
			removed = g.RemoveNode(node1)
			return
		}

		t.Run("returns true", func(t *testing.T) {
			_, removed := setup()

			if got := removed; !got {
				t.Fatalf("MutableGraph.RemoveNode: got false, want true")
			}
		})

		t.Run("leaves the other nodes alone", func(t *testing.T) {
			g, _ := setup()

			testNodes(t, g, node2, node3)
		})

		t.Run("detaches it from its adjacent nodes", func(t *testing.T) {
			g, _ := setup()

			testAdjacentNodes(t, g, node2, node3)
			testAdjacentNodes(t, g, node3, node2)
		})

		t.Run("removes the connected edges", func(t *testing.T) {
			g, _ := setup()

			tt.testEdges(t, g, graph.EndpointPairOf(node2, node3))
		})
	})
}

func (tt tester) testMutableGraphRemovingAbsentNode(t *testing.T) {
	t.Run("removing an absent node", func(t *testing.T) {
		setup := func() (g MutableGraph[int], removed bool) {
			g = tt.emptyMutableGraph()
			g.AddNode(node1)
			removed = g.RemoveNode(nodeNotInGraph)
			return
		}

		t.Run("returns false", func(t *testing.T) {
			_, removed := setup()

			if got := removed; got {
				t.Fatalf("MutableGraph.RemoveNode: got true, want false")
			}
		})

		t.Run("leaves all the nodes alone", func(t *testing.T) {
			g, _ := setup()

			testNodes(t, g, node1)
		})
	})
}

func (tt tester) testMutableGraphPuttingNewEdge(t *testing.T) {
	t.Run("putting a new edge returns true", func(t *testing.T) {
		if got := tt.emptyMutableGraph().PutEdge(node1, node2); !got {
			t.Fatalf("MutableGraph.PutEdge: got false, want true")
		}
	})
}

func (tt tester) testMutableGraphPuttingExistingEdge(t *testing.T) {
	t.Run("putting an existing edge returns false", func(t *testing.T) {
		g := tt.emptyMutableGraph()
		g.PutEdge(node1, node2)

		if got := g.PutEdge(node1, node2); got {
			t.Fatalf("MutableGraph.PutEdge: got true, want false")
		}
	})
}

func (tt tester) testMutableGraphPuttingTwoAntiParallelEdges(t *testing.T) {
	t.Run(
		"putting two anti-parallel edges and removing one of the nodes",
		func(t *testing.T) {
			setup := func() MutableGraph[int] {
				g := tt.emptyMutableGraph()
				g.PutEdge(node1, node2)
				g.PutEdge(node2, node1)
				g.RemoveNode(node1)

				return g
			}

			t.Run("leaves the other node alone", func(t *testing.T) {
				g := setup()

				testNodes(t, g, node2)
			})

			t.Run("removes both edges", func(t *testing.T) {
				g := setup()

				tt.testEdges(t, g)
			})
		},
	)
}

func (tt tester) testMutableGraphRemovingExistingEdge(t *testing.T) {
	t.Run(
		"removing an existing edge",
		func(t *testing.T) {
			setup := func() (g MutableGraph[int], removed bool) {
				g = tt.emptyMutableGraph()
				g.PutEdge(node1, node2)
				g.PutEdge(node1, node3)
				removed = g.RemoveEdge(node1, node2)
				return
			}

			t.Run("returns true", func(t *testing.T) {
				_, removed := setup()

				if got := removed; !got {
					t.Fatalf("MutableGraph.RemoveEdge: got false, want true")
				}
			})

			t.Run("detaches the two nodes", func(t *testing.T) {
				g, _ := setup()

				testSuccessors(t, g, node1, node3)
				testPredecessors(t, g, node3, node1)
				testPredecessors(t, g, node2)
			})

			t.Run("leaves the other edges alone", func(t *testing.T) {
				g, _ := setup()

				tt.testEdges(t, g, graph.EndpointPairOf(node1, node3))
			})
		},
	)
}

func (tt tester) testMutableGraphRemovingAbsentEdgeWithExistingSource(
	t *testing.T,
) {
	t.Run(
		"removing an absent edge with an existing source",
		func(t *testing.T) {
			setup := func() (g MutableGraph[int], removed bool) {
				g = tt.emptyMutableGraph()
				g.PutEdge(node1, node2)
				removed = g.RemoveEdge(node1, nodeNotInGraph)
				return
			}

			t.Run("returns false", func(t *testing.T) {
				_, removed := setup()

				if got := removed; got {
					t.Fatalf("MutableGraph.RemoveEdge: got true, want false")
				}
			})

			t.Run("leaves the existing nodes alone", func(t *testing.T) {
				g, _ := setup()

				testSuccessors(t, g, node1, node2)
				testPredecessors(t, g, node2, node1)
			})
		},
	)
}

func (tt tester) testMutableGraphRemovingAbsentEdgeWithExistingTarget(
	t *testing.T,
) {
	t.Run(
		"removing an absent edge with an existing target",
		func(t *testing.T) {
			setup := func() (g MutableGraph[int], removed bool) {
				g = tt.emptyMutableGraph()
				g.PutEdge(node1, node2)
				removed = g.RemoveEdge(nodeNotInGraph, node2)
				return
			}

			t.Run("returns false", func(t *testing.T) {
				_, removed := setup()

				if got := removed; got {
					t.Fatalf("MutableGraph.RemoveEdge: got true, want false")
				}
			})

			t.Run("leaves the existing nodes alone", func(t *testing.T) {
				g, _ := setup()

				testSuccessors(t, g, node1, node2)
				testPredecessors(t, g, node2, node1)
			})
		},
	)
}

func (tt tester) testMutableGraphRemovingAbsentEdgeWithTwoExistingNodes(
	t *testing.T,
) {
	t.Run(
		"removing an absent edge with two existing nodes",
		func(t *testing.T) {
			setup := func() (g MutableGraph[int], removed bool) {
				g = tt.emptyMutableGraph()
				g.AddNode(node1)
				g.AddNode(node2)
				removed = g.RemoveEdge(node1, node2)
				return
			}

			t.Run("returns false", func(t *testing.T) {
				_, removed := setup()

				if got := removed; got {
					t.Fatalf("MutableGraph.RemoveEdge: got true, want false")
				}
			})

			t.Run("leaves the existing nodes alone", func(t *testing.T) {
				g, _ := setup()

				testNodes(t, g, node1, node2)
			})
		},
	)
}

func (tt tester) testDirectedGraph() {
	tt.t.Run("directed graph", func(t *testing.T) {
		t.Run("says it is directed", func(t *testing.T) {
			g := tt.emptyGraph()

			if got := g.IsDirected(); !got {
				t.Fatalf("Graph.IsDirected: got false, want true")
			}
		})

		t.Run("putting an edge", func(t *testing.T) {
			g := func() Graph[int] {
				g := tt.emptyGraph()
				g = tt.putEdge(g, node1, node2)
				return g
			}

			t.Run(
				"makes the first node have no predecessors",
				func(t *testing.T) {
					testPredecessors(t, g(), node1)
				},
			)

			t.Run(
				"makes the second node have no successors",
				func(t *testing.T) {
					testSuccessors(t, g(), node2)
				},
			)

			t.Run(
				"makes the first node have an in-degree of 0",
				func(t *testing.T) {
					testInDegree(t, g(), node1, 0)
				},
			)

			t.Run(
				"makes the second node have an out-degree of 0",
				func(t *testing.T) {
					testOutDegree(t, g(), node2, 0)
				},
			)

			t.Run("does not connect the second node to the first", func(t *testing.T) {
				testHasNoEdgeConnecting(t, g(), node2, node1)
			})
		})

		t.Run(
			"putting two connected edges makes the common node have a "+
				"degree of 2",
			func(t *testing.T) {
				g := tt.emptyGraph()
				g = tt.putEdge(g, node1, node2)
				g = tt.putEdge(g, node2, node3)

				testDegree(t, g, node2, 2)
			},
		)
	})
}

func (tt tester) testUndirectedGraph() {
	tt.t.Run("undirected graph", func(t *testing.T) {
		t.Run("says it is not directed", func(t *testing.T) {
			g := tt.emptyGraph()

			if got := g.IsDirected(); got {
				t.Fatalf("Graph.IsDirected: got true, want false")
			}
		})

		t.Run("putting an edge", func(t *testing.T) {
			g := func() Graph[int] {
				g := tt.emptyGraph()
				g = tt.putEdge(g, node1, node2)
				return g
			}

			t.Run(
				"makes the first node the predecessor of the second",
				func(t *testing.T) {
					testPredecessors(t, g(), node1, node2)
				},
			)

			t.Run(
				"makes the first node a predecessor of the second",
				func(t *testing.T) {
					testPredecessors(t, g(), node1, node2)
				},
			)

			t.Run(
				"makes the second node a successor of the first",
				func(t *testing.T) {
					testSuccessors(t, g(), node2, node1)
				},
			)

			t.Run(
				"makes the first node have an in-degree of 1",
				func(t *testing.T) {
					testInDegree(t, g(), node1, 1)
				},
			)

			t.Run(
				"makes the second node have an out-degree of 1",
				func(t *testing.T) {
					testOutDegree(t, g(), node2, 1)
				},
			)

			t.Run("connects the second node to the first", func(t *testing.T) {
				testHasEdgeConnecting(t, g(), node2, node1)
			})
		})
	})
}

func (tt tester) testSelfLoopingGraph() {
	tt.t.Run("self-looping graph", func(t *testing.T) {
		t.Run("says it allows self loops", func(t *testing.T) {
			g := tt.emptyGraph()

			if got := g.AllowsSelfLoops(); !got {
				t.Fatalf("Graph.AllowsSelfLoops: got false, want true")
			}
		})

		t.Run("putting a self-loop edge", func(t *testing.T) {
			g := func() Graph[int] {
				g := tt.emptyGraph()
				g = tt.putEdge(g, node1, node1)
				return g
			}

			t.Run(
				"makes the shared node its own adjacent node",
				func(t *testing.T) {
					testAdjacentNodes(t, g(), node1, node1)
				},
			)

			t.Run(
				"makes the shared node have a degree of 2 because the "+
					"edge touches the node twice",
				func(t *testing.T) {
					testDegree(t, g(), node1, 2)
				},
			)
		})
	})
}

func (tt tester) testSelfLoopDisallowingGraph() {
	tt.t.Run("self-loop-disallowing graph", func(t *testing.T) {
		t.Run(
			"says it disallows self-loops",
			func(t *testing.T) {
				g := tt.emptyGraph()

				if got := g.AllowsSelfLoops(); got {
					t.Fatalf("Graph.AllowsSelfLoops: got true, want false")
				}
			},
		)
	})
}

func (tt tester) testMutableSelfLoopingGraph() {
	tt.t.Run("mutable self-looping graph", func(t *testing.T) {
		t.Run(
			"removing a self-looping node removes the self-loop edge",
			func(t *testing.T) {
				g := tt.emptyMutableGraph()
				g.PutEdge(node1, node1)
				g.RemoveNode(node1)

				tt.testEdges(t, g)
			},
		)
	})
}

func (tt tester) testUndirectedSelfLoopingGraph() {
	tt.t.Run("undirected self-looping graph", func(t *testing.T) {
		tt.testStringRepresentations(t, false, true)
	})
}

func (tt tester) testUndirectedSelfLoopDisallowingGraph() {
	tt.t.Run("undirected self-loop-disallowing graph", func(t *testing.T) {
		tt.testStringRepresentations(t, false, false)
	})
}

func (tt tester) testDirectedSelfLoopingGraph() {
	tt.t.Run("directed self-looping graph", func(t *testing.T) {
		tt.testStringRepresentations(t, true, true)
	})
}

func (tt tester) testDirectedSelfLoopDisallowingGraph() {
	tt.t.Run("directed self-loop-disallowing graph", func(t *testing.T) {
		tt.testStringRepresentations(t, true, false)
	})
}

//nolint:revive
func (tt tester) testStringRepresentations(
	t *testing.T,
	directed bool,
	allowsSelfLoops bool,
) {
	t.Run("has an empty graph string representation", func(t *testing.T) {
		want := "isDirected: " +
			strconv.FormatBool(directed) +
			", allowsSelfLoops: " +
			strconv.FormatBool(allowsSelfLoops) +
			", nodes: [], edges: []"
		if got := tt.emptyGraph().String(); got != want {
			t.Errorf("Graph.String: got %q, want %q", got, want)
		}
	})

	t.Run(
		"adding a node makes a non-empty graph string representation",
		func(t *testing.T) {
			g := tt.emptyGraph()
			g = tt.addNode(g, node1)

			want := "isDirected: " +
				strconv.FormatBool(directed) +
				", allowsSelfLoops: " +
				strconv.FormatBool(allowsSelfLoops) +
				", nodes: [1], edges: []"
			if got := g.String(); got != want {
				t.Errorf("Graph.String: got %q, want %q", got, want)
			}
		},
	)

	t.Run(
		"putting an edge makes a non-empty graph string representation",
		func(t *testing.T) {
			g := tt.emptyGraph()
			g = tt.putEdge(g, node1, node2)

			var wantAny []string
			for _, nodes := range []string{"[1, 2]", "[2, 1]"} {
				var wantAnyEdges []string
				if directed {
					wantAnyEdges = []string{"[<1 -> 2>]"}
				} else {
					wantAnyEdges = []string{"[<1 -> 2>]", "[<2 -> 1>]"}
				}
				for _, edges := range wantAnyEdges {
					wantAny = append(
						wantAny,
						"isDirected: "+
							strconv.FormatBool(directed)+
							", allowsSelfLoops: "+
							strconv.FormatBool(allowsSelfLoops)+
							", nodes: "+
							nodes+
							", edges: "+
							edges,
					)
				}
			}
			if got := g.String(); !slices.Contains(wantAny, got) {
				t.Errorf("Graph.String: got %q, want any of %q", got, wantAny)
			}
		},
	)
}

func complement(nodes []int) []int {
	all := []int{node1, node2, node3, nodeNotInGraph}
	return slices.DeleteFunc(all, func(value int) bool {
		return slices.Contains(nodes, value)
	})
}

func testNodes(
	t *testing.T,
	g Graph[int],
	expectedValues ...int,
) {
	t.Helper()

	testNodeSet(t, graphNodesName, g.Nodes(), expectedValues...)
}

func testAdjacentNodes(
	t *testing.T,
	g Graph[int],
	node int,
	expectedValues ...int,
) {
	t.Helper()

	testNodeSet(
		t,
		graphAdjacentNodesName,
		g.AdjacentNodes(node),
		expectedValues...,
	)
}

func testPredecessors(
	t *testing.T,
	g Graph[int],
	node int,
	expectedValues ...int,
) {
	t.Helper()

	testNodeSet(
		t,
		graphPredecessorsName,
		g.Predecessors(node),
		expectedValues...,
	)
}

func testSuccessors(
	t *testing.T,
	g Graph[int],
	node int,
	expectedValues ...int,
) {
	t.Helper()

	testNodeSet(t, graphSuccessorsName, g.Successors(node), expectedValues...)
}

func testNodeSet(
	t *testing.T,
	setName string,
	s set.Set[int],
	expectedValues ...int,
) {
	t.Helper()

	internalsettest.Len(t, setName, s, len(expectedValues))
	internalsettest.All(t, setName, s, expectedValues)
	internalsettest.Contains(t, setName, s, expectedValues)
	internalsettest.DoesNotContain(t, setName, s, complement(expectedValues))
	internalsettest.String(t, setName, s, expectedValues)
}

func (tt tester) testEdges(
	t *testing.T,
	g Graph[int],
	expectedEdges ...graph.EndpointPair[int],
) {
	t.Helper()

	edgeSetTester{
		t:             t,
		setName:       graphEdgesName,
		edges:         g.Edges(),
		directed:      tt.directed,
		expectedEdges: expectedEdges,
	}.test()
}

func (tt tester) testIncidentEdges(
	t *testing.T,
	g Graph[int],
	node int,
	expectedEdges ...graph.EndpointPair[int],
) {
	t.Helper()

	edgeSetTester{
		t:             t,
		setName:       graphIncidentEdgesName,
		edges:         g.IncidentEdges(node),
		directed:      tt.directed,
		expectedEdges: expectedEdges,
	}.test()
}

func testDegree(
	t *testing.T,
	g Graph[int],
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
	g Graph[int],
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
	g Graph[int],
	node int,
	expectedDegree int,
) {
	t.Helper()

	if got, want := g.OutDegree(node), expectedDegree; got != want {
		t.Errorf("Graph.OutDegree: got %d, want %d", got, want)
	}
}

func testHasEdgeConnecting(
	t *testing.T,
	g Graph[int],
	source, target int,
) {
	if got := g.HasEdgeConnecting(source, target); !got {
		t.Errorf("Graph.HasEdgeConnecting: got false, want true")
	}
	if got := g.HasEdgeConnectingEndpoints(
		graph.EndpointPairOf(source, target),
	); !got {
		t.Errorf("Graph.HasEdgeConnectingEndpoints: got false, want true")
	}
}

func testHasNoEdgeConnecting(
	t *testing.T,
	g Graph[int],
	source, target int,
) {
	if got := g.HasEdgeConnecting(source, target); got {
		t.Errorf("Graph.HasEdgeConnecting: got true, want false")
	}
	if got := g.HasEdgeConnectingEndpoints(
		graph.EndpointPairOf(source, target),
	); got {
		t.Errorf("Graph.HasEdgeConnectingEndpoints: got true, want false")
	}
}
