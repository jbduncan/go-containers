package graph_test

import (
	"fmt"

	"github.com/onsi/gomega/types"
	"go-containers/container/graph"
	"go-containers/container/set"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
		const (
			node1          = 1
			node2          = 2
			node3          = 3
			nodeNotInGraph = 1_000
		)

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
				if !graphIsMutable() {
					Skip("Graph is not mutable.")
				}

				Expect(graphAsMutable().AddNode(node1)).To(BeTrue())
				Expect(graphFixture.Nodes()).To(beSetThatConsistsOf(node1))
			})
		})

		Context("when adding an existing node", func() {
			It("returns false", func() {
				if !graphIsMutable() {
					Skip("Graph is not mutable.")
				}

				graphFixture = addNode(graphFixture, node1)
				Expect(graphAsMutable().AddNode(node1)).To(BeFalse())
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

		Context("when finding the predecessors of non-existent node", func() {
			It("returns an error", func() {
				Expect(graphFixture.Predecessors(nodeNotInGraph)).
					Error().
					To(MatchError(fmt.Sprintf("%d: node not an element of this graph", nodeNotInGraph)))
			})
		})

		Context("when finding the successors of non-existent node", func() {
			It("returns an error", func() {
				Expect(graphFixture.Successors(nodeNotInGraph)).
					Error().
					To(MatchError(fmt.Sprintf("%d: node not an element of this graph", nodeNotInGraph)))
			})
		})

		Context("when finding the adjacent nodes of non-existent node", func() {
			It("returns an error", func() {
				Expect(graphFixture.AdjacentNodes(nodeNotInGraph)).
					Error().
					To(MatchError(fmt.Sprintf("%d: node not an element of this graph", nodeNotInGraph)))
			})
		})

		Context("when finding the incident edges of non-existent node", func() {
			It("returns an error", func() {
				Expect(graphFixture.IncidentEdges(nodeNotInGraph)).
					Error().
					To(MatchError(fmt.Sprintf("%d: node not an element of this graph", nodeNotInGraph)))
			})
		})

		Context("when finding the degree of non-existent node", func() {
			It("returns an error", func() {
				Expect(graphFixture.Degree(nodeNotInGraph)).
					Error().
					To(MatchError(fmt.Sprintf("%d: node not an element of this graph", nodeNotInGraph)))
			})
		})

		Context("when finding the in degree of a non-existent node", func() {
			It("returns an error", func() {
				Expect(graphFixture.InDegree(nodeNotInGraph)).
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

	return WithTransform(toSlice[int], ConsistOf(all))
}

func beSetThatIsEmpty[T comparable]() types.GomegaMatcher {
	return WithTransform(toSlice[T], BeEmpty())
}

func beSetThatIsNotMutable[T comparable]() types.GomegaMatcher {
	return WithTransform(
		func(s set.Set[T]) bool {
			_, mutable := s.(set.MutableSet[T])
			return mutable
		},
		BeFalse())
}

func toSlice[T comparable](s set.Set[T]) []T {
	var result []T

	s.ForEach(func(elem T) {
		result = append(result, elem)
	})

	return result
}
