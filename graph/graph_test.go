package graph_test

import (
	"fmt"
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
	graphTests(graphName, func() graph.Graph[int] { return createGraph() }, containersMode)
}

func graphTests(
	graphName string,
	createGraph func() graph.Graph[int],
	containersMode ContainersMode,
) {
	Context(fmt.Sprintf("%s: given a graph", graphName), func() {
		var grph graph.Graph[int]

		graphAsMutable := func() graph.MutableGraph[int] {
			// This function is only ever called after skipIfGraphIsNotMutable
			//nolint:revive
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

			It("has both nodes", func() {
				testSet(grph.Nodes(), node1, node2)
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
				skipIfGraphIsNotMutable()
				grph = addNode(grph, node1)

				removed = graphAsMutable().RemoveNode(nodeNotInGraph)
			})

			It("returns false", func() {
				Expect(removed).To(BeFalse())
			})

			It("leaves all the nodes alone", func() {
				testSet(grph.Nodes(), node1)
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

			It("reports that both nodes have a degree of 1", func() {
				Expect(grph.Degree(node1)).To(Equal(1))
				Expect(grph.Degree(node2)).To(Equal(1))
			})

			It("has just that edge", func() {
				// Set.Len()
				Expect(grph.Edges()).To(HaveLenOf(1))

				// Set.ForEach()
				// Uses boolean assertions to avoid unreadable error messages
				// from this nested matcher
				matcher := HaveForEachThatConsistsOf[graph.EndpointPair[int]](
					BeEquivalentToUsingEqualMethod(
						newEndpointPair(grph, node1, node2)))
				Expect(matcher.Match(grph.Edges())).To(
					BeTrue(),
					"to consist of %v according to graph.EndpointPair.Equal()",
					newEndpointPair(grph, node1, node2))

				// Set.Contains()
				Expect(grph.Edges()).To(
					Contain(newEndpointPair(grph, node1, node2)))
				Expect(grph.Edges()).ToNot(
					Contain(newEndpointPairWithOtherOrder(grph, node1, node2)))
				Expect(grph.Edges()).ToNot(
					Contain(
						graph.NewOrderedEndpointPair(
							nodeNotInGraph, nodeNotInGraph)))
				Expect(grph.Edges()).ToNot(
					Contain(
						graph.NewUnorderedEndpointPair(
							nodeNotInGraph, nodeNotInGraph)))
			})

			It("has an incident edge connecting the first node to the second node", func() {
				// Set.Len()
				Expect(grph.IncidentEdges(node1)).To(HaveLenOf(1))

				// Set.ForEach()
				// Uses boolean assertions to avoid unreadable error messages
				// from this nested matcher
				matcher := HaveForEachThatConsistsOf[graph.EndpointPair[int]](
					BeEquivalentToUsingEqualMethod(
						newEndpointPair(grph, node1, node2)))
				Expect(matcher.Match(grph.IncidentEdges(node1))).To(
					BeTrue(),
					"to consist of %v according to graph.EndpointPair.Equal()",
					newEndpointPair(grph, node1, node2))
				matcher = HaveForEachThatConsistsOf[graph.EndpointPair[int]](
					BeEquivalentToUsingEqualMethod(
						newEndpointPair(grph, node2, node1)))
				Expect(matcher.Match(grph.IncidentEdges(node1))).To(
					BeTrue(),
					"to consist of %v according to graph.EndpointPair.Equal()",
					newEndpointPair(grph, node2, node1))

				// Set.Contains()
				Expect(grph.IncidentEdges(node1)).To(
					Contain(newEndpointPair(grph, node1, node2)))
				Expect(grph.IncidentEdges(node1)).To(
					Contain(newEndpointPair(grph, node2, node1)))
				Expect(grph.IncidentEdges(node1)).ToNot(
					Contain(
						graph.NewOrderedEndpointPair(
							nodeNotInGraph, nodeNotInGraph)))
				Expect(grph.IncidentEdges(node1)).ToNot(
					Contain(
						graph.NewUnorderedEndpointPair(
							nodeNotInGraph, nodeNotInGraph)))
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

			It("has both edges", func() {
				// Set.Len()
				Expect(grph.Edges()).To(HaveLenOf(2))

				// Set.ForEach()
				// Uses boolean assertions to avoid unreadable error messages
				// from this nested matcher
				matcher := HaveForEachThatConsistsOf[graph.EndpointPair[int]](
					BeEquivalentToUsingEqualMethod(
						newEndpointPair(grph, node1, node2)),
					BeEquivalentToUsingEqualMethod(
						newEndpointPair(grph, node1, node3)))
				Expect(matcher.Match(grph.Edges())).To(
					BeTrue(),
					"to consist of %v according to graph.EndpointPair.Equal()",
					[]graph.EndpointPair[int]{
						newEndpointPair(grph, node1, node2),
						newEndpointPair(grph, node1, node3),
					})

				// Set.Contains()
				Expect(grph.Edges()).To(
					Contain(newEndpointPair(grph, node1, node2)))
				Expect(grph.Edges()).ToNot(
					Contain(newEndpointPairWithOtherOrder(grph, node1, node2)))
				Expect(grph.Edges()).To(
					Contain(newEndpointPair(grph, node1, node3)))
				Expect(grph.Edges()).ToNot(
					Contain(newEndpointPairWithOtherOrder(grph, node1, node3)))
				Expect(grph.Edges()).ToNot(
					Contain(
						graph.NewOrderedEndpointPair(
							nodeNotInGraph, nodeNotInGraph)))
				Expect(grph.Edges()).ToNot(
					Contain(
						graph.NewUnorderedEndpointPair(
							nodeNotInGraph, nodeNotInGraph)))
			})

			It("has two incident edges connected to the common node", func() {
				// Set.Len()
				Expect(grph.IncidentEdges(node1)).To(HaveLenOf(2))

				// Set.ForEach()
				// Uses boolean assertions to avoid unreadable error messages
				// from this nested matcher
				matcher := HaveForEachThatConsistsOf[graph.EndpointPair[int]](
					BeEquivalentToUsingEqualMethod(
						newEndpointPair(grph, node1, node2)),
					BeEquivalentToUsingEqualMethod(
						newEndpointPair(grph, node1, node3)))
				Expect(matcher.Match(grph.IncidentEdges(node1))).To(
					BeTrue(),
					"to consist of %v according to graph.EndpointPair.Equal()",
					[]graph.EndpointPair[int]{
						newEndpointPair(grph, node1, node2),
						newEndpointPair(grph, node1, node3),
					})

				// Set.Contains()
				Expect(grph.IncidentEdges(node1)).To(
					Contain(newEndpointPair(grph, node1, node2)))
				Expect(grph.IncidentEdges(node1)).ToNot(
					Contain(newEndpointPairWithOtherOrder(grph, node1, node2)))
				Expect(grph.IncidentEdges(node1)).To(
					Contain(newEndpointPair(grph, node1, node3)))
				Expect(grph.IncidentEdges(node1)).ToNot(
					Contain(newEndpointPairWithOtherOrder(grph, node1, node3)))
				Expect(grph.IncidentEdges(node1)).ToNot(
					Contain(
						graph.NewOrderedEndpointPair(
							nodeNotInGraph, nodeNotInGraph)))
				Expect(grph.IncidentEdges(node1)).ToNot(
					Contain(
						graph.NewUnorderedEndpointPair(
							nodeNotInGraph, nodeNotInGraph)))
			})
		})

		Context("when putting two anti-parallel edges", func() {
			Context("and removing one of the nodes", func() {
				It("leaves the other node alone", func() {
					skipIfGraphIsNotMutable()
					grph = putEdge(grph, node1, node2)
					grph = putEdge(grph, node2, node1)
					graphAsMutable().RemoveNode(node1)

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
				skipIfGraphIsNotMutable()
				grph = putEdge(grph, node1, node2)
				grph = putEdge(grph, node1, node3)

				removed = graphAsMutable().RemoveEdge(node1, node2)
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
				skipIfGraphIsNotMutable()
				grph = putEdge(grph, node1, node2)

				removed = graphAsMutable().RemoveEdge(node1, nodeNotInGraph)
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
				skipIfGraphIsNotMutable()
				grph = putEdge(grph, node1, node2)

				removed = graphAsMutable().RemoveEdge(nodeNotInGraph, node2)
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
				skipIfGraphIsNotMutable()
				grph = addNode(grph, node1)
				grph = addNode(grph, node2)

				removed = graphAsMutable().RemoveEdge(node1, node2)
			})

			It("returns false", func() {
				Expect(removed).To(BeFalse())
			})

			It("leaves the existing nodes alone", func() {
				testSet(grph.Nodes(), node1, node2)
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

	undirectedGraphTests(graphName, createGraph, containersMode)
}

func undirectedGraphTests(
	graphName string,
	createGraph func() graph.Graph[int],
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

		Context("when putting two connected edges", func() {
			BeforeEach(func() {
				grph = putEdge(grph, node1, node2)
				grph = putEdge(grph, node1, node3)
			})

			It("has both edges", func() {
				// TODO: Write equivalent test for directed graphs

				// Set.ForEach()
				// Uses boolean assertions to avoid unreadable error messages
				// from this nested matcher
				matcher := HaveForEachThatConsistsOf[graph.EndpointPair[int]](
					BeEquivalentToUsingEqualMethod(
						newEndpointPair(grph, node2, node1)),
					BeEquivalentToUsingEqualMethod(
						newEndpointPair(grph, node3, node1)))
				Expect(matcher.Match(grph.Edges())).To(
					BeTrue(),
					"to consist of %v according to graph.EndpointPair.Equal()",
					[]graph.EndpointPair[int]{
						newEndpointPair(grph, node2, node1),
						newEndpointPair(grph, node3, node1),
					})

				// Set.Contains()
				Expect(grph.Edges()).To(
					Contain(newEndpointPair(grph, node2, node1)))
				Expect(grph.Edges()).ToNot(
					Contain(newEndpointPairWithOtherOrder(grph, node2, node1)))
				Expect(grph.Edges()).To(
					Contain(newEndpointPair(grph, node3, node1)))
				Expect(grph.Edges()).ToNot(
					Contain(newEndpointPairWithOtherOrder(grph, node3, node1)))
				Expect(grph.Edges()).ToNot(
					Contain(
						graph.NewOrderedEndpointPair(
							nodeNotInGraph, nodeNotInGraph)))
				Expect(grph.Edges()).ToNot(
					Contain(
						graph.NewUnorderedEndpointPair(
							nodeNotInGraph, nodeNotInGraph)))

				// Set.String()
				Expect(grph.Edges()).To(
					HaveStringReprThatIsAnyOf(
						"[[1, 2], [1, 3]]",
						"[[1, 2], [3, 1]]",
						"[[2, 1], [1, 3]]",
						"[[2, 1], [3, 1]]",
						"[[1, 3], [1, 2]]",
						"[[1, 3], [2, 1]]",
						"[[3, 1], [1, 2]]",
						"[[3, 1], [2, 1]]"))
			})

			It("has two incident edges connected to the common node", func() {
				// TODO: Write equivalent test for directed graphs

				// Set.ForEach()
				// Uses boolean assertions to avoid unreadable error messages
				// from this nested matcher
				matcher := HaveForEachThatConsistsOf[graph.EndpointPair[int]](
					BeEquivalentToUsingEqualMethod(
						newEndpointPair(grph, node2, node1)),
					BeEquivalentToUsingEqualMethod(
						newEndpointPair(grph, node3, node1)))
				Expect(matcher.Match(grph.IncidentEdges(node1))).To(
					BeTrue(),
					"to consist of %v according to graph.EndpointPair.Equal()",
					[]graph.EndpointPair[int]{
						newEndpointPair(grph, node2, node1),
						newEndpointPair(grph, node3, node1),
					})

				// Set.Contains()
				Expect(grph.IncidentEdges(node1)).To(
					Contain(newEndpointPair(grph, node2, node1)))
				Expect(grph.IncidentEdges(node1)).ToNot(
					Contain(newEndpointPairWithOtherOrder(grph, node2, node1)))
				Expect(grph.IncidentEdges(node1)).To(
					Contain(newEndpointPair(grph, node3, node1)))
				Expect(grph.IncidentEdges(node1)).ToNot(
					Contain(newEndpointPairWithOtherOrder(grph, node3, node1)))
				Expect(grph.IncidentEdges(node1)).ToNot(
					Contain(
						graph.NewOrderedEndpointPair(
							nodeNotInGraph, nodeNotInGraph)))
				Expect(grph.IncidentEdges(node1)).ToNot(
					Contain(
						graph.NewUnorderedEndpointPair(
							nodeNotInGraph, nodeNotInGraph)))

				// Set.String()
				Expect(grph.IncidentEdges(node1)).To(
					HaveStringReprThatIsAnyOf(
						"[[1, 2], [1, 3]]",
						"[[1, 2], [3, 1]]",
						"[[2, 1], [1, 3]]",
						"[[2, 1], [3, 1]]",
						"[[1, 3], [1, 2]]",
						"[[1, 3], [2, 1]]",
						"[[3, 1], [1, 2]]",
						"[[3, 1], [2, 1]]"))
			})
		})

		Context("when putting one edge", func() {
			BeforeEach(func() {
				grph = putEdge(grph, node1, node2)
			})

			It("has just that undirected edge", func() {
				// TODO: Write equivalent test for directed graphs

				// Set.ForEach()
				// Uses boolean assertions to avoid unreadable error messages
				// from this nested matcher
				matcher := HaveForEachThatConsistsOf[graph.EndpointPair[int]](
					BeEquivalentToUsingEqualMethod(
						newEndpointPair(grph, node2, node1)))
				Expect(matcher.Match(grph.Edges())).To(
					BeTrue(),
					"to consist of %v according to graph.EndpointPair.Equal()",
					newEndpointPair(grph, node2, node1))

				// Set.Contains()
				Expect(grph.Edges()).To(
					Contain(newEndpointPair(grph, node2, node1)))
				Expect(grph.Edges()).ToNot(
					Contain(newEndpointPairWithOtherOrder(grph, node2, node1)))
				Expect(grph.Edges()).ToNot(
					Contain(
						graph.NewOrderedEndpointPair(
							nodeNotInGraph, nodeNotInGraph)))
				Expect(grph.Edges()).ToNot(
					Contain(
						graph.NewUnorderedEndpointPair(
							nodeNotInGraph, nodeNotInGraph)))

				// Set.String()
				Expect(grph.Edges()).To(
					HaveStringReprThatIsAnyOf("[[1, 2]]", "[[2, 1]]"))
			})

			It("sees both nodes as predecessors of each other", func() {
				testSet(grph.Predecessors(node2), node1)
				testSet(grph.Predecessors(node1), node2)
			})

			It("sees both nodes as successors of each other", func() {
				testSet(grph.Successors(node1), node2)
				testSet(grph.Successors(node2), node1)
			})

			It("has an incident edge connecting the second node to the first node", func() {
				// Set.Len()
				Expect(grph.IncidentEdges(node2)).To(HaveLenOf(1))

				// Set.ForEach()
				// Uses boolean assertions to avoid unreadable error messages
				// from this nested matcher
				matcher := HaveForEachThatConsistsOf[graph.EndpointPair[int]](
					BeEquivalentToUsingEqualMethod(
						newEndpointPair(grph, node1, node2)))
				Expect(matcher.Match(grph.IncidentEdges(node2))).To(
					BeTrue(),
					"to consist of %v according to graph.EndpointPair.Equal()",
					newEndpointPair(grph, node1, node2))
				matcher = HaveForEachThatConsistsOf[graph.EndpointPair[int]](
					BeEquivalentToUsingEqualMethod(
						newEndpointPair(grph, node2, node1)))
				Expect(matcher.Match(grph.IncidentEdges(node2))).To(
					BeTrue(),
					"to consist of %v according to graph.EndpointPair.Equal()",
					newEndpointPair(grph, node2, node1))

				// Set.Contains()
				Expect(grph.IncidentEdges(node2)).To(
					Contain(newEndpointPair(grph, node1, node2)))
				Expect(grph.IncidentEdges(node2)).ToNot(
					Contain(newEndpointPairWithOtherOrder(grph, node1, node2)))
				Expect(grph.IncidentEdges(node2)).To(
					Contain(newEndpointPair(grph, node2, node1)))
				Expect(grph.IncidentEdges(node2)).ToNot(
					Contain(newEndpointPairWithOtherOrder(grph, node2, node1)))
				Expect(grph.IncidentEdges(node2)).ToNot(
					Contain(
						graph.NewOrderedEndpointPair(
							nodeNotInGraph, nodeNotInGraph)))
				Expect(grph.IncidentEdges(node2)).ToNot(
					Contain(
						graph.NewUnorderedEndpointPair(
							nodeNotInGraph, nodeNotInGraph)))

				// Set.String()
				Expect(grph.IncidentEdges(node2)).To(
					HaveStringReprThatIsAnyOf("[[1, 2]]", "[[2, 1]]"))
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

					testSet(grph.AdjacentNodes(node1), node1)
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

func newEndpointPair[N comparable](g graph.Graph[N], nodeU N, nodeV N) graph.EndpointPair[N] {
	if g.IsDirected() {
		return graph.NewOrderedEndpointPair(nodeU, nodeV)
	}
	return graph.NewUnorderedEndpointPair(nodeU, nodeV)
}

func newEndpointPairWithOtherOrder[N comparable](g graph.Graph[N], nodeU N, nodeV N) graph.EndpointPair[N] {
	if g.IsDirected() {
		return graph.NewUnorderedEndpointPair(nodeU, nodeV)
	}
	return graph.NewOrderedEndpointPair(nodeU, nodeV)
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
		expectedValueStrs = append(expectedValueStrs, fmt.Sprintf("%v", v))
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
