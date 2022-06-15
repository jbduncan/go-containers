package graph_test

import (
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
func graphTests(createGraph func() graph.Graph[int], addNode func(g graph.Graph[int], n int) graph.Graph[int]) {
	Context("given a graph", func() {
		const (
			n1 = 1
			n2 = 2
		)

		var (
			g          graph.Graph[int]
			gAsMutable graph.MutableGraph[int]
			gIsMutable bool
		)

		BeforeEach(func() {
			g = createGraph()
			gAsMutable, gIsMutable = g.(graph.MutableGraph[int])

			Expect(g.Nodes()).To(beSetThatIsEmpty())
			// TODO: Uncomment when working on Graph.Edges() method
			// Expect(g.Edges()).To(beSetThatIsEmpty())
		})

		Context("when adding one node", func() {
			It("contains just the node", func() {
				g = addNode(g, n1)
				Expect(g.Nodes()).To(beSetThatConsistsOf(n1))
			})
		})

		Context("when not adding any nodes", func() {
			It("contains no nodes", func() {
				Expect(g.Nodes()).To(beSetThatIsEmpty())
			})
		})

		Context("when adding two nodes", func() {
			It("the graph contains both nodes", func() {
				g = addNode(g, n1)
				g = addNode(g, n2)
				Expect(g.Nodes()).To(beSetThatConsistsOf(n1, n2))
			})
		})

		Context("when adding a new node", func() {
			It("contains just the node", func() {
				if !gIsMutable {
					Skip("Graph is not mutable.")
				}

				Expect(gAsMutable.AddNode(n1)).To(BeTrue())
				Expect(g.Nodes()).To(beSetThatConsistsOf(n1))
			})
		})

		Context("when adding an existing node", func() {
			It("contains no additional nodes", func() {
				if !gIsMutable {
					Skip("Graph is not mutable.")
				}

				g = addNode(g, n1)
				Expect(gAsMutable.AddNode(n1)).To(BeFalse())
				Expect(g.Nodes()).To(beSetThatConsistsOf(n1))
			})
		})

		Context("when checking the mutability of the nodes set", func() {
			It("is not mutable", func() {
				Expect(g.Nodes()).To(beSetThatIsNotMutable())
			})
		})
	})
}

func beSetThatConsistsOf(first int, others ...int) types.GomegaMatcher {
	all := []int{first}
	all = append(all, others...)

	return WithTransform(toSlice, ConsistOf(all))
}

func beSetThatIsEmpty() types.GomegaMatcher {
	return WithTransform(toSlice, BeEmpty())
}

func beSetThatIsNotMutable() types.GomegaMatcher {
	return WithTransform(
		func(s set.Set[int]) bool {
			_, mutable := s.(set.MutableSet[int])
			return mutable
		},
		BeFalse())
}

func toSlice(s set.Set[int]) []int {
	var result []int
	s.ForEach(func(elem int) {
		result = append(result, elem)
	})
	return result
}
