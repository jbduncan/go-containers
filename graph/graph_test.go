package graph_test

import (
	"fmt"

	"github.com/jbduncan/go-containers/graph"
	. "github.com/jbduncan/go-containers/internal/matchers"
	"github.com/jbduncan/go-containers/set"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// TODO: Move mutableGraphTests to a graphtester package
// TODO: Migrate mutableGraphTests to a struct with a constructor function like:
//       graphtester.ForMutableGraph(
//           name string,
//           newGraph func() graph.MutableGraph[int],
//           ...graphtester.Option options)
//       ...where options is any of:
//         - WhereContainersAreViews
//         - WhereContainersAreCopies
//       ...and the default is `WhereContainersAreViews`

var _ = Describe("Graphs", func() {
	mutableGraphTests(
		"Undirected graph",
		func() graph.MutableGraph[int] {
			return graph.Undirected[int]().Build()
		},
		ContainersAreViews)
	mutableGraphTests(
		"Undirected graph allowing self-loops",
		func() graph.MutableGraph[int] {
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

// mutableGraphTests produces a suite of Ginkgo test cases for testing implementations of the
// MutableGraph interface. MutableGraph instances created for testing are to have int nodes.
//
// Test cases that should be handled similarly in any graph implementation are included in this
// function; for example, testing that the `Nodes()` method returns the set of the nodes in the
// graph. Test cases related to specific implementations of the MutableGraph interface are
// explicitly not tested.
//
// TODO: Move to a public package for graph testing utilities
func mutableGraphTests(
	graphName string,
	createGraph func() graph.MutableGraph[int],
	containersMode ContainersMode,
) {
	// TODO: Consider replacing addNode with graphAsMutable().AddNode
	addNode := func(grph graph.Graph[int], node int) graph.Graph[int] {
		graphAsMutable := grph.(graph.MutableGraph[int])
		graphAsMutable.AddNode(node)

		return grph
	}
	// TODO: Consider replacing putEdge with graphAsMutable().PutEdge
	putEdge := func(grph graph.Graph[int], node1 int, node2 int) graph.Graph[int] {
		graphAsMutable := grph.(graph.MutableGraph[int])
		graphAsMutable.PutEdge(node1, node2)

		return grph
	}

	graphTests(graphName, func() graph.Graph[int] { return createGraph() }, addNode, putEdge, containersMode)
}

func graphTests(
	graphName string,
	createGraph func() graph.Graph[int],
	addNode func(g graph.Graph[int], n int) graph.Graph[int],
	putEdge func(g graph.Graph[int], n1 int, n2 int) graph.Graph[int],
	containersMode ContainersMode,
) {
	Context(fmt.Sprintf("%s: given a graph", graphName), func() {
		var grph graph.Graph[int]

		graphAsMutable := func() graph.MutableGraph[int] {
			result, _ := grph.(graph.MutableGraph[int])

			return result
		}

		skipIfGraphIsNotMutable := func() {
			_, mutable := grph.(graph.MutableGraph[int])

			if !mutable {
				Skip("Graph is not mutable")
			}
		}

		BeforeEach(func() {
			assertContainersMode(containersMode)

			grph = createGraph()
		})

		AfterEach(func() {
			validateGraphState(grph)
		})

		It("has nodes with ForEach() that emits nothing", func() {
			Expect(grph.Nodes()).To(HaveForEachThatEmitsNothing[int]())
		})

		It("has nodes with an empty string representation", func() {
			Expect(grph.Nodes()).To(HaveStringRepr("[]"))
		})

		It("has nodes that do not contain anything", func() {
			Expect(grph.Nodes()).ToNot(Or(Contain(node1), Contain(node2), Contain(node3)))
		})

		It("has nodes with length of 0", func() {
			Expect(grph.Nodes()).To(HaveLenOfZero())
		})

		It("has edges with ForEach() that emits nothing", func() {
			Expect(grph.Edges()).To(HaveForEachThatEmitsNothing[graph.EndpointPair[int]]())
		})

		It("has edges with an empty string representation", func() {
			Expect(grph.Edges()).To(HaveStringRepr("[]"))
		})

		It("has edges that do not contain anything", func() {
			Expect(grph.Edges()).ToNot(Contain(graph.NewUnorderedEndpointPair(node1, node2)))
		})

		It("has edges with length of 0", func() {
			Expect(grph.Edges()).To(HaveLenOfZero())
		})

		It("has an unmodifiable nodes set view", func() {
			if containersMode != ContainersAreViews {
				Skip("Graph.Nodes() is not expected to return an unmodifiable view")
			}

			nodes := grph.Nodes()
			Expect(nodes).To(BeNonMutableSet[int]())

			grph = addNode(grph, node1)
			Expect(nodes).To(Contain(node1))
		})

		// TODO: Write an equivalent test to above for ContainersAreCopies

		It("has an unmodifiable edges set view", func() {
			if containersMode != ContainersAreViews {
				Skip("Graph.Edges() is not expected to return an unmodifiable view")
			}

			edges := grph.Edges()
			Expect(edges).To(BeNonMutableSet[graph.EndpointPair[int]]())

			grph = putEdge(grph, node1, node2)
			Expect(edges).To(Contain(newEndpointPair(grph, node1, node2)))
		})

		// TODO: Write an equivalent test to above for ContainersAreCopies

		Context("when adding one node", func() {
			BeforeEach(func() {
				grph = addNode(grph, node1)
			})

			It("has nodes with ForEach() that emits just the node", func() {
				Expect(grph.Nodes()).To(HaveForEachThatConsistsOf[int](node1))
			})

			It("has nodes with a single element string representation", func() {
				Expect(grph.Nodes()).To(HaveStringRepr(fmt.Sprintf("[%d]", node1)))
			})

			It("has nodes that contains just the node", func() {
				Expect(grph.Nodes()).To(Contain(node1))
				Expect(grph.Nodes()).ToNot(Contain(node2))
				Expect(grph.Nodes()).ToNot(Contain(nodeNotInGraph))
			})

			It("has nodes with length of 1", func() {
				Expect(grph.Nodes()).To(HaveLenOf(1))
			})

			It("reports that the node has no adjacent nodes", func() {
				Expect(grph.AdjacentNodes(node1)).To(HaveForEachThatEmitsNothing[int]())
			})

			It("reports that the node has no predecessors", func() {
				Expect(grph.Predecessors(node1)).To(HaveForEachThatEmitsNothing[int]())
			})

			It("reports that the node has no successors", func() {
				Expect(grph.Successors(node1)).To(HaveForEachThatEmitsNothing[int]())
			})

			It("reports that the node has no incident edges", func() {
				Expect(grph.IncidentEdges(node1)).To(HaveForEachThatEmitsNothing[graph.EndpointPair[int]]())
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

			It("has an unmodifiable adjacent nodes set view", func() {
				if containersMode != ContainersAreViews {
					Skip("Graph.AdjacentNodes() is not expected to return an unmodifiable view")
				}

				adjacentNodes := grph.AdjacentNodes(node1)
				Expect(adjacentNodes).To(BeNonMutableSet[int]())

				grph = putEdge(grph, node1, node2)
				Expect(adjacentNodes).To(Contain(node2))
			})

			// TODO: Write an equivalent test to above for ContainersAreCopies

			It("had an unmodifiable predecessors set view", func() {
				if containersMode != ContainersAreViews {
					Skip("Graph.Predecessors() is not expected to return an unmodifiable view")
				}

				predecessors := grph.Predecessors(node1)
				Expect(predecessors).To(BeNonMutableSet[int]())

				grph = putEdge(grph, node2, node1)
				Expect(predecessors).To(Contain(node2))
			})

			// TODO: Write an equivalent test to above for ContainersAreCopies

			It("has an unmodifiable successors set view", func() {
				if containersMode != ContainersAreViews {
					Skip("Graph.Successors() is not expected to return an unmodifiable view")
				}

				successors := grph.Successors(node1)
				Expect(successors).To(BeNonMutableSet[int]())

				grph = putEdge(grph, node1, node2)
				Expect(successors).To(Contain(node2))
			})

			// TODO: Write an equivalent test to above for ContainersAreCopies

			It("has an unmodifiable incident edges set view", func() {
				if containersMode != ContainersAreViews {
					Skip("Graph.IncidentEdges() is not expected to return an unmodifiable view")
				}

				incidentEdges := grph.IncidentEdges(node1)
				Expect(incidentEdges).To(BeNonMutableSet[graph.EndpointPair[int]]())

				grph = putEdge(grph, node1, node2)
				Expect(incidentEdges).To(Contain(newEndpointPair(grph, node1, node2)))
			})

			// TODO: Write an equivalent test to above for ContainersAreCopies
		})

		Context("when adding two nodes", func() {
			BeforeEach(func() {
				grph = addNode(grph, node1)
				grph = addNode(grph, node2)
			})

			It("has nodes with ForEach() that emits both nodes", func() {
				Expect(grph.Nodes()).To(HaveForEachThatConsistsOf[int](node1, node2))
			})

			It("has nodes with a two element string representation", func() {
				Expect(grph.Nodes()).To(
					HaveStringRepr(
						BeElementOf(
							fmt.Sprintf("[%d, %d]", node1, node2),
							fmt.Sprintf("[%d, %d]", node2, node1))))
			})

			It("has nodes that contains both nodes", func() {
				Expect(grph.Nodes()).To(Contain(node1))
				Expect(grph.Nodes()).To(Contain(node2))
				Expect(grph.Nodes()).ToNot(Contain(nodeNotInGraph))
			})

			It("has nodes with length of 2", func() {
				Expect(grph.Nodes()).To(HaveLenOf(2))
			})
		})

		Context("when adding a new node", func() {
			It("returns true", func() {
				skipIfGraphIsNotMutable()

				Expect(graphAsMutable().AddNode(node1)).To(BeTrue())
			})
		})

		Context("when adding an existing node", func() {
			It("returns false", func() {
				skipIfGraphIsNotMutable()
				grph = addNode(grph, node1)

				Expect(graphAsMutable().AddNode(node1)).To(BeFalse())
			})
		})

		Context("when removing an existing node", func() {
			var removed bool

			BeforeEach(func() {
				skipIfGraphIsNotMutable()
				grph = putEdge(grph, node1, node2)
				grph = putEdge(grph, node3, node1)

				removed = graphAsMutable().RemoveNode(node1)
			})

			It("returns true", func() {
				Expect(removed).To(BeTrue())
			})

			It("it leaves the other nodes alone", func() {
				Expect(grph.Nodes()).To(HaveForEachThatConsistsOf[int](node2, node3))
			})

			It("removes its connections to its adjacent nodes", func() {
				Expect(grph.AdjacentNodes(node2)).To(HaveForEachThatEmitsNothing[int]())
				Expect(grph.AdjacentNodes(node3)).To(HaveForEachThatEmitsNothing[int]())
			})

			It("removes the connected edges", func() {
				Expect(grph.Edges()).To(HaveForEachThatEmitsNothing[graph.EndpointPair[int]]())
			})
		})

		Context("when removing an absent node", func() {
			var removed bool

			BeforeEach(func() {
				skipIfGraphIsNotMutable()
				grph = addNode(grph, node1)

				removed = graphAsMutable().RemoveNode(nodeNotInGraph)
			})

			It("returns false", func() {
				Expect(removed).To(BeFalse())
			})

			It("leaves all the nodes alone", func() {
				Expect(grph.Nodes()).To(HaveForEachThatConsistsOf[int](node1))
			})
		})

		Context("when putting one edge", func() {
			BeforeEach(func() {
				grph = putEdge(grph, node1, node2)
			})

			It("reports that both nodes are adjacent to each other", func() {
				Expect(grph.AdjacentNodes(node1)).To(HaveForEachThatConsistsOf[int](node2))
				Expect(grph.AdjacentNodes(node2)).To(HaveForEachThatConsistsOf[int](node1))
			})

			It("reports that both nodes have a degree of 1", func() {
				Expect(grph.Degree(node1)).To(Equal(1))
				Expect(grph.Degree(node2)).To(Equal(1))
			})

			It("has edges with ForEach() that emits just that edge", func() {
				// Uses boolean assertion to avoid unreadable error messages
				// from this nested matcher
				matcher := HaveForEachThatConsistsOf[graph.EndpointPair[int]](
					BeEquivalentToUsingEqualMethod(
						graph.NewUnorderedEndpointPair(node1, node2)))
				result, _ := matcher.Match(grph.Edges())
				Expect(result).To(
					BeTrue(),
					"to consist of [node1, node2] (unordered endpoint pair)")
			})

			It("has edges that contains just that edge", func() {
				// Uses boolean assertions to avoid unreadable error messages.
				containsNode1ToNode2 := grph.Edges().Contains(graph.NewUnorderedEndpointPair(node1, node2))
				Expect(containsNode1ToNode2).To(
					BeTrue(),
					"to contain [node1, node2] (unordered endpoint pair)")

				containsNode2ToNode1 := grph.Edges().Contains(graph.NewUnorderedEndpointPair(node2, node1))
				Expect(containsNode2ToNode1).To(
					BeTrue(),
					"to contain [node2, node1] (unordered endpoint pair)")

				containsNode1ToNode3 := grph.Edges().Contains(graph.NewUnorderedEndpointPair(node1, node3))
				Expect(containsNode1ToNode3).To(
					BeFalse(),
					"to not contain [node1, node3] (unordered endpoint pair)")

				containsNode3ToNode1 := grph.Edges().Contains(graph.NewUnorderedEndpointPair(node3, node1))
				Expect(containsNode3ToNode1).To(
					BeFalse(),
					"to not contain [node3, node1] (unordered endpoint pair)")
			})

			It("has edges with length of 1", func() {
				Expect(grph.Edges()).To(HaveLenOf(1))
			})

			It("has edges with a single element string representation", func() {
				Expect(grph.Edges()).To(HaveStringRepr(BeElementOf("[[1, 2]]", "[[2, 1]]")))
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
				Expect(grph.AdjacentNodes(node1)).To(HaveForEachThatConsistsOf[int](node2, node3))
			})

			It("returns both edges on iteration", func() {
				// Uses boolean assertion to avoid unreadable error messages from nested matcher
				matcher := HaveForEachThatConsistsOf[graph.EndpointPair[int]](
					BeEquivalentToUsingEqualMethod(
						graph.NewUnorderedEndpointPair(node1, node2)),
					BeEquivalentToUsingEqualMethod(
						graph.NewUnorderedEndpointPair(node1, node3)))
				result, _ := matcher.Match(grph.Edges())
				Expect(result).To(
					BeTrue(),
					"to consist of [node1, node2] and [node1, node3] (unordered endpoint pairs)")
			})

			It("has edges with length of 2", func() {
				Expect(grph.Edges()).To(HaveLenOf(2))
			})
		})

		Context("when putting two anti-parallel edges", func() {
			Context("and removing one of the nodes", func() {
				It("leaves the other node alone", func() {
					skipIfGraphIsNotMutable()
					grph = putEdge(grph, node1, node2)
					grph = putEdge(grph, node2, node1)
					graphAsMutable().RemoveNode(node1)

					Expect(grph.Nodes()).To(HaveForEachThatConsistsOf[int](node2))
				})

				It("removes both edges", func() {
					Expect(grph.Edges()).To(HaveForEachThatEmitsNothing[graph.EndpointPair[int]]())
				})
			})
		})

		Context("when removing an existing edge", func() {
			var removed bool

			BeforeEach(func() {
				skipIfGraphIsNotMutable()
				grph = putEdge(grph, node1, node2)
				grph = putEdge(grph, node1, node3)

				removed = graphAsMutable().RemoveEdge(node1, node2)
			})

			It("returns true", func() {
				Expect(removed).To(BeTrue())
			})

			It("removes the connection between its nodes", func() {
				Expect(grph.Successors(node1)).To(HaveForEachThatConsistsOf[int](node3))
				Expect(grph.Predecessors(node3)).To(HaveForEachThatConsistsOf[int](node1))
				Expect(grph.Predecessors(node2)).To(HaveForEachThatEmitsNothing[int]())
			})
		})

		Context("when removing an absent edge with an existing nodeU", func() {
			var removed bool

			BeforeEach(func() {
				skipIfGraphIsNotMutable()
				grph = putEdge(grph, node1, node2)

				removed = graphAsMutable().RemoveEdge(node1, nodeNotInGraph)
			})

			It("returns false", func() {
				Expect(removed).To(BeFalse())
			})

			It("leaves the existing nodes alone", func() {
				Expect(grph.Successors(node1)).To(Contain(node2))
				Expect(grph.Predecessors(node2)).To(Contain(node1))
			})
		})

		Context("when removing an absent edge with an existing nodeV", func() {
			var removed bool

			BeforeEach(func() {
				skipIfGraphIsNotMutable()
				grph = putEdge(grph, node1, node2)

				removed = graphAsMutable().RemoveEdge(nodeNotInGraph, node2)
			})

			It("returns false", func() {
				Expect(removed).To(BeFalse())
			})

			It("leaves the existing nodes alone", func() {
				Expect(grph.Successors(node1)).To(Contain(node2))
				Expect(grph.Predecessors(node2)).To(Contain(node1))
			})
		})

		Context("when removing an absent edge with two existing nodes", func() {
			var removed bool

			BeforeEach(func() {
				skipIfGraphIsNotMutable()
				grph = addNode(grph, node1)
				grph = addNode(grph, node2)

				removed = graphAsMutable().RemoveEdge(node1, node2)
			})

			It("returns false", func() {
				Expect(removed).To(BeFalse())
			})

			It("leaves the existing nodes alone", func() {
				Expect(grph.Nodes()).
					To(HaveForEachThatConsistsOf[int](node1, node2))
			})
		})

		Context("when finding the predecessors of an absent node", func() {
			It("returns an empty set", func() {
				Expect(grph.Predecessors(nodeNotInGraph)).
					To(HaveForEachThatEmitsNothing[int]())
			})
		})

		Context("when finding the successors of an absent node", func() {
			It("returns an empty set", func() {
				Expect(grph.Successors(nodeNotInGraph)).
					To(HaveForEachThatEmitsNothing[int]())
			})
		})

		Context("when finding the adjacent nodes of an absent node", func() {
			It("returns an empty set", func() {
				Expect(grph.AdjacentNodes(nodeNotInGraph)).
					To(HaveForEachThatEmitsNothing[int]())
			})
		})

		Context("when finding the incident edges of an absent node", func() {
			It("returns an empty set", func() {
				Expect(grph.IncidentEdges(nodeNotInGraph)).
					To(HaveForEachThatEmitsNothing[graph.EndpointPair[int]]())
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

	undirectedGraphTests(graphName, createGraph, addNode, putEdge, containersMode)
}

func undirectedGraphTests(
	graphName string,
	createGraph func() graph.Graph[int],
	addNode func(g graph.Graph[int], n int) graph.Graph[int],
	putEdge func(g graph.Graph[int], n1 int, n2 int) graph.Graph[int],
	containersMode ContainersMode,
) {
	Context(fmt.Sprintf("%s: given an undirected graph", graphName), func() {
		var grph graph.Graph[int]

		skipIfGraphAllowsSelfLoops := func() {
			if grph.AllowsSelfLoops() {
				Skip("Graph allows self-loops")
			}
		}

		skipIfGraphDisallowsSelfLoops := func() {
			if !grph.AllowsSelfLoops() {
				Skip("Graph disallows self-loops")
			}
		}

		BeforeEach(func() {
			assertContainersMode(containersMode)

			grph = createGraph()
			if grph.IsDirected() {
				Skip("Graph is directed")
			}
		})

		Context("when putting one edge", func() {
			BeforeEach(func() {
				grph = putEdge(grph, node1, node2)
			})

			It("sees both nodes as predecessors of each other", func() {
				Expect(grph.Predecessors(node2)).To(HaveForEachThatConsistsOf[int](node1))
				Expect(grph.Predecessors(node1)).To(HaveForEachThatConsistsOf[int](node2))
			})

			It("sees both nodes as successors of each other", func() {
				Expect(grph.Successors(node1)).To(HaveForEachThatConsistsOf[int](node2))
				Expect(grph.Successors(node2)).To(HaveForEachThatConsistsOf[int](node1))
			})

			It("has an incident edge connecting the first node to the second", func() {
				Expect(grph.IncidentEdges(node1)).To(
					HaveForEachThatConsistsOf[graph.EndpointPair[int]](
						BeEquivalentToUsingEqualMethod(
							graph.NewUnorderedEndpointPair(node1, node2))))
			})

			It("has an incident edge connecting the second node to the first", func() {
				Expect(grph.IncidentEdges(node2)).To(
					HaveForEachThatConsistsOf[graph.EndpointPair[int]](
						BeEquivalentToUsingEqualMethod(
							graph.NewUnorderedEndpointPair(node2, node1))))
			})

			It("has an incident edge with a single edge string representation", func() {
				Expect(grph.IncidentEdges(node1)).To(HaveStringRepr(BeElementOf("[[1, 2]]", "[[2, 1]]")))
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

		Context("when the graph disallows self-loops", func() {
			Context("and putting one self-loop edge", func() {
				It("panics", func() {
					skipIfGraphAllowsSelfLoops()

					Expect(func() { grph = putEdge(grph, node1, node1) }).
						To(PanicWith("self-loops are disallowed"))
				})
			})
		})

		Context("when the graph allows self-loops", func() {
			Context("and putting one self-loop edge", func() {
				It("sees the shared node as its own adjacent node", func() {
					skipIfGraphDisallowsSelfLoops()

					grph = putEdge(grph, node1, node1)

					Expect(grph.AdjacentNodes(node1)).To(HaveForEachThatConsistsOf[int](node1))
				})
			})
		})

		// TODO: Implement tests for stable ordering when NodeOrder()/IncidentEdgeOrder()
		//       is introduced.
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

func validateGraphState(grph graph.Graph[int]) {
	// TODO: Pending implementation of graph.CopyOf
	// expectStronglyEquivalent(graph, graph.CopyOf(graph))
	// TODO: Pending implementation of graph.ImmutableCopyOf
	// expectStronglyEquivalent(graph, graph.ImmutableCopyOf(graph))

	sanityCheckString(grph)

	allEndpointPairs := set.New[graph.EndpointPair[int]]()

	sanityCheckIntSet(grph.Nodes()).ForEach(func(node int) {
		// TODO: Pending implementation of Graph.String()
		// Expect(nodeString).To(ContainSubstring(strconv.Itoa(node)))

		sanityCheckConnections(grph, node, allEndpointPairs)
	})

	sanityCheckEdges(grph, allEndpointPairs)
}

// TODO: Move to its own test
func sanityCheckString(grph graph.Graph[int]) {
	// TODO: Pending implementation of Graph.String()
	// graphString := graph.String()
	// TODO: Pending implementation of Graph.String() and Graph.IsDirected()
	// Expect(graphString).To(ContainSubstring("isDirected: %v", graph.IsDirected()))
	// TODO: Pending implementation of Graph.String() and Graph.AllowsSelfLoops()
	// Expect(graphString).To(ContainSubstring("allowsSelfLoops: %v", graph.AllowsSelfLoops()))

	// TODO: Pending implementation of Graph.String()
	// nodeStart := strings.Index(graphString, "nodes:")
	// edgeStart := strings.Index(graphString, "edges:")
	// nodeString := graphString[nodeStart:edgeStart] // safe because contents are ASCII
}

// TODO: Move to own tests
func sanityCheckConnections(grph graph.Graph[int], node int, allEndpointPairs set.MutableSet[graph.EndpointPair[int]]) {
	if grph.IsDirected() {
		Expect(grph.Degree(node)).To(Equal(grph.InDegree(node) + grph.OutDegree(node)))
		Expect(grph.Predecessors(node)).To(HaveLenOf(grph.InDegree(node)))
		Expect(grph.Successors(node)).To(HaveLenOf(grph.OutDegree(node)))
	} else {
		nodeConnectedToSelf := grph.AdjacentNodes(node).Contains(node)
		selfLoopCount := 0
		if nodeConnectedToSelf {
			selfLoopCount = 1
		}
		Expect(grph.Degree(node)).To(Equal(grph.AdjacentNodes(node).Len() + selfLoopCount))
		Expect(grph.Predecessors(node)).To(HaveForEachThatConsistsOfElementsIn(grph.AdjacentNodes(node)))
		Expect(grph.Successors(node)).To(HaveForEachThatConsistsOfElementsIn(grph.AdjacentNodes(node)))
		Expect(grph.InDegree(node)).To(Equal(grph.Degree(node)))
		Expect(grph.OutDegree(node)).To(Equal(grph.Degree(node)))
	}

	sanityCheckIntSet(grph.AdjacentNodes(node)).ForEach(func(adjacentNode int) {
		if !grph.AllowsSelfLoops() {
			Expect(node).ToNot(Equal(adjacentNode))
		}
		Expect(
			grph.Predecessors(node).Contains(adjacentNode) ||
				grph.Successors(node).Contains(adjacentNode)).
			To(BeTrue())
	})

	sanityCheckIntSet(grph.Successors(node)).ForEach(func(successor int) {
		endpointPair := newEndpointPair(grph, node, successor)
		allEndpointPairs.Add(endpointPair)
		Expect(grph.Predecessors(successor)).To(Contain(node))
		Expect(grph.HasEdgeConnecting(node, successor)).To(BeTrue())
		Expect(grph.IncidentEdges(node)).To(Contain(endpointPair))
		if !grph.IsDirected() {
			reversedEndpointPair := newEndpointPair(grph, endpointPair.NodeV(), endpointPair.NodeU())
			Expect(grph.IncidentEdges(node)).To(Contain(reversedEndpointPair))
		}
	})

	sanityCheckEndpointPairSet(grph.IncidentEdges(node)).ForEach(func(endpoints graph.EndpointPair[int]) {
		if grph.IsDirected() {
			Expect(grph.HasEdgeConnecting(endpoints.Source(), endpoints.Target())).To(BeTrue())
		} else {
			Expect(grph.HasEdgeConnecting(endpoints.NodeU(), endpoints.NodeV())).To(BeTrue())
		}
	})
}

func sanityCheckEdges(grph graph.Graph[int], allEndpointPairs set.MutableSet[graph.EndpointPair[int]]) {
	sanityCheckEndpointPairSet(grph.Edges())
	// TODO: Move to own tests, one for empty graph, one for multiple-edge graph
	Expect(grph.Edges()).ToNot(Contain(newEndpointPair(grph, nodeNotInGraph, nodeNotInGraph)))
}

//lint:ignore U1000 this function will be used eventually
func expectStronglyEquivalent(first graph.Graph[int], second graph.Graph[int]) {
	// Properties not covered by Graph.Equal()
	Expect(first.AllowsSelfLoops()).To(Equal(second.AllowsSelfLoops()))
	// TODO: Pending implementation of Graph.NodeOrder()
	// Expect(first).To(haveNodeOrder(second.NodeOrder()))

	// TODO: Pending implementation of Graph.Equal()
	// Expect(first).To(beGraphEqualTo(second))
}

// TODO: Remove these sanity checks once all methods of all set-returning
//  methods of graph.Graph are tested.

// In some cases, graphs may return custom sets that define their own method implementations. Verify that
// these sets are consistent with the elements produced by their ForEach.
func sanityCheckIntSet(values set.Set[int]) set.Set[int] {
	Expect(values).To(HaveLenOf(forEachCount(values)))
	values.ForEach(func(elem int) {
		Expect(values).To(Contain(elem))
	})
	Expect(values).ToNot(Contain(nodeNotInGraph))
	// TODO: Pending introduction of set.CopyOf()
	// Expect(values).To(HaveForEachThatConsistsOfElementsIn(values.CopyOf(values)))
	return values
}

// TODO: Move below function to its own tests

// In some cases, graphs may return custom sets that define their own method implementations. Verify that
// these sets are consistent with the elements produced by their ForEach.
func sanityCheckEndpointPairSet(values set.Set[graph.EndpointPair[int]]) set.Set[graph.EndpointPair[int]] {
	Expect(values).To(HaveLenOf(forEachCount(values)))
	values.ForEach(func(elem graph.EndpointPair[int]) {
		Expect(values).To(Contain(elem))
	})
	Expect(values).ToNot(Contain(graph.NewOrderedEndpointPair(nodeNotInGraph, nodeNotInGraph)))
	Expect(values).ToNot(Contain(graph.NewUnorderedEndpointPair(nodeNotInGraph, nodeNotInGraph)))
	// TODO: Pending introduction of set.CopyOf()
	// Expect(values).To(HaveForEachThatConsistsOfElementsIn(values.CopyOf(values)))
	return values
}

func newEndpointPair[N comparable](grph graph.Graph[N], nodeU N, nodeV N) graph.EndpointPair[N] {
	if grph.IsDirected() {
		return graph.NewOrderedEndpointPair(nodeU, nodeV)
	}
	return graph.NewUnorderedEndpointPair(nodeU, nodeV)
}

func forEachCount[T comparable](set set.Set[T]) int {
	var result int

	set.ForEach(func(elem T) {
		result++
	})

	return result
}
