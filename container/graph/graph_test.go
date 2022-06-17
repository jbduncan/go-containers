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
		func(g graph.Graph[int], n int) graph.Graph[int] {
			g.(graph.MutableGraph[int]).AddNode(n)
			return g
		},
		func(g graph.Graph[int], n1 int, n2 int) graph.Graph[int] {
			g.(graph.MutableGraph[int]).PutEdge(n1, n2)
			return g
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
			n1             = 1
			n2             = 2
			n3             = 3
			nodeNotInGraph = 1_000
		)

		var (
			g          graph.Graph[int]
			gAsMutable graph.MutableGraph[int]
			gIsMutable bool
		)

		BeforeEach(func() {
			g = createGraph()
			gAsMutable, gIsMutable = g.(graph.MutableGraph[int])

			Expect(g.Nodes()).To(beSetThatIsEmpty[int]())
			// TODO: Uncomment when working on Graph.Edges() method
			// Expect(g.Edges()).To(beSetThatIsEmpty())
		})

		It("contains no nodes", func() {
			Expect(g.Nodes()).To(beSetThatIsEmpty[int]())
		})

		Context("when adding one node", func() {
			It("contains just the node", func() {
				g = addNode(g, n1)
				Expect(g.Nodes()).To(beSetThatConsistsOf(n1))
			})

			It("reports that the node has no adjacent nodes", func() {
				g = addNode(g, n1)
				Expect(g.AdjacentNodes(n1)).To(beSetThatIsEmpty[int]())
			})

			It("reports that the node has no predecessors", func() {
				g = addNode(g, n1)
				Expect(g.Predecessors(n1)).To(beSetThatIsEmpty[int]())
			})

			It("reports that the node has no successors", func() {
				g = addNode(g, n1)
				Expect(g.Successors(n1)).To(beSetThatIsEmpty[int]())
			})

			It("reports that the node has no incident edges", func() {
				g = addNode(g, n1)
				Expect(g.IncidentEdges(n1)).To(beSetThatIsEmpty[graph.EndpointPair[int]]())
			})

			It("reports that the node has a degree of 0", func() {
				g = addNode(g, n1)
				Expect(g.Degree(n1)).To(BeZero())
			})

			It("reports that the node has an in degree of 0", func() {
				g = addNode(g, n1)
				Expect(g.InDegree(n1)).To(BeZero())
			})
		})

		Context("when adding two nodes", func() {
			It("contains both nodes", func() {
				g = addNode(g, n1)
				g = addNode(g, n2)
				Expect(g.Nodes()).To(beSetThatConsistsOf(n1, n2))
			})
		})

		Context("when adding a new node", func() {
			It("returns true", func() {
				if !gIsMutable {
					Skip("Graph is not mutable.")
				}

				Expect(gAsMutable.AddNode(n1)).To(BeTrue())
				Expect(g.Nodes()).To(beSetThatConsistsOf(n1))
			})
		})

		Context("when adding an existing node", func() {
			It("returns false", func() {
				if !gIsMutable {
					Skip("Graph is not mutable.")
				}

				g = addNode(g, n1)
				Expect(gAsMutable.AddNode(n1)).To(BeFalse())
				Expect(g.Nodes()).To(beSetThatConsistsOf(n1))
			})
		})

		Context("when adding one edge", func() {
			It("reports both nodes as being adjacent", func() {
				g = putEdge(g, n1, n2)
				Expect(g.AdjacentNodes(n1)).To(beSetThatConsistsOf(n2))
				Expect(g.AdjacentNodes(n2)).To(beSetThatConsistsOf(n1))
			})

			It("reports both nodes as having a degree of 1", func() {
				g = putEdge(g, n1, n2)
				Expect(g.Degree(n1)).To(Equal(1))
				Expect(g.Degree(n2)).To(Equal(1))
			})
		})

		Context("when adding two connected edges", func() {
			It("reports that the common node has a degree of 2", func() {
				g = putEdge(g, n1, n2)
				g = putEdge(g, n1, n3)
				Expect(g.Degree(n1)).To(Equal(2))
			})
		})

		Context("when finding the predecessors of non-existent node", func() {
			It("returns an error", func() {
				Expect(g.Predecessors(nodeNotInGraph)).
					Error().
					To(MatchError(fmt.Sprintf("node %d not an element of this graph", nodeNotInGraph)))
			})
		})

		Context("when finding the successors of non-existent node", func() {
			It("returns an error", func() {
				Expect(g.Successors(nodeNotInGraph)).
					Error().
					To(MatchError(fmt.Sprintf("node %d not an element of this graph", nodeNotInGraph)))
			})
		})

		Context("when finding the adjacent nodes of non-existent node", func() {
			It("returns an error", func() {
				Expect(g.AdjacentNodes(nodeNotInGraph)).
					Error().
					To(MatchError(fmt.Sprintf("node %d not an element of this graph", nodeNotInGraph)))
			})
		})

		Context("when finding the incident edges of non-existent node", func() {
			It("returns an error", func() {
				Expect(g.IncidentEdges(nodeNotInGraph)).
					Error().
					To(MatchError(fmt.Sprintf("node %d not an element of this graph", nodeNotInGraph)))
			})
		})

		Context("when finding the degree of non-existent node", func() {
			It("returns an error", func() {
				Expect(g.Degree(nodeNotInGraph)).
					Error().
					To(MatchError(fmt.Sprintf("node %d not an element of this graph", nodeNotInGraph)))
			})
		})

		Context("when finding the in degree of a non-existent node", func() {
			It("returns an error", func() {
				Expect(g.InDegree(nodeNotInGraph)).
					Error().
					To(MatchError(fmt.Sprintf("node %d not an element of this graph", nodeNotInGraph)))
			})
		})

		Context("when checking the mutability of the Nodes set", func() {
			It("is not mutable", func() {
				Expect(g.Nodes()).To(beSetThatIsNotMutable[int]())
			})
		})

		Context("when checking the mutability of the AdjacentNodes set", func() {
			It("is not mutable", func() {
				g = addNode(g, n1)
				Expect(g.AdjacentNodes(n1)).To(beSetThatIsNotMutable[int]())
			})
		})

		Context("when checking the mutability of the Predecessors set", func() {
			It("is not mutable", func() {
				g = addNode(g, n1)
				Expect(g.Predecessors(n1)).To(beSetThatIsNotMutable[int]())
			})
		})

		Context("when checking the mutability of the Successors set", func() {
			It("is not mutable", func() {
				g = addNode(g, n1)
				Expect(g.Successors(n1)).To(beSetThatIsNotMutable[int]())
			})
		})

		Context("when checking the mutability of the IncidentEdges set", func() {
			It("is not mutable", func() {
				g = addNode(g, n1)
				Expect(g.IncidentEdges(n1)).To(beSetThatIsNotMutable[graph.EndpointPair[int]]())
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
