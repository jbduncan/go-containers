package graph_test

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/jbduncan/go-containers/graph"
	. "github.com/jbduncan/go-containers/internal/matchers"
	"github.com/jbduncan/go-containers/set"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Graphs", func() {
	graphTests(
		"graph.Undirected[int]().Build()",
		func() graph.Graph[int] {
			return graph.Undirected[int]().Build()
		},
		Mutable,
		Undirected,
		DisallowsSelfLoops)
	graphTests(
		"graph.Undirected[int]().AllowsSelfLoops(true).Build()",
		func() graph.Graph[int] {
			return graph.Undirected[int]().AllowsSelfLoops(true).Build()
		},
		Mutable,
		Undirected,
		AllowsSelfLoops)
	graphTests(
		"graph.Directed[int]().Build()",
		func() graph.Graph[int] {
			return graph.Directed[int]().Build()
		},
		Mutable,
		Directed,
		DisallowsSelfLoops)
	graphTests(
		"graph.Directed[int]().AllowsSelfLoops(true).Build()",
		func() graph.Graph[int] {
			return graph.Directed[int]().AllowsSelfLoops(true).Build()
		},
		Mutable,
		Directed,
		AllowsSelfLoops)
})

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

func addNode(grph graph.Graph[int], node int) graph.Graph[int] {
	if graphAsMutable, ok := grph.(graph.MutableGraph[int]); ok {
		graphAsMutable.AddNode(node)
	}

	return grph
}

func putEdge(grph graph.Graph[int], node1 int, node2 int) graph.Graph[int] {
	if graphAsMutable, ok := grph.(graph.MutableGraph[int]); ok {
		graphAsMutable.PutEdge(node1, node2)
	}

	return grph
}

// graphTests produces a suite of Ginkgo test cases for testing implementations of the Graph and
// MutableGraph interfaces. Graph instances created for testing are to have int nodes.
//
// Test cases that should be handled similarly in any graph implementation are included in this
// function; for example, testing that the `Nodes()` method returns the set of the nodes in the
// graph. Details of specific implementations of the Graph and MutableGraph interfaces are
// explicitly not tested.
func graphTests(
	graphName string,
	createGraph func() graph.Graph[int],
	mutability Mutability,
	directionMode DirectionMode,
	selfLoopsMode SelfLoopsMode,
) {
	Context(fmt.Sprintf("%s: given a graph", graphName), func() {
		var grph graph.Graph[int]

		BeforeEach(func() {
			grph = createGraph()
		})

		It("has no nodes", func() {
			testSet(grph.Nodes())
		})

		It("has no edges", func() {
			testEmptyEdges(grph.Edges())
		})

		Context("when adding one node", func() {
			BeforeEach(func() {
				grph = addNode(grph, node1)
			})

			It("has just that node", func() {
				testSet(grph.Nodes(), node1)
			})

			It("reports that the node has no adjacent nodes", func() {
				testSet(grph.AdjacentNodes(node1))
			})

			It("reports that the node has no predecessors", func() {
				testSet(grph.Predecessors(node1))
			})

			It("reports that the node has no successors", func() {
				testSet(grph.Successors(node1))
			})

			It("reports that the node has no incident edges", func() {
				testEmptyEdges(grph.IncidentEdges(node1))
			})

			It("reports that the node has a degree of 0", func() {
				Expect(grph.Degree(node1)).To(BeZero())
			})

			It("reports that the node has an in degree of 0", func() {
				Expect(grph.InDegree(node1)).To(BeZero())
			})

			It("reports that the node has an out degree of 0", func() {
				Expect(grph.OutDegree(node1)).To(BeZero())
			})
		})

		Context("when adding two nodes", func() {
			BeforeEach(func() {
				grph = addNode(grph, node1)
				grph = addNode(grph, node2)
			})

			It("has both nodes", func() {
				testSet(grph.Nodes(), node1, node2)
			})
		})

		Context("when putting one edge", func() {
			BeforeEach(func() {
				grph = putEdge(grph, node1, node2)
			})

			It("reports that both nodes are adjacent to each other", func() {
				testSet(grph.AdjacentNodes(node1), node2)
				testSet(grph.AdjacentNodes(node2), node1)
			})

			It("sees that the second node has the first node as a predecessor", func() {
				testSet(grph.Predecessors(node2), node1)
			})

			It("sees that the first node has the second node as a successor", func() {
				testSet(grph.Successors(node1), node2)
			})

			It("reports that both nodes have a degree of 1", func() {
				Expect(grph.Degree(node1)).To(Equal(1))
				Expect(grph.Degree(node2)).To(Equal(1))
			})

			It("has an in degree of 1 for the second node", func() {
				Expect(grph.InDegree(node2)).To(Equal(1))
			})

			It("has an out degree of 1 for the first node", func() {
				Expect(grph.OutDegree(node1)).To(Equal(1))
			})

			It("has an incident edge connecting the first node to the second node", func() {
				testSingleEdge(
					grph,
					grph.IncidentEdges(node1),
					graph.EndpointPairOf(node1, node2))
			})

			It("has an incident edge connecting the second node to the first node", func() {
				testSingleEdge(
					grph,
					grph.IncidentEdges(node2),
					graph.EndpointPairOf(node1, node2))
			})

			It("has just one edge", func() {
				testSingleEdge(
					grph,
					grph.Edges(),
					graph.EndpointPairOf(node1, node2))
			})

			It("sees the first node as being connected to the second", func() {
				Expect(grph.HasEdgeConnecting(node1, node2)).
					To(BeTrue())

				Expect(grph.HasEdgeConnectingEndpoints(graph.EndpointPairOf(node1, node2))).
					To(BeTrue())
			})

			It("does not see the first node as being connected to any other node", func() {
				Expect(grph.HasEdgeConnecting(node1, nodeNotInGraph)).
					To(BeFalse())
				Expect(grph.HasEdgeConnecting(nodeNotInGraph, node1)).
					To(BeFalse())

				Expect(grph.HasEdgeConnectingEndpoints(graph.EndpointPairOf(node1, nodeNotInGraph))).
					To(BeFalse())
				Expect(grph.HasEdgeConnectingEndpoints(graph.EndpointPairOf(nodeNotInGraph, node1))).
					To(BeFalse())
			})

			It("does not see the second node as being connected to any other node", func() {
				Expect(grph.HasEdgeConnecting(node2, nodeNotInGraph)).
					To(BeFalse())

				Expect(grph.HasEdgeConnectingEndpoints(graph.EndpointPairOf(node2, nodeNotInGraph))).
					To(BeFalse())
			})
		})

		Context("when putting the same edge twice", func() {
			It("has only one edge", func() {
				grph = putEdge(grph, node1, node2)
				grph = putEdge(grph, node1, node2)

				testSingleEdge(
					grph,
					grph.Edges(),
					graph.EndpointPairOf(node1, node2))
			})
		})

		Context("when putting two connected edges with the same source node", func() {
			BeforeEach(func() {
				grph = putEdge(grph, node1, node2)
				grph = putEdge(grph, node1, node3)
			})

			It("reports that the common node has a degree of 2", func() {
				Expect(grph.Degree(node1)).To(Equal(2))
			})

			It("reports that the common node has two successors", func() {
				testSet(grph.Successors(node1), node2, node3)
			})

			It("reports the two unique nodes as adjacent to the common one", func() {
				testSet(grph.AdjacentNodes(node1), node2, node3)
			})

			It("has two edges sharing a common node", func() {
				testTwoEdges(
					grph,
					grph.Edges(),
					graph.EndpointPairOf(node1, node2),
					graph.EndpointPairOf(node1, node3))
			})

			It("has two incident edges connected to the common node", func() {
				testTwoEdges(
					grph,
					grph.IncidentEdges(node1),
					graph.EndpointPairOf(node1, node2),
					graph.EndpointPairOf(node1, node3))
			})

			It("reports that the common node has an out degree of 2", func() {
				Expect(grph.OutDegree(node1)).To(Equal(2))
			})
		})

		Context("when putting two connected edges with the same target node", func() {
			BeforeEach(func() {
				grph = putEdge(grph, node1, node2)
				grph = putEdge(grph, node3, node2)
			})

			It("reports that the common node has two predecessors", func() {
				testSet(grph.Predecessors(node2), node1, node3)
			})

			It("reports that the common node has an in degree of 2", func() {
				Expect(grph.InDegree(node2)).To(Equal(2))
			})

			It("has two incident edges connected to the common node", func() {
				testTwoEdges(
					grph,
					grph.IncidentEdges(node2),
					graph.EndpointPairOf(node1, node2),
					graph.EndpointPairOf(node3, node2))
			})
		})

		Context("when finding the predecessors of an absent node", func() {
			It("returns an empty set", func() {
				testSet(grph.Predecessors(nodeNotInGraph))
			})
		})

		Context("when finding the successors of an absent node", func() {
			It("returns an empty set", func() {
				testSet(grph.Successors(nodeNotInGraph))
			})
		})

		Context("when finding the adjacent nodes of an absent node", func() {
			It("returns an empty set", func() {
				testSet(grph.AdjacentNodes(nodeNotInGraph))
			})
		})

		Context("when finding the incident edges of an absent node", func() {
			It("returns an empty set", func() {
				testEmptyEdges(grph.IncidentEdges(nodeNotInGraph))
			})
		})

		Context("when finding the degree of an absent node", func() {
			It("returns zero", func() {
				Expect(grph.Degree(nodeNotInGraph)).To(BeZero())
			})
		})

		Context("when finding the in degree of an absent node", func() {
			It("returns zero", func() {
				Expect(grph.InDegree(nodeNotInGraph)).To(BeZero())
			})
		})

		Context("when finding the out degree of of an absent node", func() {
			It("returns zero", func() {
				Expect(grph.OutDegree(nodeNotInGraph)).To(BeZero())
			})
		})

		It("has an unmodifiable nodes set view", func() {
			nodes := grph.Nodes()
			Expect(nodes).To(BeNonMutableSet[int]())

			grph = addNode(grph, node1)
			testSet(nodes, node1)
		})

		It("has an unmodifiable adjacent nodes set view", func() {
			adjacentNodes := grph.AdjacentNodes(node1)
			Expect(adjacentNodes).To(BeNonMutableSet[int]())

			grph = putEdge(grph, node1, node2)
			grph = putEdge(grph, node3, node1)
			testSet(adjacentNodes, node2, node3)
		})

		It("has an unmodifiable predecessors set view", func() {
			predecessors := grph.Predecessors(node1)
			Expect(predecessors).To(BeNonMutableSet[int]())

			grph = putEdge(grph, node2, node1)
			testSet(predecessors, node2)
		})

		It("has an unmodifiable successors set view", func() {
			successors := grph.Successors(node1)
			Expect(successors).To(BeNonMutableSet[int]())

			grph = putEdge(grph, node1, node2)
			testSet(successors, node2)
		})

		It("has an unmodifiable set view of edges", func() {
			edges := grph.Edges()
			Expect(edges).To(BeNonMutableSet[graph.EndpointPair[int]]())

			grph = putEdge(grph, node1, node2)
			testSingleEdge(
				grph,
				grph.Edges(),
				graph.EndpointPairOf(node1, node2))
		})

		It("has an unmodifiable set view of incident edges", func() {
			incidentEdges := grph.IncidentEdges(node1)
			Expect(incidentEdges).To(BeNonMutableSet[graph.EndpointPair[int]]())

			grph = putEdge(grph, node1, node2)
			testSingleEdge(
				grph,
				grph.IncidentEdges(node1),
				graph.EndpointPairOf(node1, node2))
		})
	})

	if mutability == Mutable {
		mutableGraphTests(graphName, createGraph)
	}
	if mutability == Immutable {
		immutableGraphTests(graphName, createGraph)
	}

	if selfLoopsMode == AllowsSelfLoops {
		allowsSelfLoopsGraphTests(graphName, createGraph)

		if mutability == Mutable {
			allowsSelfLoopsMutableGraphTests(graphName, createGraph)
		}
	}
	if selfLoopsMode == DisallowsSelfLoops {
		disallowsSelfLoopsGraphTests(graphName, createGraph)
	}

	if directionMode == Undirected {
		undirectedGraphTests(graphName, createGraph)

		if selfLoopsMode == AllowsSelfLoops {
			undirectedAllowsSelfLoopGraphTests(graphName, createGraph)
		}
		if selfLoopsMode == DisallowsSelfLoops {
			undirectedDisallowsSelfLoopGraphTests(graphName, createGraph)
		}
	}

	if directionMode == Directed {
		directedGraphTests(graphName, createGraph)

		if selfLoopsMode == AllowsSelfLoops {
			directedAllowsSelfLoopGraphTests(graphName, createGraph)
		}
		if selfLoopsMode == DisallowsSelfLoops {
			directedDisallowsSelfLoopGraphTests(graphName, createGraph)
		}
	}
}

func mutableGraphTests(graphName string, createGraph func() graph.Graph[int]) {
	Context(fmt.Sprintf("%s: given a mutable graph", graphName), func() {
		var grph graph.MutableGraph[int]

		createGraphAsMutable := func() graph.MutableGraph[int] {
			g := createGraph()

			mutG, ok := g.(graph.MutableGraph[int])
			if !ok {
				panic("g is expected to implement graph.MutableGraph, but it doesn't")
			}

			return mutG
		}

		BeforeEach(func() {
			grph = createGraphAsMutable()
		})

		Context("when adding a new node", func() {
			It("returns true", func() {
				Expect(grph.AddNode(node1)).To(BeTrue())
			})
		})

		Context("when adding an existing node", func() {
			It("returns false", func() {
				grph.AddNode(node1)

				Expect(grph.AddNode(node1)).To(BeFalse())
			})
		})

		Context("when removing an existing node", func() {
			var removed bool

			BeforeEach(func() {
				grph.PutEdge(node1, node2)
				grph.PutEdge(node3, node1)
				grph.PutEdge(node2, node3)

				removed = grph.RemoveNode(node1)
			})

			It("returns true", func() {
				Expect(removed).To(BeTrue())
			})

			It("it leaves the other nodes alone", func() {
				testSet(grph.Nodes(), node2, node3)
			})

			It("removes its connections to its adjacent nodes", func() {
				testSet(grph.AdjacentNodes(node2), node3)
				testSet(grph.AdjacentNodes(node3), node2)
			})

			It("removes the connected edges", func() {
				testSingleEdge(
					grph,
					grph.Edges(),
					graph.EndpointPairOf(node2, node3))
			})
		})

		Context("when removing an absent node", func() {
			var removed bool

			BeforeEach(func() {
				grph.AddNode(node1)

				removed = grph.RemoveNode(nodeNotInGraph)
			})

			It("returns false", func() {
				Expect(removed).To(BeFalse())
			})

			It("leaves all the nodes alone", func() {
				testSet(grph.Nodes(), node1)
			})
		})

		Context("when putting a new edge", func() {
			It("returns true", func() {
				result := grph.PutEdge(node1, node2)

				Expect(result).To(BeTrue())
			})
		})

		Context("when putting an existing edge", func() {
			It("returns false", func() {
				grph.PutEdge(node1, node2)

				result := grph.PutEdge(node1, node2)

				Expect(result).To(BeFalse())
			})
		})

		Context("when putting two anti-parallel edges", func() {
			Context("and removing one of the nodes", func() {
				BeforeEach(func() {
					grph.PutEdge(node1, node2)
					grph.PutEdge(node2, node1)
					grph.RemoveNode(node1)
				})

				It("leaves the other node alone", func() {
					testSet(grph.Nodes(), node2)
				})

				It("removes both edges", func() {
					testEmptyEdges(grph.Edges())
				})
			})
		})

		Context("when removing an existing edge", func() {
			It("returns true", func() {
				grph.PutEdge(node1, node2)
				grph.PutEdge(node1, node3)

				removed := grph.RemoveEdge(node1, node2)

				Expect(removed).To(BeTrue())
			})

			It("removes the connection between its nodes", func() {
				grph.PutEdge(node1, node2)
				grph.PutEdge(node1, node3)

				grph.RemoveEdge(node1, node2)

				testSet(grph.Successors(node1), node3)
				testSet(grph.Predecessors(node3), node1)
				testSet(grph.Predecessors(node2))
			})

			It("removes the edge", func() {
				grph.PutEdge(node1, node2)

				grph.RemoveEdge(node1, node2)

				testEmptyEdges(grph.Edges())
			})
		})

		Context("when removing an absent edge with an existing source", func() {
			var removed bool

			BeforeEach(func() {
				grph.PutEdge(node1, node2)

				removed = grph.RemoveEdge(node1, nodeNotInGraph)
			})

			It("returns false", func() {
				Expect(removed).To(BeFalse())
			})

			It("leaves the existing nodes alone", func() {
				testSet(grph.Successors(node1), node2)
				testSet(grph.Predecessors(node2), node1)
			})
		})

		Context("when removing an absent edge with an existing target", func() {
			var removed bool

			BeforeEach(func() {
				grph.PutEdge(node1, node2)

				removed = grph.RemoveEdge(nodeNotInGraph, node2)
			})

			It("returns false", func() {
				Expect(removed).To(BeFalse())
			})

			It("leaves the existing nodes alone", func() {
				testSet(grph.Successors(node1), node2)
				testSet(grph.Predecessors(node2), node1)
			})
		})

		Context("when removing an absent edge with two existing nodes", func() {
			var removed bool

			BeforeEach(func() {
				grph.AddNode(node1)
				grph.AddNode(node2)

				removed = grph.RemoveEdge(node1, node2)
			})

			It("returns false", func() {
				Expect(removed).To(BeFalse())
			})

			It("leaves the existing nodes alone", func() {
				testSet(grph.Nodes(), node1, node2)
			})
		})
	})
}

func immutableGraphTests(graphName string, createGraph func() graph.Graph[int]) {
	Context(fmt.Sprintf("%s: given an immutable graph", graphName), func() {
		var _ graph.Graph[int]

		createGraphAsImmutable := func() graph.Graph[int] {
			g := createGraph()

			_, ok := g.(graph.MutableGraph[int])
			if ok {
				panic("g is expected to not implement graph.MutableGraph, but it does")
			}

			return g
		}

		BeforeEach(func() {
			_ = createGraphAsImmutable()
		})
	})
}

func undirectedGraphTests(
	graphName string,
	createGraph func() graph.Graph[int],
) {
	Context(fmt.Sprintf("%s: given an undirected graph", graphName), func() {
		var grph graph.Graph[int]

		BeforeEach(func() {
			grph = createGraph()
		})

		It("is undirected", func() {
			Expect(grph.IsDirected()).To(BeFalse())
		})

		Context("when putting one edge", func() {
			BeforeEach(func() {
				grph = putEdge(grph, node1, node2)
			})

			It("sees that the first node has the second node as a predecessor", func() {
				testSet(grph.Predecessors(node1), node2)
			})

			It("sees that the second node has the first node as a successor", func() {
				testSet(grph.Successors(node2), node1)
			})

			It("has an in degree of 1 for the first node", func() {
				Expect(grph.InDegree(node1)).To(Equal(1))
			})

			It("has an out degree of 1 for the second node", func() {
				Expect(grph.OutDegree(node2)).To(Equal(1))
			})

			It("sees the second node as being connected to the first", func() {
				Expect(grph.HasEdgeConnecting(node2, node1)).
					To(BeTrue())

				Expect(grph.HasEdgeConnectingEndpoints(graph.EndpointPairOf(node2, node1))).
					To(BeTrue())
			})
		})
	})
}

func directedGraphTests(
	graphName string,
	createGraph func() graph.Graph[int],
) {
	Context(fmt.Sprintf("%s: given a directed graph", graphName), func() {
		var grph graph.Graph[int]

		BeforeEach(func() {
			grph = createGraph()
		})

		It("is directed", func() {
			Expect(grph.IsDirected()).To(BeTrue())
		})

		Context("when putting one edge", func() {
			BeforeEach(func() {
				grph = putEdge(grph, node1, node2)
			})

			It("sees that the first node has no predecessors", func() {
				testSet(grph.Predecessors(node1))
			})

			It("sees that the second node has no successors", func() {
				testSet(grph.Successors(node2))
			})

			It("has an in degree of 0 for the first node", func() {
				Expect(grph.InDegree(node1)).To(BeZero())
			})

			It("has an out degree of 0 for the second node", func() {
				Expect(grph.OutDegree(node2)).To(BeZero())
			})

			It("does not see the second node as being connected to the first node", func() {
				Expect(grph.HasEdgeConnecting(node2, node1)).
					To(BeFalse())

				Expect(grph.HasEdgeConnectingEndpoints(graph.EndpointPairOf(node2, node1))).
					To(BeFalse())
			})
		})

		Context("when putting two connected edges that form a line graph", func() {
			It("reports that the common node has a degree of 2", func() {
				grph = putEdge(grph, node1, node2)
				grph = putEdge(grph, node2, node3)

				Expect(grph.Degree(node2)).To(Equal(2))
			})
		})
	})
}

func allowsSelfLoopsGraphTests(graphName string, createGraph func() graph.Graph[int]) {
	Context(fmt.Sprintf("%s: given a graph that allows self loops", graphName), func() {
		var grph graph.Graph[int]

		BeforeEach(func() {
			grph = createGraph()
		})

		It("allows self loops", func() {
			Expect(grph.AllowsSelfLoops()).To(BeTrue())
		})

		Context("when putting one self-loop edge", func() {
			BeforeEach(func() {
				grph = putEdge(grph, node1, node1)
			})

			It("sees the shared node as its own adjacent node", func() {
				testSet(grph.AdjacentNodes(node1), node1)
			})

			It(
				"reports that the shared node has a degree of 2 because the edge touches the node twice",
				func() {
					Expect(grph.Degree(node1)).To(Equal(2))
				})
		})
	})
}

func allowsSelfLoopsMutableGraphTests(graphName string, createGraph func() graph.Graph[int]) {
	Context(fmt.Sprintf("%s: given a mutable graph that allows self loops", graphName), func() {
		var grph graph.MutableGraph[int]

		createGraphAsMutable := func() graph.MutableGraph[int] {
			g := createGraph()

			mutG, ok := g.(graph.MutableGraph[int])
			if !ok {
				panic("g is expected to implement graph.MutableGraph, but it doesn't")
			}

			return mutG
		}

		BeforeEach(func() {
			grph = createGraphAsMutable()
		})

		Context("when removing a self-looping node", func() {
			It("reports that the self-loop edge is gone", func() {
				grph.PutEdge(node1, node1)

				grph.RemoveNode(node1)

				testEmptyEdges(grph.Edges())
			})
		})
	})
}

func disallowsSelfLoopsGraphTests(graphName string, createGraph func() graph.Graph[int]) {
	Context(fmt.Sprintf("%s: given a graph that disallows self loops", graphName), func() {
		var grph graph.Graph[int]

		BeforeEach(func() {
			grph = createGraph()
		})

		It("disallows self loops", func() {
			Expect(grph.AllowsSelfLoops()).To(BeFalse())
		})

		Context("when putting one self-loop edge", func() {
			It("panics", func() {
				Expect(func() { grph = putEdge(grph, node1, node1) }).
					To(PanicWith("self-loops are disallowed"))
			})
		})
	})
}

func undirectedAllowsSelfLoopGraphTests(graphName string, createGraph func() graph.Graph[int]) {
	Context(fmt.Sprintf("%s: given an undirected graph that allows self loops", graphName), func() {
		var grph graph.Graph[int]

		BeforeEach(func() {
			grph = createGraph()
		})

		It("has an appropriate string representation", func() {
			Expect(grph).To(
				HaveStringRepr(
					"isDirected: false, allowsSelfLoops: true, nodes: [], edges: []"))
		})

		Context("when adding one node", func() {
			It("has an appropriate string representation", func() {
				grph = addNode(grph, node1)

				Expect(grph).To(
					HaveStringRepr(
						"isDirected: false, allowsSelfLoops: true, nodes: [1], edges: []"))
			})
		})

		Context("when putting one edge", func() {
			It("has an appropriate string representation", func() {
				grph = putEdge(grph, node1, node2)

				Expect(grph).To(
					HaveStringReprThatIsAnyOf(
						"isDirected: false, allowsSelfLoops: true, nodes: [1, 2], edges: [<1 -> 2>]",
						"isDirected: false, allowsSelfLoops: true, nodes: [2, 1], edges: [<1 -> 2>]",
						"isDirected: false, allowsSelfLoops: true, nodes: [1, 2], edges: [<2 -> 1>]",
						"isDirected: false, allowsSelfLoops: true, nodes: [2, 1], edges: [<2 -> 1>]"))
			})
		})
	})
}

func undirectedDisallowsSelfLoopGraphTests(graphName string, createGraph func() graph.Graph[int]) {
	Context(fmt.Sprintf("%s: given an undirected graph that disallows self loops", graphName), func() {
		var grph graph.Graph[int]

		BeforeEach(func() {
			grph = createGraph()
		})

		It("has an appropriate string representation", func() {
			Expect(grph).To(
				HaveStringRepr(
					"isDirected: false, allowsSelfLoops: false, nodes: [], edges: []"))
		})

		Context("when adding one node", func() {
			It("has an appropriate string representation", func() {
				grph = addNode(grph, node1)

				Expect(grph).To(
					HaveStringRepr(
						"isDirected: false, allowsSelfLoops: false, nodes: [1], edges: []"))
			})
		})

		Context("when putting one edge", func() {
			It("has an appropriate string representation", func() {
				grph = putEdge(grph, node1, node2)

				Expect(grph).To(
					HaveStringReprThatIsAnyOf(
						"isDirected: false, allowsSelfLoops: false, nodes: [1, 2], edges: [<1 -> 2>]",
						"isDirected: false, allowsSelfLoops: false, nodes: [2, 1], edges: [<1 -> 2>]",
						"isDirected: false, allowsSelfLoops: false, nodes: [1, 2], edges: [<2 -> 1>]",
						"isDirected: false, allowsSelfLoops: false, nodes: [2, 1], edges: [<2 -> 1>]"))
			})
		})
	})
}

func directedAllowsSelfLoopGraphTests(graphName string, createGraph func() graph.Graph[int]) {
	Context(fmt.Sprintf("%s: given a directed graph that allows self loops", graphName), func() {
		var grph graph.Graph[int]

		BeforeEach(func() {
			grph = createGraph()
		})

		It("has an appropriate string representation", func() {
			Expect(grph).To(
				HaveStringRepr(
					"isDirected: true, allowsSelfLoops: true, nodes: [], edges: []"))
		})

		Context("when adding one node", func() {
			It("has an appropriate string representation", func() {
				grph = addNode(grph, node1)

				Expect(grph).To(
					HaveStringRepr(
						"isDirected: true, allowsSelfLoops: true, nodes: [1], edges: []"))
			})
		})

		Context("when putting one edge", func() {
			It("has an appropriate string representation", func() {
				grph = putEdge(grph, node1, node2)

				Expect(grph).To(
					HaveStringReprThatIsAnyOf(
						"isDirected: true, allowsSelfLoops: true, nodes: [1, 2], edges: [<1 -> 2>]",
						"isDirected: true, allowsSelfLoops: true, nodes: [2, 1], edges: [<1 -> 2>]"))
			})
		})

		Context("when putting one self-loop edge", func() {
			It("reports that the shared node only has one incident edge with itself", func() {
				grph = putEdge(grph, node1, node1)

				testSingleEdge(
					grph,
					grph.IncidentEdges(node1),
					graph.EndpointPairOf(node1, node1))
			})
		})
	})
}

func directedDisallowsSelfLoopGraphTests(graphName string, createGraph func() graph.Graph[int]) {
	Context(fmt.Sprintf("%s: given a directed graph that disallows self loops", graphName), func() {
		var grph graph.Graph[int]

		BeforeEach(func() {
			grph = createGraph()
		})

		It("has an appropriate string representation", func() {
			Expect(grph).To(
				HaveStringRepr(
					"isDirected: true, allowsSelfLoops: false, nodes: [], edges: []"))
		})

		Context("when adding one node", func() {
			It("has an appropriate string representation", func() {
				grph = addNode(grph, node1)

				Expect(grph).To(
					HaveStringRepr(
						"isDirected: true, allowsSelfLoops: false, nodes: [1], edges: []"))
			})
		})

		Context("when putting one edge", func() {
			It("has an appropriate string representation", func() {
				grph = putEdge(grph, node1, node2)

				Expect(grph).To(
					HaveStringReprThatIsAnyOf(
						"isDirected: true, allowsSelfLoops: false, nodes: [1, 2], edges: [<1 -> 2>]",
						"isDirected: true, allowsSelfLoops: false, nodes: [2, 1], edges: [<1 -> 2>]"))
			})
		})
	})
}

func testSet(s set.Set[int], expectedValues ...int) {
	// Set.Len()
	Expect(s).To(HaveLenOf(len(expectedValues)))

	// Set.All()
	if len(expectedValues) == 0 {
		Expect(s).To(HaveAllThatEmitsNothing[int]())
	} else {
		Expect(s).To(HaveAllThatConsistsOfElementsInSlice(expectedValues))
	}

	// Set.Contains()
	for _, value := range []int{node1, node2, node3} {
		if slices.Contains(expectedValues, value) {
			Expect(s).To(Contain(value))
		} else {
			Expect(s).ToNot(Contain(value))
		}
	}
	Expect(s).ToNot(Contain(nodeNotInGraph))

	// Set.String()
	str := s.String()
	Expect(str).To(HavePrefix("["))
	Expect(str).To(HaveSuffix("]"))

	expectedValueStrs := make([]string, 0, len(expectedValues))
	for _, v := range expectedValues {
		expectedValueStrs = append(expectedValueStrs, strconv.Itoa(v))
	}

	trimmed := strings.Trim(str, "[]")
	actualValueStrs := strings.SplitN(trimmed, ", ", len(expectedValues))

	Expect(actualValueStrs).To(
		ConsistOf(expectedValueStrs),
		"to find all elements in string repr of set")
}

func testEmptyEdges(edges set.Set[graph.EndpointPair[int]]) {
	// Set.Len()
	Expect(edges).To(HaveLenOfZero())

	// Set.All()
	Expect(edges).To(
		HaveAllThatEmitsNothing[graph.EndpointPair[int]]())

	// Set.Contains()
	Expect(edges).ToNot(
		Contain(
			graph.EndpointPairOf(
				nodeNotInGraph, nodeNotInGraph)))

	// Set.String()
	Expect(edges).To(HaveStringRepr("[]"))
}

func testSingleEdge(
	g graph.Graph[int],
	endpointPairs set.Set[graph.EndpointPair[int]],
	expectedEndpointPair graph.EndpointPair[int],
) {
	// Set.Len()
	Expect(endpointPairs).To(HaveLenOf(1))

	// Set.All()
	if g.IsDirected() {
		Expect(slices.Collect(endpointPairs.All())).
			To(HaveExactElements(expectedEndpointPair))
	} else {
		Expect(slices.Collect(endpointPairs.All())).
			To(Or(
				HaveExactElements(expectedEndpointPair),
				HaveExactElements(reverseOf(expectedEndpointPair))))
	}

	// Set.Contains()
	Expect(endpointPairs).To(Contain(expectedEndpointPair))
	if expectedEndpointPair != reverseOf(expectedEndpointPair) {
		if g.IsDirected() {
			Expect(endpointPairs).ToNot(Contain(reverseOf(expectedEndpointPair)))
		} else {
			Expect(endpointPairs).To(Contain(reverseOf(expectedEndpointPair)))
		}
	}
	Expect(endpointPairs).ToNot(
		Contain(graph.EndpointPairOf(nodeNotInGraph, nodeNotInGraph)))

	// Set.String()
	if g.IsDirected() {
		Expect(endpointPairs).To(
			HaveStringRepr(fmt.Sprintf("[%v]", expectedEndpointPair)))
	} else {
		Expect(endpointPairs).To(
			HaveStringReprThatIsAnyOf(
				fmt.Sprintf("[%v]", expectedEndpointPair),
				fmt.Sprintf("[%v]", reverseOf(expectedEndpointPair))))
	}
}

func testTwoEdges(
	g graph.Graph[int],
	endpointPairs set.Set[graph.EndpointPair[int]],
	firstExpectedPair graph.EndpointPair[int],
	secondExpectedPair graph.EndpointPair[int],
) {
	// Set.Len()
	Expect(endpointPairs).To(HaveLenOf(2))

	// Set.All()
	if g.IsDirected() {
		Expect(slices.Collect(endpointPairs.All())).
			To(ConsistOf(firstExpectedPair, secondExpectedPair))
	} else {
		Expect(slices.Collect(endpointPairs.All())).
			To(
				ConsistOf(
					Or(
						Equal(firstExpectedPair),
						Equal(reverseOf(firstExpectedPair)),
					),
					Or(
						Equal(secondExpectedPair),
						Equal(reverseOf(secondExpectedPair)),
					),
				))
	}

	// Set.Contains()
	Expect(endpointPairs).To(Contain(firstExpectedPair))
	Expect(endpointPairs).To(Contain(secondExpectedPair))
	if g.IsDirected() {
		Expect(endpointPairs).ToNot(Contain(reverseOf(firstExpectedPair)))
		Expect(endpointPairs).ToNot(Contain(reverseOf(secondExpectedPair)))
	} else {
		Expect(endpointPairs).To(Contain(reverseOf(firstExpectedPair)))
		Expect(endpointPairs).To(Contain(reverseOf(secondExpectedPair)))
	}
	Expect(endpointPairs).ToNot(
		Contain(graph.EndpointPairOf(nodeNotInGraph, nodeNotInGraph)))

	// Set.String()
	if g.IsDirected() {
		Expect(endpointPairs).To(
			HaveStringReprThatIsAnyOf(
				fmt.Sprintf("[%v, %v]", firstExpectedPair, secondExpectedPair),
				fmt.Sprintf("[%v, %v]", secondExpectedPair, firstExpectedPair)))
	} else {
		Expect(endpointPairs).To(
			HaveStringReprThatIsAnyOf(
				fmt.Sprintf("[%v, %v]", firstExpectedPair, secondExpectedPair),
				fmt.Sprintf("[%v, %v]", firstExpectedPair, reverseOf(secondExpectedPair)),
				fmt.Sprintf("[%v, %v]", reverseOf(firstExpectedPair), secondExpectedPair),
				fmt.Sprintf("[%v, %v]", reverseOf(firstExpectedPair), reverseOf(secondExpectedPair)),
				fmt.Sprintf("[%v, %v]", secondExpectedPair, firstExpectedPair),
				fmt.Sprintf("[%v, %v]", secondExpectedPair, reverseOf(firstExpectedPair)),
				fmt.Sprintf("[%v, %v]", reverseOf(secondExpectedPair), firstExpectedPair),
				fmt.Sprintf("[%v, %v]", reverseOf(secondExpectedPair), reverseOf(firstExpectedPair))))
	}
}

func reverseOf(endpointPair graph.EndpointPair[int]) graph.EndpointPair[int] {
	return graph.EndpointPairOf(endpointPair.Target(), endpointPair.Source())
}
