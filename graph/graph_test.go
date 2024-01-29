package graph_test

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jbduncan/go-containers/graph"
	. "github.com/jbduncan/go-containers/internal/matchers"
	"github.com/jbduncan/go-containers/set"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/exp/slices"
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
		"graph.Undirected[int]().AllowsSelfLoops().Build()",
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

		FIt("has no nodes", func() {
			testSet(grph.Nodes())
		})

		FIt("has no edges", func() {
			testEmptyEdges(grph.Edges())
		})

		Context("when adding one node", func() {
			BeforeEach(func() {
				grph = addNode(grph, node1)
			})

			FIt("has just that node", func() {
				testSet(grph.Nodes(), node1)
			})

			FIt("reports that the node has no adjacent nodes", func() {
				testSet(grph.AdjacentNodes(node1))
			})

			FIt("reports that the node has no predecessors", func() {
				testSet(grph.Predecessors(node1))
			})

			FIt("reports that the node has no successors", func() {
				testSet(grph.Successors(node1))
			})

			FIt("reports that the node has no incident edges", func() {
				testEmptyEdges(grph.IncidentEdges(node1))
			})

			FIt("reports that the node has a degree of 0", func() {
				Expect(grph.Degree(node1)).To(BeZero())
			})

			FIt("reports that the node has an in degree of 0", func() {
				Expect(grph.InDegree(node1)).To(BeZero())
			})

			FIt("reports that the node has an out degree of 0", func() {
				Expect(grph.OutDegree(node1)).To(BeZero())
			})
		})

		Context("when adding two nodes", func() {
			BeforeEach(func() {
				grph = addNode(grph, node1)
				grph = addNode(grph, node2)
			})

			FIt("has both nodes", func() {
				testSet(grph.Nodes(), node1, node2)
			})
		})

		Context("when putting one edge", func() {
			FIt("reports that both nodes are adjacent to each other", func() {
				grph = putEdge(grph, node1, node2)

				testSet(grph.AdjacentNodes(node1), node2)
				testSet(grph.AdjacentNodes(node2), node1)
			})

			It("reports that both nodes have a degree of 1", func() {
				grph = putEdge(grph, node1, node2)

				Expect(grph.Degree(node1)).To(Equal(1))
				Expect(grph.Degree(node2)).To(Equal(1))
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

		Context("when putting two connected edges with the same source node", func() {
			BeforeEach(func() {
				grph = putEdge(grph, node1, node2)
				grph = putEdge(grph, node1, node3)
			})

			It("reports that the common node has a degree of 2", func() {
				Expect(grph.Degree(node1)).To(Equal(2))
			})

			FIt("reports that the common node has two successors", func() {
				testSet(grph.Successors(node1), node2, node3)
			})

			FIt("reports the two unique nodes as adjacent to the common one", func() {
				testSet(grph.AdjacentNodes(node1), node2, node3)
			})
		})

		Context("when putting two connected edges with the same target node", func() {
			FIt("reports that the common node has two predecessors", func() {
				grph = putEdge(grph, node1, node2)
				grph = putEdge(grph, node3, node2)

				testSet(grph.Predecessors(node2), node1, node3)
			})
		})

		Context("when finding the predecessors of an absent node", func() {
			FIt("returns an empty set", func() {
				testSet(grph.Predecessors(nodeNotInGraph))
			})
		})

		Context("when finding the successors of an absent node", func() {
			FIt("returns an empty set", func() {
				testSet(grph.Successors(nodeNotInGraph))
			})
		})

		Context("when finding the adjacent nodes of an absent node", func() {
			FIt("returns an empty set", func() {
				testSet(grph.AdjacentNodes(nodeNotInGraph))
			})
		})

		Context("when finding the incident edges of an absent node", func() {
			FIt("returns an empty set", func() {
				testEmptyEdges(grph.IncidentEdges(nodeNotInGraph))
			})
		})

		Context("when finding the degree of an absent node", func() {
			FIt("returns zero", func() {
				Expect(grph.Degree(nodeNotInGraph)).To(BeZero())
			})
		})

		Context("when finding the in degree of an absent node", func() {
			FIt("returns zero", func() {
				Expect(grph.InDegree(nodeNotInGraph)).To(BeZero())
			})
		})

		Context("when finding the out degree of of an absent node", func() {
			FIt("returns zero", func() {
				Expect(grph.OutDegree(nodeNotInGraph)).To(BeZero())
			})
		})

		FIt("has an unmodifiable nodes set view", func() {
			nodes := grph.Nodes()
			Expect(nodes).To(BeNonMutableSet[int]())

			grph = addNode(grph, node1)
			testSet(nodes, node1)
		})

		FIt("has an unmodifiable adjacent nodes set view", func() {
			adjacentNodes := grph.AdjacentNodes(node1)
			Expect(adjacentNodes).To(BeNonMutableSet[int]())

			grph = putEdge(grph, node1, node2)
			grph = putEdge(grph, node3, node1)
			testSet(adjacentNodes, node2, node3)
		})

		FIt("had an unmodifiable predecessors set view", func() {
			predecessors := grph.Predecessors(node1)
			Expect(predecessors).To(BeNonMutableSet[int]())

			grph = putEdge(grph, node2, node1)
			testSet(predecessors, node2)
		})

		FIt("has an unmodifiable successors set view", func() {
			successors := grph.Successors(node1)
			Expect(successors).To(BeNonMutableSet[int]())

			grph = putEdge(grph, node1, node2)
			testSet(successors, node2)
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

				removed = grph.RemoveNode(node1)
			})

			It("returns true", func() {
				Expect(removed).To(BeTrue())
			})

			It("it leaves the other nodes alone", func() {
				testSet(grph.Nodes(), node2, node3)
			})

			It("removes its connections to its adjacent nodes", func() {
				testSet(grph.AdjacentNodes(node2))
				testSet(grph.AdjacentNodes(node3))
			})

			It("removes the connected edges", func() {
				testEmptyEdges(grph.Edges())
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
			var removed bool

			BeforeEach(func() {
				grph.PutEdge(node1, node2)
				grph.PutEdge(node1, node3)

				removed = grph.RemoveEdge(node1, node2)
			})

			It("returns true", func() {
				Expect(removed).To(BeTrue())
			})

			It("removes the connection between its nodes", func() {
				testSet(grph.Successors(node1), node3)
				testSet(grph.Predecessors(node3), node1)
				testSet(grph.Predecessors(node2))
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
	FContext(fmt.Sprintf("%s: given an undirected graph", graphName), func() {
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

			It("has just one edge", func() {
				testSingleEdgeForUndirectedGraph(grph.Edges())
			})

			It("has an incident edge connecting the first node to the second node", func() {
				testSingleEdgeForUndirectedGraph(grph.IncidentEdges(node1))
			})

			It("sees that the second node has the first node as a predecessor", func() {
				testSet(grph.Predecessors(node2), node1)
			})

			It("sees that the first node has the second node as a predecessor", func() {
				testSet(grph.Predecessors(node1), node2)
			})

			It("sees that the second node has the first node as a successor", func() {
				testSet(grph.Successors(node2), node1)
			})

			It("sees that the first node has the second node as a successor", func() {
				testSet(grph.Successors(node1), node2)
			})

			It("has an incident edge connecting the second node to the first node", func() {
				testSingleEdgeForUndirectedGraph(grph.IncidentEdges(node2))
			})

			It("has an in degree of 1 for the first node", func() {
				Expect(grph.InDegree(node1)).To(Equal(1))
			})

			It("has an in degree of 1 for the second node", func() {
				Expect(grph.InDegree(node2)).To(Equal(1))
			})

			It("has an out degree of 1 for the first node", func() {
				Expect(grph.OutDegree(node1)).To(Equal(1))
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

		Context("when putting two connected edges with the same source node", func() {
			BeforeEach(func() {
				grph = putEdge(grph, node1, node2)
				grph = putEdge(grph, node1, node3)
			})

			It("has two edges sharing a common node", func() {
				testTwoEdgesForUndirectedGraphs(grph.Edges())
			})

			It("has two incident edges connected to the common node", func() {
				testTwoEdgesForUndirectedGraphs(grph.IncidentEdges(node1))
			})
		})

		It("has an unmodifiable set view of edges", func() {
			edges := grph.Edges()
			Expect(edges).To(BeNonMutableSet[graph.EndpointPair[int]]())

			grph = putEdge(grph, node1, node2)
			testSingleEdgeForUndirectedGraph(edges)
		})

		It("has an unmodifiable set view of incident edges", func() {
			incidentEdges := grph.IncidentEdges(node1)
			Expect(incidentEdges).To(BeNonMutableSet[graph.EndpointPair[int]]())

			grph = putEdge(grph, node1, node2)
			testSingleEdgeForUndirectedGraph(incidentEdges)
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

		FIt("is directed", func() {
			Expect(grph.IsDirected()).To(BeTrue())
		})

		Context("when putting one edge", func() {
			BeforeEach(func() {
				grph = putEdge(grph, node1, node2)
			})

			FIt("sees that the second node has the first node as a predecessor", func() {
				testSet(grph.Predecessors(node2), node1)
			})

			FIt("sees that the first node has no predecessors", func() {
				testSet(grph.Predecessors(node1))
			})

			FIt("sees that the first node has the second node as a successor", func() {
				testSet(grph.Successors(node1), node2)
			})

			FIt("sees that the second node has no successors", func() {
				testSet(grph.Successors(node2))
			})

			It("has just one edge", func() {
				testSingleEdgeForDirectedGraph(grph.Edges())
			})

			It("has an incident edge connecting the first node to the second node", func() {
				testSingleEdgeForDirectedGraph(grph.IncidentEdges(node1))
			})

			FIt("has an in degree of 0 for the first node", func() {
				Expect(grph.InDegree(node1)).To(BeZero())
			})

			FIt("has an in degree of 1 for the second node", func() {
				Expect(grph.InDegree(node2)).To(Equal(1))
			})

			FIt("has an out degree of 1 for the first node", func() {
				Expect(grph.OutDegree(node1)).To(Equal(1))
			})

			FIt("has an out degree of 0 for the second node", func() {
				Expect(grph.OutDegree(node2)).To(BeZero())
			})
		})

		Context("when putting two connected edges with the same source node", func() {
			BeforeEach(func() {
				grph = putEdge(grph, node1, node2)
				grph = putEdge(grph, node1, node3)
			})

			It("has two edges sharing a common node", func() {
				testTwoEdgesForDirectedGraphs(grph.Edges())
			})

			It("has two incident edges connected to the common node", func() {
				testTwoEdgesForDirectedGraphs(grph.IncidentEdges(node1))
			})

			FIt("has an out degree of 2 for the common node", func() {
				Expect(grph.OutDegree(node1)).To(Equal(2))
			})
		})

		Context("when putting two connected edges with the same target node", func() {
			FIt("has an in degree of 2 for the common node", func() {
				grph = putEdge(grph, node1, node2)
				grph = putEdge(grph, node3, node2)

				Expect(grph.InDegree(node2)).To(Equal(2))
			})
		})

		It("has an unmodifiable set view of edges", func() {
			edges := grph.Edges()
			Expect(edges).To(BeNonMutableSet[graph.EndpointPair[int]]())

			grph = putEdge(grph, node1, node2)
			testSingleEdgeForDirectedGraph(edges)
		})

		It("has an unmodifiable set view of incident edges", func() {
			incidentEdges := grph.IncidentEdges(node1)
			Expect(incidentEdges).To(BeNonMutableSet[graph.EndpointPair[int]]())

			grph = putEdge(grph, node1, node2)
			testSingleEdgeForDirectedGraph(incidentEdges)
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
			FIt("sees the shared node as its own adjacent node", func() {
				grph = putEdge(grph, node1, node1)

				testSet(grph.AdjacentNodes(node1), node1)
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

		FIt("disallows self loops", func() {
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
		It("has an appropriate string representation", func() {
			grph := createGraph()

			Expect(grph).To(
				HaveStringRepr(
					"isDirected: true, allowsSelfLoops: true, nodes: [], edges: []"))
		})
	})
}

func directedDisallowsSelfLoopGraphTests(graphName string, createGraph func() graph.Graph[int]) {
	Context(fmt.Sprintf("%s: given a directed graph that disallows self loops", graphName), func() {
		It("has an appropriate string representation", func() {
			grph := createGraph()

			Expect(grph).To(
				HaveStringRepr(
					"isDirected: true, allowsSelfLoops: false, nodes: [], edges: []"))
		})
	})
}

func testSet(s set.Set[int], expectedValues ...int) {
	// Set.Len()
	Expect(s).To(HaveLenOf(len(expectedValues)))

	// Set.ForEach()
	if len(expectedValues) == 0 {
		Expect(s).To(HaveForEachThatEmitsNothing[int]())
	} else {
		Expect(s).To(HaveForEachThatConsistsOfElementsInSlice(expectedValues))
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

	// Set.ForEach()
	Expect(edges).To(
		HaveForEachThatEmitsNothing[graph.EndpointPair[int]]())

	// Set.Contains()
	Expect(edges).ToNot(
		Contain(
			graph.EndpointPairOf(
				nodeNotInGraph, nodeNotInGraph)))

	// Set.String()
	Expect(edges).To(HaveStringRepr("[]"))
}

func testSingleEdgeForUndirectedGraph(endpointPairs set.Set[graph.EndpointPair[int]]) {
	// Set.Len()
	Expect(endpointPairs).To(HaveLenOf(1))

	// Set.ForEach()
	Expect(ForEachToSlice(endpointPairs)).
		To(Or(
			HaveExactElements(graph.EndpointPairOf(node1, node2)),
			HaveExactElements(graph.EndpointPairOf(node2, node1))))

	// Set.Contains()
	Expect(endpointPairs).To(Contain(graph.EndpointPairOf(node1, node2)))
	Expect(endpointPairs).To(Contain(graph.EndpointPairOf(node2, node1)))
	Expect(endpointPairs).ToNot(
		Contain(graph.EndpointPairOf(nodeNotInGraph, nodeNotInGraph)))

	// Set.String()
	Expect(endpointPairs).To(
		HaveStringReprThatIsAnyOf("[<1 -> 2>]", "[<2 -> 1>]"))
}

func testSingleEdgeForDirectedGraph(endpointPairs set.Set[graph.EndpointPair[int]]) {
	// Set.Len()
	Expect(endpointPairs).To(HaveLenOf(1))

	// Set.ForEach()
	Expect(ForEachToSlice(endpointPairs)).
		To(HaveExactElements(graph.EndpointPairOf(node1, node2)))

	// Set.Contains()
	Expect(endpointPairs).To(Contain(graph.EndpointPairOf(node1, node2)))
	Expect(endpointPairs).ToNot(Contain(graph.EndpointPairOf(node2, node1)))
	Expect(endpointPairs).ToNot(
		Contain(graph.EndpointPairOf(nodeNotInGraph, nodeNotInGraph)))

	// Set.String()
	Expect(endpointPairs).To(HaveStringRepr("[<1 -> 2>]"))
}

func testTwoEdgesForUndirectedGraphs(endpointPairs set.Set[graph.EndpointPair[int]]) {
	// Set.Len()
	Expect(endpointPairs).To(HaveLenOf(2))

	// Set.ForEach()
	Expect(ForEachToSlice(endpointPairs)).
		To(
			ConsistOf(
				Or(
					Equal(graph.EndpointPairOf(node1, node2)),
					Equal(graph.EndpointPairOf(node2, node1)),
				),
				Or(
					Equal(graph.EndpointPairOf(node1, node3)),
					Equal(graph.EndpointPairOf(node3, node1)),
				),
			))

	// Set.Contains()
	Expect(endpointPairs).To(Contain(graph.EndpointPairOf(node1, node2)))
	Expect(endpointPairs).To(Contain(graph.EndpointPairOf(node1, node3)))
	Expect(endpointPairs).To(Contain(graph.EndpointPairOf(node2, node1)))
	Expect(endpointPairs).To(Contain(graph.EndpointPairOf(node3, node1)))
	Expect(endpointPairs).ToNot(
		Contain(graph.EndpointPairOf(nodeNotInGraph, nodeNotInGraph)))

	// Set.String()
	Expect(endpointPairs).To(
		HaveStringReprThatIsAnyOf(
			"[<1 -> 2>, <1 -> 3>]",
			"[<1 -> 2>, <3 -> 1>]",
			"[<2 -> 1>, <1 -> 3>]",
			"[<2 -> 1>, <3 -> 1>]",
			"[<1 -> 3>, <1 -> 2>]",
			"[<1 -> 3>, <2 -> 1>]",
			"[<3 -> 1>, <1 -> 2>]",
			"[<3 -> 1>, <2 -> 1>]"))
}

func testTwoEdgesForDirectedGraphs(endpointPairs set.Set[graph.EndpointPair[int]]) {
	// Set.Len()
	Expect(endpointPairs).To(HaveLenOf(2))

	// Set.ForEach()
	Expect(ForEachToSlice(endpointPairs)).
		To(
			ConsistOf(
				Equal(graph.EndpointPairOf(node1, node2)),
				Equal(graph.EndpointPairOf(node1, node3)),
			))

	// Set.Contains()
	Expect(endpointPairs).To(Contain(graph.EndpointPairOf(node1, node2)))
	Expect(endpointPairs).To(Contain(graph.EndpointPairOf(node1, node3)))
	Expect(endpointPairs).ToNot(Contain(graph.EndpointPairOf(node2, node1)))
	Expect(endpointPairs).ToNot(Contain(graph.EndpointPairOf(node3, node1)))
	Expect(endpointPairs).ToNot(
		Contain(graph.EndpointPairOf(nodeNotInGraph, nodeNotInGraph)))

	// Set.String()
	Expect(endpointPairs).To(
		HaveStringReprThatIsAnyOf(
			"[<1 -> 2>, <1 -> 3>]",
			"[<1 -> 3>, <1 -> 2>]"))
}
