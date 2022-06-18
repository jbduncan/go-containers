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
	graphTests(
		func() graph.Graph[int] {
			return graph.Undirected[int]().Build()
		},
		func(graphFixture graph.Graph[int], n int) graph.Graph[int] {
			graphAsMutable, _ := graphFixture.(graph.MutableGraph[int])
			graphAsMutable.AddNode(n)

			return graphFixture
		},
		func(graphFixture graph.Graph[int], node1 int, node2 int) graph.Graph[int] {
			graphAsMutable, _ := graphFixture.(graph.MutableGraph[int])
			graphAsMutable.PutEdge(node1, node2)

			return graphFixture
		},
	)
})

const (
	node1          = 1
	node2          = 2
	node3          = 3
	nodeNotInGraph = 1_000
)

// graphTests produces a suite of Ginkgo test cases for testing implementations of Graph interface. Graph instances
// created for testing should have int nodes.
//
// Test cases that should be handled similarly in any graph implementation are included in this function. For example,
// testing that `Nodes()` method returns the set of the nodes in the graph. The following test cases are explicitly not
// tested:
//   - Test cases related to whether the graph is directed or undirected.
//   - Test cases related to the specific implementation of the Graph interface.
func graphTests(
	createGraph func() graph.Graph[int],
	addNode func(g graph.Graph[int], n int) graph.Graph[int],
	putEdge func(g graph.Graph[int], n1 int, n2 int) graph.Graph[int]) {

	Context("given a graph", func() {
		var (
			graphFixture graph.Graph[int]
		)

		graphIsMutable := func() bool {
			_, result := graphFixture.(graph.MutableGraph[int])

			return result
		}

		graphAsMutable := func() graph.MutableGraph[int] {
			result, _ := graphFixture.(graph.MutableGraph[int])

			return result
		}

		skipIfGraphIsNotMutable := func() {
			if !graphIsMutable() {
				Skip("Graph is not mutable.")
			}
		}

		BeforeEach(func() {
			graphFixture = createGraph()

			Expect(graphFixture.Nodes()).To(beSetThatIsEmpty[int]())
			// TODO: Uncomment when working on Graph.Edges() method
			// Expect(graphFixture.Edges()).To(beSetThatIsEmpty())
		})

		It("contains no nodes", func() {
			Expect(graphFixture.Nodes()).To(beSetThatIsEmpty[int]())
		})

		Context("when adding one node", func() {
			It("contains just the node", func() {
				graphFixture = addNode(graphFixture, node1)
				Expect(graphFixture.Nodes()).To(beSetThatConsistsOf(node1))
			})

			It("reports that the node has no adjacent nodes", func() {
				graphFixture = addNode(graphFixture, node1)
				Expect(graphFixture.AdjacentNodes(node1)).To(beSetThatIsEmpty[int]())
			})

			It("reports that the node has no predecessors", func() {
				graphFixture = addNode(graphFixture, node1)
				Expect(graphFixture.Predecessors(node1)).To(beSetThatIsEmpty[int]())
			})

			It("reports that the node has no successors", func() {
				graphFixture = addNode(graphFixture, node1)
				Expect(graphFixture.Successors(node1)).To(beSetThatIsEmpty[int]())
			})

			It("reports that the node has no incident edges", func() {
				graphFixture = addNode(graphFixture, node1)
				Expect(graphFixture.IncidentEdges(node1)).To(beSetThatIsEmpty[graph.EndpointPair[int]]())
			})

			It("reports that the node has a degree of 0", func() {
				graphFixture = addNode(graphFixture, node1)
				Expect(graphFixture.Degree(node1)).To(BeZero())
			})

			It("reports that the node has an in degree of 0", func() {
				graphFixture = addNode(graphFixture, node1)
				Expect(graphFixture.InDegree(node1)).To(BeZero())
			})

			It("reports that the node has an out degree of 0", func() {
				graphFixture = addNode(graphFixture, node1)
				Expect(graphFixture.OutDegree(node1)).To(BeZero())
			})
		})

		Context("when adding two nodes", func() {
			It("contains both nodes", func() {
				graphFixture = addNode(graphFixture, node1)
				graphFixture = addNode(graphFixture, node2)
				Expect(graphFixture.Nodes()).To(beSetThatConsistsOf(node1, node2))
			})
		})

		Context("when adding a new node", func() {
			It("returns true", func() {
				skipIfGraphIsNotMutable()

				Expect(graphAsMutable().AddNode(node1)).To(BeTrue())
				Expect(graphFixture.Nodes()).To(beSetThatConsistsOf(node1))
			})
		})

		Context("when adding an existing node", func() {
			It("returns false", func() {
				skipIfGraphIsNotMutable()

				graphFixture = addNode(graphFixture, node1)
				Expect(graphAsMutable().AddNode(node1)).To(BeFalse())
			})

			It("does not add the node again", func() {
				graphFixture = addNode(graphFixture, node1)
				graphFixture = addNode(graphFixture, node1)
				Expect(graphFixture.Nodes()).To(beSetThatConsistsOf(node1))
			})
		})

		Context("when removing an existing node", func() {
			It("returns true", func() {
				skipIfGraphIsNotMutable()

				graphFixture = addNode(graphFixture, node1)
				Expect(graphAsMutable().RemoveNode(node1)).To(BeTrue())
			})

			It("removes the node", func() {
				skipIfGraphIsNotMutable()

				graphFixture = addNode(graphFixture, node1)
				graphAsMutable().RemoveNode(node1)
				Expect(graphFixture.Nodes()).To(beSetThatIsEmpty[int]())
			})

			It("removes the connections to its adjacent nodes", func() {
				skipIfGraphIsNotMutable()

				graphFixture = putEdge(graphFixture, node1, node2)
				graphFixture = putEdge(graphFixture, node3, node1)

				graphAsMutable().RemoveNode(node1)

				Expect(graphFixture.AdjacentNodes(node2)).To(beSetThatIsEmpty[int]())
				Expect(graphFixture.AdjacentNodes(node3)).To(beSetThatIsEmpty[int]())
			})

			// TODO: Pending implementation of Graph.Edges()
			//It("removes the connected edges", func() {
			//	skipIfGraphIsNotMutable()
			//
			//	graphFixture = putEdge(graphFixture, node1, node2)
			//	graphFixture = putEdge(graphFixture, node3, node1)
			//
			//	graphAsMutable().RemoveNode(node1)
			//
			//	Expect(graphFixture.Edges()).To(beSetThatIsEmpty[int]())
			//})

			Context("and querying the node after removal", func() {
				It("returns an error", func() {
					skipIfGraphIsNotMutable()

					graphFixture = addNode(graphFixture, node1)

					graphAsMutable().RemoveNode(node1)

					Expect(graphFixture.AdjacentNodes(node1)).
						Error().
						To(MatchError(fmt.Sprintf("%d: node not an element of this graph", node1)))
				})
			})
		})

		Context("when removing one node from two anti-parallel edges", func() {
			It("leaves the other node alone", func() {
				skipIfGraphIsNotMutable()

				graphFixture = putEdge(graphFixture, node1, node2)
				graphFixture = putEdge(graphFixture, node2, node1)

				graphAsMutable().RemoveNode(node1)

				Expect(graphFixture.Nodes()).To(beSetThatConsistsOf(node2))
			})

			// TODO: Pending implementation of Graph.Edges()
			//It("removes both edges", func() {
			//	skipIfGraphIsNotMutable()
			//
			//	graphFixture = putEdge(graphFixture, node1, node2)
			//	graphFixture = putEdge(graphFixture, node2, node1)
			//
			//	graphAsMutable().RemoveNode(node1)
			//
			//	Expect(graphFixture.Edges()).To(beSetThatIsEmpty[int]())
			//})
		})

		Context("when removing an absent node", func() {
			It("returns false", func() {
				skipIfGraphIsNotMutable()

				Expect(graphAsMutable().RemoveNode(nodeNotInGraph)).To(BeFalse())
			})

			It("leaves the existing nodes alone", func() {
				skipIfGraphIsNotMutable()

				graphFixture = addNode(graphFixture, node1)

				graphAsMutable().RemoveNode(nodeNotInGraph)
				Expect(graphFixture.Nodes()).To(beSetThatConsistsOf(node1))
			})
		})

		Context("when adding one edge", func() {
			It("reports both nodes as being adjacent", func() {
				graphFixture = putEdge(graphFixture, node1, node2)
				Expect(graphFixture.AdjacentNodes(node1)).To(beSetThatConsistsOf(node2))
				Expect(graphFixture.AdjacentNodes(node2)).To(beSetThatConsistsOf(node1))
			})

			It("reports both nodes as having a degree of 1", func() {
				graphFixture = putEdge(graphFixture, node1, node2)
				Expect(graphFixture.Degree(node1)).To(Equal(1))
				Expect(graphFixture.Degree(node2)).To(Equal(1))
			})
		})

		Context("when adding two connected edges", func() {
			It("reports that the common node has a degree of 2", func() {
				graphFixture = putEdge(graphFixture, node1, node2)
				graphFixture = putEdge(graphFixture, node1, node3)
				Expect(graphFixture.Degree(node1)).To(Equal(2))
			})
		})

		Context("when finding the predecessors of an absent node", func() {
			It("returns an error", func() {
				Expect(graphFixture.Predecessors(nodeNotInGraph)).
					Error().
					To(MatchError(fmt.Sprintf("%d: node not an element of this graph", nodeNotInGraph)))
			})
		})

		Context("when finding the successors of an absent node", func() {
			It("returns an error", func() {
				Expect(graphFixture.Successors(nodeNotInGraph)).
					Error().
					To(MatchError(fmt.Sprintf("%d: node not an element of this graph", nodeNotInGraph)))
			})
		})

		Context("when finding the adjacent nodes of an absent node", func() {
			It("returns an error", func() {
				Expect(graphFixture.AdjacentNodes(nodeNotInGraph)).
					Error().
					To(MatchError(fmt.Sprintf("%d: node not an element of this graph", nodeNotInGraph)))
			})
		})

		Context("when finding the incident edges of an absent node", func() {
			It("returns an error", func() {
				Expect(graphFixture.IncidentEdges(nodeNotInGraph)).
					Error().
					To(MatchError(fmt.Sprintf("%d: node not an element of this graph", nodeNotInGraph)))
			})
		})

		Context("when finding the degree of an absent node", func() {
			It("returns an error", func() {
				Expect(graphFixture.Degree(nodeNotInGraph)).
					Error().
					To(MatchError(fmt.Sprintf("%d: node not an element of this graph", nodeNotInGraph)))
			})
		})

		Context("when finding the in degree of an absent node", func() {
			It("returns an error", func() {
				Expect(graphFixture.InDegree(nodeNotInGraph)).
					Error().
					To(MatchError(fmt.Sprintf("%d: node not an element of this graph", nodeNotInGraph)))
			})
		})

		Context("when finding the out degree of of an absent node", func() {
			It("returns an error", func() {
				Expect(graphFixture.OutDegree(nodeNotInGraph)).
					Error().
					To(MatchError(fmt.Sprintf("%d: node not an element of this graph", nodeNotInGraph)))
			})
		})

		Context("when checking the mutability of the Nodes set", func() {
			It("is not mutable", func() {
				Expect(graphFixture.Nodes()).To(beSetThatIsNotMutable[int]())
			})
		})

		Context("when checking the mutability of the AdjacentNodes set", func() {
			It("is not mutable", func() {
				graphFixture = addNode(graphFixture, node1)
				Expect(graphFixture.AdjacentNodes(node1)).To(beSetThatIsNotMutable[int]())
			})
		})

		Context("when checking the mutability of the Predecessors set", func() {
			It("is not mutable", func() {
				graphFixture = addNode(graphFixture, node1)
				Expect(graphFixture.Predecessors(node1)).To(beSetThatIsNotMutable[int]())
			})
		})

		Context("when checking the mutability of the Successors set", func() {
			It("is not mutable", func() {
				graphFixture = addNode(graphFixture, node1)
				Expect(graphFixture.Successors(node1)).To(beSetThatIsNotMutable[int]())
			})
		})

		Context("when checking the mutability of the IncidentEdges set", func() {
			It("is not mutable", func() {
				graphFixture = addNode(graphFixture, node1)
				Expect(graphFixture.IncidentEdges(node1)).To(beSetThatIsNotMutable[graph.EndpointPair[int]]())
			})
		})
	})
}

func beSetThatConsistsOf(first int, others ...int) types.GomegaMatcher {
	all := []int{first}
	all = append(all, others...)

	return WithTransform(ForEachToSlice[int], ConsistOf(all))
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

func forEachCount(set set.Set[int]) int {
	var result int

	set.ForEach(func(elem int) {
		result++
	})

	return result
}
