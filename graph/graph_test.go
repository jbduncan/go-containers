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

// TODO: Experiment with migrating to go test. Does it make the tests easier to read?

// TODO: Move mutableGraphTests to a graphtest package
// TODO: Migrate mutableGraphTests to a struct with a constructor function like:
//       graphtest.MutableGraph(
//           name string,
//           newGraph func() graph.MutableGraph[int],
//           ...graphtest.Option options)
//       ...where options is any of:
//         - WhereContainersAreViews
//         - WhereContainersAreCopies
//       ...and the default is `WhereContainersAreViews`

var _ = Describe("Graphs", func() {
	graphTests(
		"graph.Undirected[int]().Build()",
		func() graph.Graph[int] {
			return graph.Undirected[int]().Build()
		},
		ContainersAreViews)
	graphTests(
		"graph.Undirected[int]().AllowsSelfLoops().Build()",
		func() graph.Graph[int] {
			return graph.Undirected[int]().AllowsSelfLoops(true).Build()
		},
		ContainersAreViews)
})

const (
	node1          = 1
	node2          = 2
	node3          = 3
	nodeNotInGraph = 1_000
)

type ContainersMode int

const (
	ContainersAreViews ContainersMode = iota
	// TODO: Add additional test cases for graphs whose accessor containers (Nodes(), Edges(), etc.)
	//       are immutable copies.
	ContainersAreCopies
)

func addNode(grph graph.Graph[int], node int) graph.Graph[int] {
	// TODO: When introducing ImmutableGraph, expand addNode to recognise when grph is not
	//  mutable and return a copy of the graph with the added node instead.

	if graphAsMutable, ok := grph.(graph.MutableGraph[int]); ok {
		graphAsMutable.AddNode(node)
	}

	return grph
}

func putEdge(grph graph.Graph[int], node1 int, node2 int) graph.Graph[int] {
	// TODO: When introducing ImmutableGraph, expand putEdge to recognise when grph is not
	//  mutable and return a copy of the graph with the added edge instead.

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
//
// TODO: Move to a public package for graph testing utilities
func graphTests(
	graphName string,
	createGraph func() graph.Graph[int],
	containersMode ContainersMode,
) {
	Context(fmt.Sprintf("%s: given a graph", graphName), func() {
		var grph graph.Graph[int]

		BeforeEach(func() {
			assertContainersMode(containersMode)

			grph = createGraph()
		})

		It("has no nodes", func() {
			testSet(grph.Nodes())
		})

		It("has no edges", func() {
			testEmptyEdges(grph.Edges())
		})

		It("has an unmodifiable nodes set view", func() {
			if containersMode != ContainersAreViews {
				Skip("Graph.Nodes() is not expected to return an unmodifiable view")
			}

			nodes := grph.Nodes()
			Expect(nodes).To(BeNonMutableSet[int]())

			grph = addNode(grph, node1)
			testSet(nodes, node1)
		})

		// TODO: Write an equivalent test to above for ContainersAreCopies

		It("has an unmodifiable adjacent nodes set view", func() {
			if containersMode != ContainersAreViews {
				Skip("Graph.AdjacentNodes() is not expected to return an unmodifiable view")
			}

			adjacentNodes := grph.AdjacentNodes(node1)
			Expect(adjacentNodes).To(BeNonMutableSet[int]())

			grph = putEdge(grph, node1, node2)
			testSet(adjacentNodes, node2)
		})

		// TODO: Write an equivalent test to above for ContainersAreCopies

		It("had an unmodifiable predecessors set view", func() {
			if containersMode != ContainersAreViews {
				Skip("Graph.Predecessors() is not expected to return an unmodifiable view")
			}

			predecessors := grph.Predecessors(node1)
			Expect(predecessors).To(BeNonMutableSet[int]())

			grph = putEdge(grph, node2, node1)
			testSet(predecessors, node2)
		})

		// TODO: Write an equivalent test to above for ContainersAreCopies

		It("has an unmodifiable successors set view", func() {
			if containersMode != ContainersAreViews {
				Skip("Graph.Successors() is not expected to return an unmodifiable view")
			}

			successors := grph.Successors(node1)
			Expect(successors).To(BeNonMutableSet[int]())

			grph = putEdge(grph, node1, node2)
			testSet(successors, node2)
		})

		// TODO: Write an equivalent test to above for ContainersAreCopies

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
			It("reports that both nodes are adjacent to each other", func() {
				grph = putEdge(grph, node1, node2)

				testSet(grph.AdjacentNodes(node1), node2)
				testSet(grph.AdjacentNodes(node2), node1)
			})

			It("reports that both nodes have a degree of 1", func() {
				grph = putEdge(grph, node1, node2)

				Expect(grph.Degree(node1)).To(Equal(1))
				Expect(grph.Degree(node2)).To(Equal(1))
			})
		})

		Context("when putting two connected edges", func() {
			BeforeEach(func() {
				grph = putEdge(grph, node1, node2)
				grph = putEdge(grph, node1, node3)
			})

			It("reports that the common node has a degree of 2", func() {
				Expect(grph.Degree(node1)).To(Equal(2))
			})

			It("reports the two unique nodes as adjacent to the common one", func() {
				testSet(grph.AdjacentNodes(node1), node2, node3)
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
	})

	mutableGraphTests(graphName, createGraph, containersMode)

	undirectedGraphTests(graphName, createGraph, containersMode)

	directedGraphsTests(graphName, createGraph, containersMode)

	allowsSelfLoopsGraphTests(graphName, createGraph, containersMode)

	disallowsSelfLoopsGraphTests(graphName, createGraph, containersMode)
}

func mutableGraphTests(
	graphName string,
	createGraph func() graph.Graph[int],
	containersMode ContainersMode,
) {
	_, mutable := createGraph().(graph.MutableGraph[int])
	if !mutable {
		// skip
		return
	}

	Context(fmt.Sprintf("%s: given a mutable graph", graphName), func() {
		var grph graph.MutableGraph[int]

		createGraphAsMutable := func() graph.MutableGraph[int] {
			grph := createGraph()

			g, ok := grph.(graph.MutableGraph[int])
			if !ok {
				panic("grph is expected to be mutable but was not")
			}

			return g
		}

		BeforeEach(func() {
			assertContainersMode(containersMode)

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

		Context("when removing an absent edge with an existing nodeU", func() {
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

		Context("when removing an absent edge with an existing nodeV", func() {
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

func undirectedGraphTests(
	graphName string,
	createGraph func() graph.Graph[int],
	containersMode ContainersMode,
) {
	if createGraph().IsDirected() {
		// skip
		return
	}

	Context(fmt.Sprintf("%s: given an undirected graph", graphName), func() {
		var grph graph.Graph[int]

		BeforeEach(func() {
			assertContainersMode(containersMode)

			grph = createGraph()
		})

		It("has an unmodifiable set view of unordered edges", func() {
			if containersMode != ContainersAreViews {
				Skip("Graph.Edges() is not expected to return an unmodifiable view")
			}

			edges := grph.Edges()
			Expect(edges).To(BeNonMutableSet[graph.EndpointPair[int]]())

			grph = putEdge(grph, node1, node2)
			testSingleEdgeForUndirectedGraph(edges)
		})

		// TODO: Write an equivalent test to above for ContainersAreCopies

		It("has an unmodifiable set view of unordered incident edges", func() {
			if containersMode != ContainersAreViews {
				Skip("Graph.IncidentEdges() is not expected to return an unmodifiable view")
			}

			incidentEdges := grph.IncidentEdges(node1)
			Expect(incidentEdges).To(BeNonMutableSet[graph.EndpointPair[int]]())

			grph = putEdge(grph, node1, node2)
			testSingleEdgeForUndirectedGraph(incidentEdges)
		})

		// TODO: Write an equivalent test to above for ContainersAreCopies

		Context("when putting one edge", func() {
			BeforeEach(func() {
				grph = putEdge(grph, node1, node2)
			})

			It("has just one unordered edge", func() {
				testSingleEdgeForUndirectedGraph(grph.Edges())
			})

			It("has an unordered incident edge connecting the first node to the second node", func() {
				testSingleEdgeForUndirectedGraph(grph.IncidentEdges(node1))
			})

			It("sees both nodes as predecessors of each other", func() {
				testSet(grph.Predecessors(node2), node1)
				testSet(grph.Predecessors(node1), node2)
			})

			It("sees both nodes as successors of each other", func() {
				testSet(grph.Successors(node1), node2)
				testSet(grph.Successors(node2), node1)
			})

			It("has an unordered incident edge connecting the second node to the first node", func() {
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

			It("sees the first node as being connected to the second in an unordered fashion", func() {
				Expect(grph.HasEdgeConnecting(node1, node2)).
					To(BeTrue())

				Expect(grph.HasEdgeConnectingEndpoints(graph.NewUnorderedEndpointPair(node1, node2))).
					To(BeTrue())
			})

			It("sees the second node as being connected to the first in an unordered fashion", func() {
				Expect(grph.HasEdgeConnecting(node2, node1)).
					To(BeTrue())

				Expect(grph.HasEdgeConnectingEndpoints(graph.NewUnorderedEndpointPair(node2, node1))).
					To(BeTrue())
			})

			It("does not see the first node as being connected to any other node", func() {
				Expect(grph.HasEdgeConnecting(node1, nodeNotInGraph)).
					To(BeFalse())
				Expect(grph.HasEdgeConnecting(nodeNotInGraph, node1)).
					To(BeFalse())

				Expect(grph.HasEdgeConnectingEndpoints(graph.NewUnorderedEndpointPair(node1, nodeNotInGraph))).
					To(BeFalse())
				Expect(grph.HasEdgeConnectingEndpoints(graph.NewUnorderedEndpointPair(nodeNotInGraph, node1))).
					To(BeFalse())
			})

			It("does not see the second node as being connected to any other node", func() {
				Expect(grph.HasEdgeConnecting(node2, nodeNotInGraph)).
					To(BeFalse())

				Expect(grph.HasEdgeConnectingEndpoints(graph.NewUnorderedEndpointPair(node2, nodeNotInGraph))).
					To(BeFalse())
			})

			Context("and trying to find that edge with ordered endpoints", func() {
				It("returns false", func() {
					Expect(grph.HasEdgeConnectingEndpoints(graph.NewOrderedEndpointPair(node1, node2))).
						To(BeFalse())
				})
			})
		})

		Context("when putting two connected edges", func() {
			BeforeEach(func() {
				grph = putEdge(grph, node1, node2)
				grph = putEdge(grph, node1, node3)
			})

			It("has two unordered edges sharing a common node", func() {
				testTwoEdgesForUndirectedGraphs(grph.Edges())
			})

			It("has two unordered incident edges connected to the common node", func() {
				testTwoEdgesForUndirectedGraphs(grph.IncidentEdges(node1))
			})
		})
	})
}

func directedGraphsTests(
	graphName string,
	createGraph func() graph.Graph[int],
	containersMode ContainersMode,
) {
	if !createGraph().IsDirected() {
		// skip
		return
	}

	Context(fmt.Sprintf("%s: given a directed graph", graphName), func() {
		var grph graph.Graph[int]

		BeforeEach(func() {
			assertContainersMode(containersMode)

			grph = createGraph()
		})

		It("has an unmodifiable set view of unordered edges", func() {
			if containersMode != ContainersAreViews {
				Skip("Graph.Edges() is not expected to return an unmodifiable view")
			}

			edges := grph.Edges()
			Expect(edges).To(BeNonMutableSet[graph.EndpointPair[int]]())

			grph = putEdge(grph, node1, node2)
			testSingleEdgeForDirectedGraph(edges)
		})

		// TODO: Write an equivalent test to above for ContainersAreCopies

		It("has an unmodifiable set view of ordered incident edges", func() {
			if containersMode != ContainersAreViews {
				Skip("Graph.IncidentEdges() is not expected to return an unmodifiable view")
			}

			incidentEdges := grph.IncidentEdges(node1)
			Expect(incidentEdges).To(BeNonMutableSet[graph.EndpointPair[int]]())

			grph = putEdge(grph, node1, node2)
			testSingleEdgeForDirectedGraph(incidentEdges)
		})

		// TODO: Write an equivalent test to above for ContainersAreCopies

		Context("when putting one edge", func() {
			BeforeEach(func() {
				grph = putEdge(grph, node1, node2)
			})

			It("has just one ordered edge", func() {
				testSingleEdgeForDirectedGraph(grph.Edges())
			})

			It("has an ordered incident edge connecting the first node to the second node", func() {
				testSingleEdgeForDirectedGraph(grph.IncidentEdges(node1))
			})
		})

		Context("when putting two connected edges", func() {
			It("has two ordered edges sharing a common node", func() {
				grph = putEdge(grph, node1, node2)
				grph = putEdge(grph, node1, node3)

				testTwoEdgesForDirectedGraphs(grph.Edges())
			})

			It("has two ordered incident edges connected to the common node", func() {
				testTwoEdgesForDirectedGraphs(grph.IncidentEdges(node1))
			})
		})
	})
}

func allowsSelfLoopsGraphTests(
	graphName string,
	createGraph func() graph.Graph[int],
	containersMode ContainersMode,
) {
	if !createGraph().AllowsSelfLoops() {
		// skip
		return
	}

	Context(fmt.Sprintf("%s: given a graph that allows self loops", graphName), func() {
		BeforeEach(func() {
			assertContainersMode(containersMode)
		})

		Context("when putting one self-loop edge", func() {
			It("sees the shared node as its own adjacent node", func() {
				grph := createGraph()

				grph = putEdge(grph, node1, node1)

				testSet(grph.AdjacentNodes(node1), node1)
			})
		})
	})
}

func disallowsSelfLoopsGraphTests(
	graphName string,
	createGraph func() graph.Graph[int],
	containersMode ContainersMode,
) {
	if createGraph().AllowsSelfLoops() {
		// skip
		return
	}

	Context(fmt.Sprintf("%s: given a graph that disallows self loops", graphName), func() {
		BeforeEach(func() {
			assertContainersMode(containersMode)
		})

		Context("when putting one self-loop edge", func() {
			It("panics", func() {
				grph := createGraph()

				Expect(func() { grph = putEdge(grph, node1, node1) }).
					To(PanicWith("self-loops are disallowed"))
			})
		})
	})
}

func assertContainersMode(containersMode ContainersMode) {
	if containersMode != ContainersAreViews &&
		containersMode != ContainersAreCopies {
		Fail(
			fmt.Sprintf(
				"containersMode is neither ContainersAreViews nor "+
					"ContainersAreCopies, but %d instead",
				containersMode))
	}
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
			graph.NewOrderedEndpointPair(
				nodeNotInGraph, nodeNotInGraph)))
	Expect(edges).ToNot(
		Contain(
			graph.NewUnorderedEndpointPair(
				nodeNotInGraph, nodeNotInGraph)))

	// Set.String()
	Expect(edges).To(HaveStringRepr("[]"))
}

func testSingleEdgeForUndirectedGraph(endpointPairs set.Set[graph.EndpointPair[int]]) {
	// Set.Len()
	Expect(endpointPairs).To(HaveLenOf(1))

	// Set.ForEach()
	//
	// Uses boolean assertions to avoid unreadable error messages from this
	// nested matcher.
	//
	// Note: on undirected graphs, this assertion checks that the endpoint
	// pairs emitted by .ForEach() are equal to any of the following:
	// - [[1, 2]]
	// - [[2, 1]]
	matcher := HaveForEachThatConsistsOf[graph.EndpointPair[int]](
		BeEquivalentToUsingEqualMethod(
			graph.NewUnorderedEndpointPair(node1, node2)))
	Expect(matcher.Match(endpointPairs)).To(
		BeTrue(),
		"to consist of %v according to graph.EndpointPair.Equal()",
		[]graph.EndpointPair[int]{
			graph.NewUnorderedEndpointPair(node1, node2),
		})

	// Set.Contains()
	Expect(endpointPairs).To(
		Contain(graph.NewUnorderedEndpointPair(node1, node2)))
	Expect(endpointPairs).ToNot(
		Contain(graph.NewOrderedEndpointPair(node1, node2)))
	Expect(endpointPairs).To(
		Contain(graph.NewUnorderedEndpointPair(node2, node1)))
	Expect(endpointPairs).ToNot(
		Contain(graph.NewOrderedEndpointPair(node2, node1)))
	Expect(endpointPairs).ToNot(
		Contain(
			graph.NewOrderedEndpointPair(
				nodeNotInGraph, nodeNotInGraph)))
	Expect(endpointPairs).ToNot(
		Contain(
			graph.NewUnorderedEndpointPair(
				nodeNotInGraph, nodeNotInGraph)))

	// Set.String()
	Expect(endpointPairs).To(
		HaveStringReprThatIsAnyOf("[[1, 2]]", "[[2, 1]]"))
}

func testSingleEdgeForDirectedGraph(endpointPairs set.Set[graph.EndpointPair[int]]) {
	// Set.Len()
	Expect(endpointPairs).To(HaveLenOf(1))

	// Set.ForEach()
	//
	// Uses boolean assertions to avoid unreadable error messages from this
	// nested matcher.
	//
	// Note: on directed graphs, it checks that the endpoint pairs emitted by
	// .ForEach() are exactly equal to [<1 -> 2>].
	matcher := HaveForEachThatConsistsOf[graph.EndpointPair[int]](
		BeEquivalentToUsingEqualMethod(
			graph.NewOrderedEndpointPair(node1, node2)))
	Expect(matcher.Match(endpointPairs)).To(
		BeTrue(),
		"to consist of %v according to graph.EndpointPair.Equal()",
		[]graph.EndpointPair[int]{
			graph.NewOrderedEndpointPair(node1, node2),
		})

	// Set.Contains()
	Expect(endpointPairs).To(
		Contain(graph.NewOrderedEndpointPair(node1, node2)))
	Expect(endpointPairs).ToNot(
		Contain(graph.NewUnorderedEndpointPair(node1, node2)))
	Expect(endpointPairs).ToNot(
		Contain(graph.NewOrderedEndpointPair(node1, node2)))
	Expect(endpointPairs).ToNot(
		Contain(graph.NewUnorderedEndpointPair(node1, node2)))
	Expect(endpointPairs).ToNot(
		Contain(
			graph.NewOrderedEndpointPair(
				nodeNotInGraph, nodeNotInGraph)))
	Expect(endpointPairs).ToNot(
		Contain(
			graph.NewUnorderedEndpointPair(
				nodeNotInGraph, nodeNotInGraph)))

	// Set.String()
	Expect(endpointPairs).To(HaveStringRepr("<1 -> 2>"))
}

func testTwoEdgesForUndirectedGraphs(endpointPairs set.Set[graph.EndpointPair[int]]) {
	// Set.Len()
	Expect(endpointPairs).To(HaveLenOf(2))

	// Set.ForEach()
	//
	// Uses boolean assertions to avoid unreadable error messages from this
	// nested matcher.
	//
	// Note: on undirected graphs, this assertion checks that the endpoint
	// pairs emitted by .ForEach() are equal to any of the following:
	// - [[1, 2], [1, 3]]
	// - [[1, 2], [3, 1]]
	// - [[2, 1], [1, 3]]
	// - [[2, 1], [3, 1]]
	// - [[1, 3], [1, 2]]
	// - [[1, 3], [2, 1]]
	// - [[3, 1], [1, 2]]
	// - [[3, 1], [2, 1]]
	matcher := HaveForEachThatConsistsOf[graph.EndpointPair[int]](
		BeEquivalentToUsingEqualMethod(
			graph.NewUnorderedEndpointPair(node1, node2)),
		BeEquivalentToUsingEqualMethod(
			graph.NewUnorderedEndpointPair(node1, node3)))
	Expect(matcher.Match(endpointPairs)).To(
		BeTrue(),
		"to consist of %v according to graph.EndpointPair.Equal()",
		[]graph.EndpointPair[int]{
			graph.NewUnorderedEndpointPair(node1, node2),
			graph.NewUnorderedEndpointPair(node1, node3),
		})

	// Set.Contains()
	Expect(endpointPairs).To(
		Contain(graph.NewUnorderedEndpointPair(node1, node2)))
	Expect(endpointPairs).To(
		Contain(graph.NewUnorderedEndpointPair(node1, node3)))
	Expect(endpointPairs).To(
		Contain(graph.NewUnorderedEndpointPair(node2, node1)))
	Expect(endpointPairs).To(
		Contain(graph.NewUnorderedEndpointPair(node3, node1)))
	Expect(endpointPairs).ToNot(
		Contain(graph.NewOrderedEndpointPair(node1, node2)))
	Expect(endpointPairs).ToNot(
		Contain(graph.NewOrderedEndpointPair(node1, node3)))
	Expect(endpointPairs).ToNot(
		Contain(graph.NewOrderedEndpointPair(node2, node1)))
	Expect(endpointPairs).ToNot(
		Contain(graph.NewOrderedEndpointPair(node3, node1)))
	Expect(endpointPairs).ToNot(
		Contain(
			graph.NewOrderedEndpointPair(
				nodeNotInGraph, nodeNotInGraph)))
	Expect(endpointPairs).ToNot(
		Contain(
			graph.NewUnorderedEndpointPair(
				nodeNotInGraph, nodeNotInGraph)))

	// Set.String()
	Expect(endpointPairs).To(
		HaveStringReprThatIsAnyOf(
			"[[1, 2], [1, 3]]",
			"[[1, 2], [3, 1]]",
			"[[2, 1], [1, 3]]",
			"[[2, 1], [3, 1]]",
			"[[1, 3], [1, 2]]",
			"[[1, 3], [2, 1]]",
			"[[3, 1], [1, 2]]",
			"[[3, 1], [2, 1]]"))
}

func testTwoEdgesForDirectedGraphs(endpointPairs set.Set[graph.EndpointPair[int]]) {
	// Set.Len()
	Expect(endpointPairs).To(HaveLenOf(2))

	// Set.ForEach()
	//
	// Uses boolean assertions to avoid unreadable error messages from this
	// nested matcher.
	//
	// Note: on undirected graphs, this assertion checks that the endpoint
	// pairs emitted by .ForEach() are equal to any of the following:
	// - [[1, 2], [1, 3]]
	// - [[1, 2], [3, 1]]
	// - [[2, 1], [1, 3]]
	// - [[2, 1], [3, 1]]
	// - [[1, 3], [1, 2]]
	// - [[1, 3], [2, 1]]
	// - [[3, 1], [1, 2]]
	// - [[3, 1], [2, 1]]
	//
	// For directed graphs, it checks that the endpoint pairs emitted by
	// .ForEach() are equal to any of the following:
	// - [<1 -> 2>, <1 -> 3>]
	// - [<1 -> 3>, <1 -> 2>]
	matcher := HaveForEachThatConsistsOf[graph.EndpointPair[int]](
		BeEquivalentToUsingEqualMethod(
			graph.NewOrderedEndpointPair(node1, node2)),
		BeEquivalentToUsingEqualMethod(
			graph.NewOrderedEndpointPair(node1, node3)))
	Expect(matcher.Match(endpointPairs)).To(
		BeTrue(),
		"to consist of %v according to graph.EndpointPair.Equal()",
		[]graph.EndpointPair[int]{
			graph.NewOrderedEndpointPair(node1, node2),
			graph.NewOrderedEndpointPair(node1, node3),
		})

	// Set.Contains()
	Expect(endpointPairs).To(
		Contain(graph.NewOrderedEndpointPair(node1, node2)))
	Expect(endpointPairs).To(
		Contain(graph.NewOrderedEndpointPair(node1, node3)))
	Expect(endpointPairs).ToNot(
		Contain(graph.NewOrderedEndpointPair(node2, node1)))
	Expect(endpointPairs).ToNot(
		Contain(graph.NewOrderedEndpointPair(node3, node1)))
	Expect(endpointPairs).ToNot(
		Contain(graph.NewUnorderedEndpointPair(node1, node2)))
	Expect(endpointPairs).ToNot(
		Contain(graph.NewUnorderedEndpointPair(node1, node3)))
	Expect(endpointPairs).ToNot(
		Contain(graph.NewUnorderedEndpointPair(node2, node1)))
	Expect(endpointPairs).ToNot(
		Contain(graph.NewUnorderedEndpointPair(node3, node1)))
	Expect(endpointPairs).ToNot(
		Contain(
			graph.NewOrderedEndpointPair(
				nodeNotInGraph, nodeNotInGraph)))
	Expect(endpointPairs).ToNot(
		Contain(
			graph.NewUnorderedEndpointPair(
				nodeNotInGraph, nodeNotInGraph)))

	// Set.String()
	Expect(endpointPairs).To(
		HaveStringReprThatIsAnyOf(
			"[<1 -> 2>, <1 -> 3>]",
			"[<1 -> 3>, <1 -> 2>]"))
}
