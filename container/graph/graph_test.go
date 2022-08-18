package graph

import (
	"fmt"

	"github.com/onsi/gomega/types"
	"go-containers/container/set"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "go-containers/internal/matchers"
)

var _ = Describe("Undirected mutable graph", func() {
	graphTests(
		func() Graph[int] {
			return Undirected[int]().Build()
		},
		func(graph Graph[int], n int) Graph[int] {
			graphAsMutable, _ := graph.(MutableGraph[int])
			graphAsMutable.AddNode(n)

			return graph
		},
		func(graph Graph[int], node1 int, node2 int) Graph[int] {
			graphAsMutable, _ := graph.(MutableGraph[int])
			graphAsMutable.PutEdge(node1, node2)

			return graph
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

// graphTests produces a suite of Ginkgo test cases for testing implementations of the Graph
// interface. Graph instances created for testing should have int nodes.
//
// Test cases that should be handled similarly in any graph implementation are included in this
// function; for example, testing that the `Nodes()` method returns the set of the nodes in the
// graph. The following test cases are explicitly not tested:
//   - Test cases related to whether the graph is directed or undirected.
//   - Test cases related to specific implementations of the Graph interface.
//
// TODO: Move to a public package for graph testing utilities
func graphTests(
	createGraph func() Graph[int],
	addNode func(g Graph[int], n int) Graph[int],
	putEdge func(g Graph[int], n1 int, n2 int) Graph[int],
	containersMode ContainersMode) {

	Context("given a graph", func() {
		var (
			graph Graph[int]
		)

		graphIsMutable := func() bool {
			_, result := graph.(MutableGraph[int])

			return result
		}

		graphAsMutable := func() MutableGraph[int] {
			result, _ := graph.(MutableGraph[int])

			return result
		}

		skipIfGraphIsNotMutable := func() {
			if !graphIsMutable() {
				Skip("Graph is not mutable.")
			}
		}

		BeforeEach(func() {
			assertContainersMode(containersMode)

			graph = createGraph()
		})

		AfterEach(func() {
			validateGraphState(graph)
		})

		It("contains no nodes", func() {
			Expect(graph.Nodes()).To(beSetThatIsEmpty[int]())
		})

		It("contains no edges", func() {
			Expect(graph.Edges()).To(beSetThatIsEmpty[EndpointPair[int]]())
		})

		It("has an unmodifiable nodes set view", func() {
			if containersMode != ContainersAreViews {
				Skip("Graph.Nodes() is not expected to return an unmodifiable view")
			}

			nodes := graph.Nodes()
			Expect(nodes).To(beSetThatIsNotMutable[int]())

			graph = addNode(graph, node1)
			// TODO: Pending implementation of keySet.Contains()
			//Expect(nodes).To(Contain(node1))
		})

		It("has an unmodifiable edges set view", func() {
			if containersMode != ContainersAreViews {
				Skip("Graph.Edges() is not expected to return an unmodifiable view")
			}

			edges := graph.Edges()
			Expect(edges).To(beSetThatIsNotMutable[EndpointPair[int]]())

			graph = putEdge(graph, node1, node2)
			// TODO: Pending uncommenting of newEndpointPair function and implementation of
			//		 Graph.Edges()
			//Expect(edges).To(Contain(newEndpointPair(graph, node1, node2)))
		})

		Context("when adding one node", func() {
			BeforeEach(func() {
				graph = addNode(graph, node1)
			})

			It("contains just the node", func() {
				Expect(graph.Nodes()).To(beSetThatConsistsOf(node1))
			})

			It("reports that the node has no adjacent nodes", func() {
				Expect(graph.AdjacentNodes(node1)).To(beSetThatIsEmpty[int]())
			})

			It("reports that the node has no predecessors", func() {
				Expect(graph.Predecessors(node1)).To(beSetThatIsEmpty[int]())
			})

			It("reports that the node has no successors", func() {
				Expect(graph.Successors(node1)).To(beSetThatIsEmpty[int]())
			})

			It("reports that the node has no incident edges", func() {
				Expect(graph.IncidentEdges(node1)).To(beSetThatIsEmpty[EndpointPair[int]]())
			})

			It("reports that the node has a degree of 0", func() {
				Expect(graph.Degree(node1)).To(BeZero())
			})

			It("reports that the node has an in degree of 0", func() {
				Expect(graph.InDegree(node1)).To(BeZero())
			})

			It("reports that the node has an out degree of 0", func() {
				Expect(graph.OutDegree(node1)).To(BeZero())
			})

			It("has an unmodifiable adjacent nodes set view", func() {
				if containersMode != ContainersAreViews {
					Skip("Graph.AdjacentNodes() is not expected to return an unmodifiable view")
				}

				adjacentNodes := graph.AdjacentNodes(node1)
				Expect(adjacentNodes).To(beSetThatIsNotMutable[int]())

				graph = putEdge(graph, node1, node2)
				// TODO: Pending uncommenting of newEndpointPair function and implementation of
				//       Graph.AdjacentNodes()
				//Expect(adjacentNodes).To(Contain(newEndpointPair(graph, node1, node2)))
			})

			It("had an unmodifiable predecessors set view", func() {
				if containersMode != ContainersAreViews {
					Skip("Graph.Predecessors() is not expected to return an unmodifiable view")
				}

				predecessors := graph.Predecessors(node1)
				Expect(predecessors).To(beSetThatIsNotMutable[int]())

				graph = putEdge(graph, node2, node1)
				// TODO: Pending implementation of Graph.Predecessors()
				//Expect(predecessors).To(Contain(node2))
			})

			It("has an unmodifiable successors set view", func() {
				if containersMode != ContainersAreViews {
					Skip("Graph.Successors() is not expected to return an unmodifiable view")
				}

				successors := graph.Successors(node1)
				Expect(successors).To(beSetThatIsNotMutable[int]())

				graph = putEdge(graph, node1, node2)
				// TODO: Pending implementation of Graph.Successors()
				//Expect(successors).To(Contain(node2))
			})

			It("has an unmodifiable incident edges set view", func() {
				if containersMode != ContainersAreViews {
					Skip("Graph.IncidentEdges() is not expected to return an unmodifiable view")
				}

				incidentEdges := graph.IncidentEdges(node1)
				Expect(incidentEdges).To(beSetThatIsNotMutable[EndpointPair[int]]())

				graph = putEdge(graph, node1, node2)
				// TODO: Pending uncommenting of newEndpointPair function and implementation of
				//       Graph.IncidentEdges()
				//Expect(incidentEdges).To(Contain(newEndpointPair(graph, node1, node2)))
			})
		})

		Context("when adding two nodes", func() {
			It("contains both nodes", func() {
				graph = addNode(graph, node1)
				graph = addNode(graph, node2)

				Expect(graph.Nodes()).To(beSetThatConsistsOf(node1, node2))
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
				graph = addNode(graph, node1)

				Expect(graphAsMutable().AddNode(node1)).To(BeFalse())
			})
		})

		Context("when removing an existing node", func() {
			var removed bool

			BeforeEach(func() {
				skipIfGraphIsNotMutable()
				graph = putEdge(graph, node1, node2)
				graph = putEdge(graph, node3, node1)

				removed = graphAsMutable().RemoveNode(node1)
			})

			It("returns true", func() {
				Expect(removed).To(BeTrue())
			})

			It("it leaves the other nodes alone", func() {
				Expect(graph.Nodes()).To(beSetThatConsistsOf(node2, node3))
			})

			It("removes its connections to its adjacent nodes", func() {
				Expect(graph.AdjacentNodes(node2)).To(beSetThatIsEmpty[int]())
				Expect(graph.AdjacentNodes(node3)).To(beSetThatIsEmpty[int]())
			})

			It("removes the connected edges", func() {
				Expect(graph.Edges()).To(beSetThatIsEmpty[EndpointPair[int]]())
			})
		})

		Context("when removing an absent node", func() {
			var removed bool

			BeforeEach(func() {
				skipIfGraphIsNotMutable()
				graph = addNode(graph, node1)

				removed = graphAsMutable().RemoveNode(nodeNotInGraph)
			})

			It("returns false", func() {
				Expect(removed).To(BeFalse())
			})

			It("leaves all the nodes alone", func() {
				Expect(graph.Nodes()).To(beSetThatConsistsOf(node1))
			})
		})

		Context("when adding one edge", func() {
			BeforeEach(func() {
				graph = putEdge(graph, node1, node2)
			})

			It("reports that both nodes are adjacent to each other", func() {
				Expect(graph.AdjacentNodes(node1)).To(beSetThatConsistsOf(node2))
				Expect(graph.AdjacentNodes(node2)).To(beSetThatConsistsOf(node1))
			})

			It("reports that both nodes have a degree of 1", func() {
				Expect(graph.Degree(node1)).To(Equal(1))
				Expect(graph.Degree(node2)).To(Equal(1))
			})
		})

		Context("when adding two connected edges", func() {
			It("reports that the common node has a degree of 2", func() {
				graph = putEdge(graph, node1, node2)
				graph = putEdge(graph, node1, node3)

				Expect(graph.Degree(node1)).To(Equal(2))
			})

			It("reports the two unique nodes as adjacent to the common one", func() {
				graph = putEdge(graph, node1, node2)
				graph = putEdge(graph, node1, node3)

				Expect(graph.AdjacentNodes(node1)).To(beSetThatConsistsOf(node2, node3))
			})
		})

		Context("when adding two anti-parallel edges", func() {
			Context("and removing one of the nodes", func() {
				It("leaves the other node alone", func() {
					skipIfGraphIsNotMutable()
					graph = putEdge(graph, node1, node2)
					graph = putEdge(graph, node2, node1)
					graphAsMutable().RemoveNode(node1)

					Expect(graph.Nodes()).To(beSetThatConsistsOf(node2))
				})

				It("removes both edges", func() {
					Expect(graph.Edges()).To(beSetThatIsEmpty[EndpointPair[int]]())
				})
			})
		})

		Context("when removing an existing edge", func() {
			var removed bool

			BeforeEach(func() {
				skipIfGraphIsNotMutable()
				graph = putEdge(graph, node1, node2)

				removed = graphAsMutable().RemoveEdge(node1, node2)
			})

			It("returns true", func() {
				Expect(removed).To(BeTrue())
			})

			It("removes the connection between its nodes", func() {
				// TODO: Pending full implementation of Graph.Successors and Graph.Predecessors
				Skip("Pending full implementation of Graph.Successors and Graph.Predecessors")

				Expect(graph.Successors(node1)).To(beSetThatIsEmpty[int]())
				Expect(graph.Predecessors(node2)).To(beSetThatIsEmpty[int]())
			})
		})

		Context("when removing an absent edge with an existing nodeU", func() {
			var removed bool

			BeforeEach(func() {
				skipIfGraphIsNotMutable()
				graph = putEdge(graph, node1, node2)

				removed = graphAsMutable().RemoveEdge(node1, nodeNotInGraph)
			})

			It("returns false", func() {
				Expect(removed).To(BeFalse())
			})

			It("leaves the existing nodes alone", func() {
				// TODO: Pending full implementation of Graph.Successors
				Skip("Pending full implementation of Graph.Successors")

				Expect(graph.Successors(node1)).To(Contain(node2))
			})
		})

		Context("when removing an absent edge with an existing nodeV", func() {
			var removed bool

			BeforeEach(func() {
				skipIfGraphIsNotMutable()
				graph = putEdge(graph, node1, node2)

				removed = graphAsMutable().RemoveEdge(nodeNotInGraph, node2)
			})

			It("returns false", func() {
				Expect(removed).To(BeFalse())
			})

			It("leaves the existing nodes alone", func() {
				// TODO: Pending full implementation of Graph.Successors
				Skip("Pending full implementation of Graph.Successors")

				Expect(graph.Successors(node1)).To(Contain(node2))
			})
		})

		Context("when finding the predecessors of an absent node", func() {
			It("returns an empty set", func() {
				Expect(graph.Predecessors(nodeNotInGraph)).
					To(beSetThatIsEmpty[int]())
			})
		})

		Context("when finding the successors of an absent node", func() {
			It("returns an empty set", func() {
				Expect(graph.Successors(nodeNotInGraph)).
					To(beSetThatIsEmpty[int]())
			})
		})

		Context("when finding the adjacent nodes of an absent node", func() {
			It("returns an empty set", func() {
				Expect(graph.AdjacentNodes(nodeNotInGraph)).
					To(beSetThatIsEmpty[int]())
			})
		})

		Context("when finding the incident edges of an absent node", func() {
			It("returns an empty set", func() {
				Expect(graph.IncidentEdges(nodeNotInGraph)).
					To(beSetThatIsEmpty[EndpointPair[int]]())
			})
		})

		Context("when finding the degree of an absent node", func() {
			It("returns zero", func() {
				Expect(graph.Degree(nodeNotInGraph)).
					To(BeZero())
			})
		})

		Context("when finding the in degree of an absent node", func() {
			It("returns zero", func() {
				Expect(graph.InDegree(nodeNotInGraph)).
					To(BeZero())
			})
		})

		Context("when finding the out degree of of an absent node", func() {
			It("returns zero", func() {
				Expect(graph.OutDegree(nodeNotInGraph)).
					To(BeZero())
			})
		})
	})

	undirectedGraphTests(createGraph, addNode, putEdge, containersMode)
}

func undirectedGraphTests(
	createGraph func() Graph[int],
	addNode func(g Graph[int], n int) Graph[int],
	putEdge func(g Graph[int], n1 int, n2 int) Graph[int],
	containersMode ContainersMode) {

	Context("given an undirected graph", func() {

		var (
			graph Graph[int]
		)

		BeforeEach(func() {
			assertContainersMode(containersMode)

			graph = createGraph()
			if graph.IsDirected() {
				Skip("graph is not undirected")
			}
		})

		Context("when adding one edge", func() {
			BeforeEach(func() {
				graph = putEdge(graph, node1, node2)
			})

			It("sees both nodes as predecessors of each other", func() {
				Expect(graph.Predecessors(node2)).To(beSetThatConsistsOf(node1))
				Expect(graph.Predecessors(node1)).To(beSetThatConsistsOf(node2))
			})

			It("sees both nodes as successors of each other", func() {
				Expect(graph.Successors(node1)).To(beSetThatConsistsOf(node2))
				Expect(graph.Successors(node2)).To(beSetThatConsistsOf(node1))
			})

			It("has an incident edge connecting the first node to the second", func() {
				incidentEdges := graph.IncidentEdges(node1)

				Expect(incidentEdges).To(
					// TODO: When EndpointPair has an Equal method, replace this assertion with a custom one
					//       that checks that the set contains only one element, where the element is "equal"
					//       according to EndpointPair.Equal. (The name "BeEquivalentTo" is already taken by
					//       Gomega. Maybe "EqualAccordingToEqualMethod"?)
					beSetThatConsistsOf(NewUnorderedEndpointPair(node1, node2)))
			})

			It("has an incident edge connecting the second node to the first", func() {
				incidentEdges := graph.IncidentEdges(node2)

				Expect(incidentEdges).To(
					// TODO: When EndpointPair has an Equal method, replace this assertion with a custom one
					//       that checks that the set contains only one element, where the element is "equal"
					//       according to EndpointPair.Equal. (The name "BeEquivalentTo" is already taken by
					//       Gomega. Maybe "EqualAccordingToEqualMethod"?)
					beSetThatConsistsOf(NewUnorderedEndpointPair(node2, node1)))
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

func validateGraphState(graph Graph[int]) {
	// TODO: Pending implementation of graph.CopyOf
	//expectStronglyEquivalent(graph, graph.CopyOf(graph))
	// TODO: Pending implementation of graph.ImmutableCopyOf
	//expectStronglyEquivalent(graph, graph.ImmutableCopyOf(graph))

	// TODO: Pending implementation of Graph.String()
	//graphString := graph.String()
	// TODO: Pending implementation of Graph.IsDirected()
	//Expect(graphString).To(ContainSubstring("isDirected: %v", graph.IsDirected()))
	// TODO: Pending implementation of Graph.AllowsSelfLoops()
	//Expect(graphString).To(ContainSubstring("allowsSelfLoops: %v", graph.AllowsSelfLoops()))

	// TODO: Pending implementation of Graph.String()
	//nodeStart := strings.Index(graphString, "nodes:")
	//edgeStart := strings.Index(graphString, "edges:")
	//nodeString := graphString[nodeStart:edgeStart] // safe because contents are ASCII

	allEndpointPairs := set.New[EndpointPair[int]]()
	_ = allEndpointPairs

	sanityCheckSet(graph.Nodes()).ForEach(func(node int) {
		// TODO: Pending implementation of Graph.String()
		//Expect(nodeString).To(ContainSubstring(strconv.Itoa(node)))

		// TODO: Pending implementation of many Graph methods
		//if graph.IsDirected() {
		//	Expect(graph.Degree(node)).To(Equal(graph.MustInDegree(node) + graph.MustOutDegree(node)))
		//	Expect(graph.Predecessors(node)).To(HaveLenOf(graph.MustInDegree(node)))
		//	Expect(graph.Successors(node)).To(HaveLenOf(graph.MustOutDegree(node)))
		//} else {
		//	nodeConnectedToSelf := must(graph.AdjacentNodes(node)).Contains(node)
		//	selfLoopCount := 0
		//	if nodeConnectedToSelf {
		//		selfLoopCount = 1
		//	}
		//	Expect(graph.Degree(node)).To(Equal(must(graph.MustAdjacentNodes(nodes)).Len() + selfLoopCount))
		//	Expect(graph.Predecessors(node)).To(BeSetEqualTo(must(graph.AdjacentNodes(nodes))))
		//	Expect(graph.Successors(node)).To(BeSetEqualTo(must(graph.AdjacentNodes(nodes))))
		//	Expect(graph.InDegree(node)).To(Equal(graph.Degree(node)))
		//	Expect(graph.OutDegree(node)).To(Equal(graph.Degree(node)))
		//}

		// TODO: Pending implementation of many Graph methods
		//sanityCheckSet(must(graph.AdjacentNodes(node))).ForEach(func(adjacentNode int) {
		//	if !graph.AllowsSelfLoops() {
		//		Expect(node).ToNot(Equal(adjacentNode))
		//	}
		//	Expect(
		//		must(graph.Predecessors(node)).Contains(adjacentNode) ||
		//			must(graph.Successors(node)).Contains(adjacentNode)).
		//		To(BeTrue())
		//})

		// TODO: Pending implementation of Graph.IsDirected() and Graph.HasEdgeConnecting()
		//sanityCheckSet(must(graph.Successors(node))).ForEach(func(successor int) {
		//	allEndpointPairs.Add(newEndpointPair(graph, node, successor))
		//	Expect(graph.Predecessors(successor)).To(Contain(node))
		//	Expect(graph.HasEdgeConnecting(node, successor)).To(BeTrue())
		//	Expect(graph.IncidentEdges(node)).To(Contain(graph, node, successor))
		//})

		// TODO: Pending implementation of Graph.IsDirected() and Graph.HasEdgeConnecting()
		//sanityCheckSet(must(graph.IncidentEdges(node))).ForEach(func(endpoints EndpointPair[int]) {
		//	if graph.IsDirected() {
		//		Expect(graph.HasEdgeConnecting(endpoints.Source(), endpoints.Target())).To(BeTrue())
		//	} else {
		//		Expect(graph.HasEdgeConnecting(endpoints.NodeU(), endpoints.NodeV())).To(BeTrue())
		//	}
		//})

		// TODO: Pending implementation of Graph.Edges()
		//sanityCheckSet(graph.Edges())
		//Expect(graph.Edges()).ToNot(Contain(newEndpointPair(graph, nodeNotInGraph, nodeNotInGraph)))
		//Expect(graph.Edges()).To(beSetThatConsistsOfElementsIn(allEndpointPairs))
	})
}

func expectStronglyEquivalent(first Graph[int], second Graph[int]) {
	// Properties not covered by Graph.Equal()
	// TODO: Pending implementation of Graph.AllowsSelfLoops()
	//Expect(first.AllowsSelfLoops()).To(Equal(second.AllowsSelfLoops()))
	// TODO: Pending implementation of Graph.NodeOrder()
	//Expect(first).To(haveNodeOrder(second.NodeOrder()))

	// TODO: Pending implementation of Graph.Equal()
	//Expect(first).To(beGraphEqualTo(second))
}

func must[N any](value N, err error) N {
	Expect(err).ToNot(HaveOccurred())
	return value
}

// TODO: Pending implementation of Graph.IsDirected()
//func newEndpointPair[N comparable](graph Graph[N], nodeU N, nodeV N) EndpointPair[N] {
//	if graph.IsDirected() {
//		return NewOrderedEndpointPair(nodeU, nodeV)
//	}
//	return NewUnorderedEndpointPair(nodeU, nodeV)
//}

// In some cases, our graph implementations return custom sets that define their own method implementations. Verify that
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
	//Expect(theSet).To(HaveStringConsistingOfElementsIn(ForEachToSlice(theSet)))
	// TODO: Pending introduction of Set.Equal()
	//Expect(theSet).To(BeSetEqualTo(set.CopyOf(theSet)))
	return set
}

func beSetThatConsistsOf[N comparable](first N, others ...N) types.GomegaMatcher {
	all := []N{first}
	all = append(all, others...)

	return WithTransform(ForEachToSlice[N], ConsistOf(all))
}

func beSetThatConsistsOfElementsIn[T comparable](set set.Set[T]) types.GomegaMatcher {
	return WithTransform(ForEachToSlice[int], ConsistOf(ForEachToSlice(set)))
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
