package graph_test

import (
	"fmt"

	"github.com/onsi/gomega/types"
	"go-containers/container/graph"
	"go-containers/container/set"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "go-containers/internal/matchers"
)

var _ = Describe("Undirected mutable graph", func() {
	mutableGraphTests(
		func() graph.MutableGraph[int] {
			return graph.Undirected[int]().Build()
		},
		ContainersAreViews)
	mutableGraphTests(
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
	createGraph func() graph.MutableGraph[int],
	containersMode ContainersMode) {

	addNode := func(grph graph.Graph[int], node int) graph.Graph[int] {
		graphAsMutable := grph.(graph.MutableGraph[int])
		graphAsMutable.AddNode(node)

		return grph
	}
	putEdge := func(grph graph.Graph[int], node1 int, node2 int) graph.Graph[int] {
		graphAsMutable := grph.(graph.MutableGraph[int])
		graphAsMutable.PutEdge(node1, node2)

		return grph
	}

	graphTests(func() graph.Graph[int] { return createGraph() }, addNode, putEdge, containersMode)
}

func graphTests(
	createGraph func() graph.Graph[int],
	addNode func(g graph.Graph[int], n int) graph.Graph[int],
	putEdge func(g graph.Graph[int], n1 int, n2 int) graph.Graph[int],
	containersMode ContainersMode) {

	Context("given a graph", func() {
		var (
			grph graph.Graph[int]
		)

		graphIsMutable := func() bool {
			_, result := grph.(graph.MutableGraph[int])

			return result
		}

		graphAsMutable := func() graph.MutableGraph[int] {
			result, _ := grph.(graph.MutableGraph[int])

			return result
		}

		skipIfGraphIsNotMutable := func() {
			if !graphIsMutable() {
				Skip("Graph is not mutable.")
			}
		}

		BeforeEach(func() {
			assertContainersMode(containersMode)

			grph = createGraph()
		})

		AfterEach(func() {
			validateGraphState(grph)
		})

		It("contains no nodes", func() {
			Expect(grph.Nodes()).To(beSetThatIsEmpty[int]())
		})

		It("contains no edges", func() {
			Expect(grph.Edges()).To(beSetThatIsEmpty[graph.EndpointPair[int]]())
		})

		It("has an unmodifiable nodes set view", func() {
			if containersMode != ContainersAreViews {
				Skip("Graph.Nodes() is not expected to return an unmodifiable view")
			}

			nodes := grph.Nodes()
			Expect(nodes).To(beSetThatIsNotMutable[int]())

			grph = addNode(grph, node1)
			// TODO: Pending implementation of Graph.Nodes().Contains()
			//Expect(nodes).To(Contain(node1))
		})

		// TODO: Write an equivalent test to above for ContainersAreCopies

		It("has an unmodifiable edges set view", func() {
			if containersMode != ContainersAreViews {
				Skip("Graph.Edges() is not expected to return an unmodifiable view")
			}

			edges := grph.Edges()
			Expect(edges).To(beSetThatIsNotMutable[graph.EndpointPair[int]]())

			grph = putEdge(grph, node1, node2)
			// TODO: Pending implementation of Graph.Edges()
			//Expect(edges).To(Contain(newEndpointPair(grph, node1, node2)))
		})

		// TODO: Write an equivalent test to above for ContainersAreCopies

		Context("when adding one node", func() {
			BeforeEach(func() {
				grph = addNode(grph, node1)
			})

			It("contains just the node", func() {
				Expect(grph.Nodes()).To(beSetThatConsistsOf(node1))
			})

			It("reports that the node has no adjacent nodes", func() {
				Expect(grph.AdjacentNodes(node1)).To(beSetThatIsEmpty[int]())
			})

			It("reports that the node has no predecessors", func() {
				Expect(grph.Predecessors(node1)).To(beSetThatIsEmpty[int]())
			})

			It("reports that the node has no successors", func() {
				Expect(grph.Successors(node1)).To(beSetThatIsEmpty[int]())
			})

			It("reports that the node has no incident edges", func() {
				Expect(grph.IncidentEdges(node1)).To(beSetThatIsEmpty[graph.EndpointPair[int]]())
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
				Expect(adjacentNodes).To(beSetThatIsNotMutable[int]())

				grph = putEdge(grph, node1, node2)
				Expect(adjacentNodes).To(Contain(node2))
			})

			// TODO: Write an equivalent test to above for ContainersAreCopies

			It("had an unmodifiable predecessors set view", func() {
				if containersMode != ContainersAreViews {
					Skip("Graph.Predecessors() is not expected to return an unmodifiable view")
				}

				predecessors := grph.Predecessors(node1)
				Expect(predecessors).To(beSetThatIsNotMutable[int]())

				grph = putEdge(grph, node2, node1)
				Expect(predecessors).To(Contain(node2))
			})

			// TODO: Write an equivalent test to above for ContainersAreCopies

			It("has an unmodifiable successors set view", func() {
				if containersMode != ContainersAreViews {
					Skip("Graph.Successors() is not expected to return an unmodifiable view")
				}

				successors := grph.Successors(node1)
				Expect(successors).To(beSetThatIsNotMutable[int]())

				grph = putEdge(grph, node1, node2)
				Expect(successors).To(Contain(node2))
			})

			// TODO: Write an equivalent test to above for ContainersAreCopies

			It("has an unmodifiable incident edges set view", func() {
				if containersMode != ContainersAreViews {
					Skip("Graph.IncidentEdges() is not expected to return an unmodifiable view")
				}

				incidentEdges := grph.IncidentEdges(node1)
				Expect(incidentEdges).To(beSetThatIsNotMutable[graph.EndpointPair[int]]())

				grph = putEdge(grph, node1, node2)
				// TODO: Pending implementation of Graph.IncidentEdges().Contains()
				//Expect(incidentEdges).To(Contain(newEndpointPair(grph, node1, node2)))
			})

			// TODO: Write an equivalent test to above for ContainersAreCopies
		})

		Context("when adding two nodes", func() {
			It("contains both nodes", func() {
				grph = addNode(grph, node1)
				grph = addNode(grph, node2)

				Expect(grph.Nodes()).To(beSetThatConsistsOf(node1, node2))
			})
		})

		Context("when adding a new node", func() {
			It("returns true", func() {
				skipIfGraphIsNotMutable()

				Expect(graphAsMutable().AddNode(node1)).To(BeTrue())
			})
		})

		Context("when adding an existing node", Ordered, func() {
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
				Expect(grph.Nodes()).To(beSetThatConsistsOf(node2, node3))
			})

			It("removes its connections to its adjacent nodes", func() {
				Expect(grph.AdjacentNodes(node2)).To(beSetThatIsEmpty[int]())
				Expect(grph.AdjacentNodes(node3)).To(beSetThatIsEmpty[int]())
			})

			It("removes the connected edges", func() {
				Expect(grph.Edges()).To(beSetThatIsEmpty[graph.EndpointPair[int]]())
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
				Expect(grph.Nodes()).To(beSetThatConsistsOf(node1))
			})
		})

		Context("when adding one edge", func() {
			BeforeEach(func() {
				grph = putEdge(grph, node1, node2)
			})

			It("reports that both nodes are adjacent to each other", func() {
				Expect(grph.AdjacentNodes(node1)).To(beSetThatConsistsOf(node2))
				Expect(grph.AdjacentNodes(node2)).To(beSetThatConsistsOf(node1))
			})

			It("reports that both nodes have a degree of 1", func() {
				Expect(grph.Degree(node1)).To(Equal(1))
				Expect(grph.Degree(node2)).To(Equal(1))
			})
		})

		Context("when adding two connected edges", func() {
			It("reports that the common node has a degree of 2", func() {
				grph = putEdge(grph, node1, node2)
				grph = putEdge(grph, node1, node3)

				Expect(grph.Degree(node1)).To(Equal(2))
			})

			It("reports the two unique nodes as adjacent to the common one", func() {
				grph = putEdge(grph, node1, node2)
				grph = putEdge(grph, node1, node3)

				Expect(grph.AdjacentNodes(node1)).To(beSetThatConsistsOf(node2, node3))
			})
		})

		Context("when adding two anti-parallel edges", func() {
			Context("and removing one of the nodes", func() {
				It("leaves the other node alone", func() {
					skipIfGraphIsNotMutable()
					grph = putEdge(grph, node1, node2)
					grph = putEdge(grph, node2, node1)
					graphAsMutable().RemoveNode(node1)

					Expect(grph.Nodes()).To(beSetThatConsistsOf(node2))
				})

				It("removes both edges", func() {
					Expect(grph.Edges()).To(beSetThatIsEmpty[graph.EndpointPair[int]]())
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
				Expect(grph.Successors(node1)).To(beSetThatConsistsOf(node3))
				Expect(grph.Predecessors(node3)).To(beSetThatConsistsOf(node1))
				Expect(grph.Predecessors(node2)).To(beSetThatIsEmpty[int]())
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
			})
		})

		Context("when finding the predecessors of an absent node", func() {
			It("returns an empty set", func() {
				Expect(grph.Predecessors(nodeNotInGraph)).
					To(beSetThatIsEmpty[int]())
			})
		})

		Context("when finding the successors of an absent node", func() {
			It("returns an empty set", func() {
				Expect(grph.Successors(nodeNotInGraph)).
					To(beSetThatIsEmpty[int]())
			})
		})

		Context("when finding the adjacent nodes of an absent node", func() {
			It("returns an empty set", func() {
				Expect(grph.AdjacentNodes(nodeNotInGraph)).
					To(beSetThatIsEmpty[int]())
			})
		})

		Context("when finding the incident edges of an absent node", func() {
			It("returns an empty set", func() {
				Expect(grph.IncidentEdges(nodeNotInGraph)).
					To(beSetThatIsEmpty[graph.EndpointPair[int]]())
			})
		})

		Context("when finding the degree of an absent node", func() {
			It("returns zero", func() {
				Expect(grph.Degree(nodeNotInGraph)).
					To(BeZero())
			})
		})

		Context("when finding the in degree of an absent node", func() {
			It("returns zero", func() {
				Expect(grph.InDegree(nodeNotInGraph)).
					To(BeZero())
			})
		})

		Context("when finding the out degree of of an absent node", func() {
			It("returns zero", func() {
				Expect(grph.OutDegree(nodeNotInGraph)).
					To(BeZero())
			})
		})
	})

	undirectedGraphTests(createGraph, addNode, putEdge, containersMode)
}

func undirectedGraphTests(
	createGraph func() graph.Graph[int],
	addNode func(g graph.Graph[int], n int) graph.Graph[int],
	putEdge func(g graph.Graph[int], n1 int, n2 int) graph.Graph[int],
	containersMode ContainersMode) {

	Context("given an undirected graph", func() {

		var (
			grph graph.Graph[int]
		)

		BeforeEach(func() {
			assertContainersMode(containersMode)

			grph = createGraph()
			if grph.IsDirected() {
				Skip("graph is not undirected")
			}
		})

		Context("when adding one edge", func() {
			BeforeEach(func() {
				grph = putEdge(grph, node1, node2)
			})

			It("sees both nodes as predecessors of each other", func() {
				Expect(grph.Predecessors(node2)).To(beSetThatConsistsOf(node1))
				Expect(grph.Predecessors(node1)).To(beSetThatConsistsOf(node2))
			})

			It("sees both nodes as successors of each other", func() {
				Expect(grph.Successors(node1)).To(beSetThatConsistsOf(node2))
				Expect(grph.Successors(node2)).To(beSetThatConsistsOf(node1))
			})

			It("has an incident edge connecting the first node to the second", func() {
				Expect(grph.IncidentEdges(node1)).To(
					// TODO: When EndpointPair has an Equal method, replace this assertion with a custom one
					//       that checks that the set contains only one element, where the element is "equal"
					//       according to EndpointPair.Equal. (The name "BeEquivalentTo" is already taken by
					//       Gomega. Maybe "EqualAccordingToEqualMethod"?)
					beSetThatConsistsOf(graph.NewUnorderedEndpointPair(node1, node2)))
			})

			It("has an incident edge connecting the second node to the first", func() {
				Expect(grph.IncidentEdges(node2)).To(
					// TODO: When EndpointPair has an Equal method, replace this assertion with a custom one
					//       that checks that the set contains only one element, where the element is "equal"
					//       according to EndpointPair.Equal. (The name "BeEquivalentTo" is already taken by
					//       Gomega. Maybe "EqualAccordingToEqualMethod"?)
					beSetThatConsistsOf(graph.NewUnorderedEndpointPair(node2, node1)))
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
	})
}

func assertContainersMode(containersMode ContainersMode) {
	if containersMode != ContainersAreViews &&
		containersMode != ContainersAreCopies {
		Fail(
			fmt.Sprintf(
				"containersMode returned neither ContainersAreViews nor "+
					"ContainersAreCopies, but %d instead",
				containersMode))
	}
}

func validateGraphState(grph graph.Graph[int]) {
	// TODO: Pending implementation of graph.CopyOf
	//expectStronglyEquivalent(graph, graph.CopyOf(graph))
	// TODO: Pending implementation of graph.ImmutableCopyOf
	//expectStronglyEquivalent(graph, graph.ImmutableCopyOf(graph))

	sanityCheckString(grph)

	allEndpointPairs := set.New[graph.EndpointPair[int]]()

	sanityCheckSet(grph.Nodes()).ForEach(func(node int) {
		// TODO: Pending implementation of Graph.String()
		//Expect(nodeString).To(ContainSubstring(strconv.Itoa(node)))

		sanityCheckConnections(grph, node, allEndpointPairs)
		sanityCheckEdges(grph)
	})
}

func sanityCheckString(grph graph.Graph[int]) {
	_ = grph
	// TODO: Pending implementation of Graph.String()
	//graphString := graph.String()
	// TODO: Pending implementation of Graph.String() and Graph.IsDirected()
	//Expect(graphString).To(ContainSubstring("isDirected: %v", graph.IsDirected()))
	// TODO: Pending implementation of Graph.String() and Graph.AllowsSelfLoops()
	//Expect(graphString).To(ContainSubstring("allowsSelfLoops: %v", graph.AllowsSelfLoops()))

	// TODO: Pending implementation of Graph.String()
	//nodeStart := strings.Index(graphString, "nodes:")
	//edgeStart := strings.Index(graphString, "edges:")
	//nodeString := graphString[nodeStart:edgeStart] // safe because contents are ASCII
}

func sanityCheckConnections(grph graph.Graph[int], node int, allEndpointPairs set.MutableSet[graph.EndpointPair[int]]) {
	// TODO: Pending implementation of many Graph methods
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
		Expect(grph.Predecessors(node)).To(beSetThatConsistsOfElementsIn(grph.AdjacentNodes(node)))
		Expect(grph.Successors(node)).To(beSetThatConsistsOfElementsIn(grph.AdjacentNodes(node)))
		Expect(grph.InDegree(node)).To(Equal(grph.Degree(node)))
		Expect(grph.OutDegree(node)).To(Equal(grph.Degree(node)))
	}

	sanityCheckSet(grph.AdjacentNodes(node)).ForEach(func(adjacentNode int) {
		// TODO: Pending implementation of Graph.AllowSelfLoops()
		//if !graph.AllowsSelfLoops() {
		//	Expect(node).ToNot(Equal(adjacentNode))
		//}
		Expect(
			grph.Predecessors(node).Contains(adjacentNode) ||
				grph.Successors(node).Contains(adjacentNode)).
			To(BeTrue())
	})

	sanityCheckSet(grph.Successors(node)).ForEach(func(successor int) {
		allEndpointPairs.Add(newEndpointPair(grph, node, successor))
		Expect(grph.Predecessors(successor)).To(Contain(node))
		Expect(grph.HasEdgeConnecting(node, successor)).To(BeTrue())
		// TODO: Pending implementation of Graph.IncidentEdges().Contains()
		//Expect(graph.IncidentEdges(node)).To(Contain(newEndpointPair(graph, node, successor)))
	})

	sanityCheckSet(grph.IncidentEdges(node)).ForEach(func(endpoints graph.EndpointPair[int]) {
		if grph.IsDirected() {
			Expect(grph.HasEdgeConnecting(endpoints.Source(), endpoints.Target())).To(BeTrue())
		} else {
			Expect(grph.HasEdgeConnecting(endpoints.NodeU(), endpoints.NodeV())).To(BeTrue())
		}
	})
}

func sanityCheckEdges(grph graph.Graph[int]) {
	sanityCheckSet(grph.Edges())
	Expect(grph.Edges()).ToNot(Contain(newEndpointPair(grph, nodeNotInGraph, nodeNotInGraph)))
	// TODO: Pending implementation of Graph.Edges()
	//Expect(grph.Edges()).To(beSetThatConsistsOfElementsIn[EndpointPair[int]](allEndpointPairs))
}

func expectStronglyEquivalent(first graph.Graph[int], second graph.Graph[int]) {
	// Properties not covered by Graph.Equal()
	// TODO: Pending implementation of Graph.AllowsSelfLoops()
	//Expect(first.AllowsSelfLoops()).To(Equal(second.AllowsSelfLoops()))
	// TODO: Pending implementation of Graph.NodeOrder()
	//Expect(first).To(haveNodeOrder(second.NodeOrder()))

	// TODO: Pending implementation of Graph.Equal()
	//Expect(first).To(beGraphEqualTo(second))
}

func newEndpointPair[N comparable](grph graph.Graph[N], nodeU N, nodeV N) graph.EndpointPair[N] {
	if grph.IsDirected() {
		return graph.NewOrderedEndpointPair(nodeU, nodeV)
	}
	return graph.NewUnorderedEndpointPair(nodeU, nodeV)
}

// In some cases, graphs may return custom sets that define their own method implementations. Verify that
// these sets are consistent with the elements produced by their ForEach.
func sanityCheckSet[N comparable](set set.Set[N]) set.Set[N] {
	// TODO: Pending implementation of keySet.Len()
	// Expect(set).To(HaveLenOf(forEachCount(set)))
	// TODO: Pending implementation of keySet.Contains()
	//set.ForEach(func(elem N) {
	//	Expect(set).To(Contain(elem))
	//})
	//Expect(set).ToNot(Contain(nodeNotInGraph))
	// TODO: Pending introduction of Set.String()
	//Expect(set).To(HaveStringConsistingOfElementsIn(set))
	// TODO: Pending introduction of Set.Equal()
	//Expect(set).To(beSetThatConsistsOfElementsIn(set.CopyOf(set)))
	return set
}

func beSetThatConsistsOf[N comparable](first N, others ...N) types.GomegaMatcher {
	all := []N{first}
	all = append(all, others...)

	return WithTransform(ForEachToSlice[N], ConsistOf(all))
}

func beSetThatConsistsOfElementsIn[T comparable](set set.Set[T]) types.GomegaMatcher {
	return WithTransform(ForEachToSlice[T], ConsistOf(ForEachToSlice(set)))
}

func beSetThatIsEmpty[T comparable]() types.GomegaMatcher {
	return WithTransform(ForEachToSlice[T], BeEmpty())
}

func beSetThatIsNotMutable[T comparable]() types.GomegaMatcher {
	return WithTransform(
		func(s set.Set[T]) bool {
			_, mutable := s.(set.MutableSet[T])
			return mutable
		},
		BeFalse())
}

func forEachCount[N comparable](set set.Set[N]) int {
	var result int

	set.ForEach(func(elem N) {
		result++
	})

	return result
}
