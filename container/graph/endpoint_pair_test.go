package graph_test

import (
	"github.com/onsi/gomega/types"
	"go-containers/container/graph"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "go-containers/internal/matchers"
)

var _ = Describe("EndpointPair", func() {
	Describe("given a new unordered endpoint pair", func() {
		var endpointPair graph.EndpointPair[string]

		BeforeEach(func() {
			endpointPair = graph.NewUnorderedEndpointPair("link", "zelda")
		})

		Context("when calling .IsOrdered()", func() {
			It("returns false", func() {
				Expect(endpointPair).ToNot(beOrdered())
			})
		})

		Context("when calling .Source()", func() {
			It("panics", func() {
				Expect(endpointPair).To(haveUnavailableSource())
			})
		})

		Context("when calling .Target()", func() {
			It("panics", func() {
				Expect(endpointPair).To(haveUnavailableTarget())
			})
		})

		Context("when calling .NodeU()", func() {
			It("returns the first node", func() {
				Expect(endpointPair).To(haveNodeU("link"))
			})
		})

		Context("when calling .NodeV()", func() {
			It("returns the second node", func() {
				Expect(endpointPair).To(haveNodeV("zelda"))
			})
		})

		Context("when calling .AdjacentNode() with NodeU", func() {
			It("returns NodeV", func() {
				Expect(endpointPair.AdjacentNode(endpointPair.NodeU())).
					To(Equal(endpointPair.NodeV()))
			})
		})

		Context("when calling .AdjacentNode() with NodeV", func() {
			It("returns NodeU", func() {
				Expect(endpointPair.AdjacentNode(endpointPair.NodeV())).
					To(Equal(endpointPair.NodeU()))
			})
		})

		Context("when calling .AdjacentNode() with a non-adjacent node", func() {
			It("panics", func() {
				Expect(func() { endpointPair.AdjacentNode("ganondorf") }).
					To(Or(
						PanicWith("EndpointPair [link, zelda] does not contain node ganondorf"),
						PanicWith("EndpointPair [zelda, link] does not contain node ganondorf")))
			})
		})

		Context("when calling .String()", func() {
			It("returns a string representation", func() {
				Expect(endpointPair).To(
					HaveStringRepr(BeElementOf("[link, zelda]", "[zelda, link]")))
			})
		})
	})

	Describe("given a new ordered endpoint pair", func() {
		var endpointPair graph.EndpointPair[string]

		BeforeEach(func() {
			endpointPair = graph.NewOrderedEndpointPair("link", "zelda")
		})

		Context("when calling .IsOrdered()", func() {
			It("returns true", func() {
				Expect(endpointPair).To(beOrdered())
			})
		})

		Context("when calling .Source()", func() {
			It("returns the first node", func() {
				Expect(endpointPair).To(haveSource("link"))
			})
		})

		Context("when calling .Target()", func() {
			It("panics", func() {
				Expect(endpointPair).To(haveTarget("zelda"))
			})
		})

		Context("when calling .NodeU()", func() {
			It("returns the first node", func() {
				Expect(endpointPair).To(haveNodeU("link"))
			})
		})

		Context("when calling .NodeV()", func() {
			It("returns the second node", func() {
				Expect(endpointPair).To(haveNodeV("zelda"))
			})
		})

		Context("when calling .AdjacentNode() with NodeU", func() {
			It("returns NodeV", func() {
				Expect(endpointPair.AdjacentNode(endpointPair.NodeU())).
					To(Equal(endpointPair.NodeV()))
			})
		})

		Context("when calling .AdjacentNode() with NodeV", func() {
			It("returns NodeU", func() {
				Expect(endpointPair.AdjacentNode(endpointPair.NodeV())).
					To(Equal(endpointPair.NodeU()))
			})
		})

		Context("when calling .AdjacentNode() with a non-adjacent node", func() {
			It("panics", func() {
				Expect(func() { endpointPair.AdjacentNode("ganondorf") }).
					To(PanicWith("EndpointPair <link -> zelda> does not contain node ganondorf"))
			})
		})

		Context("when calling .String()", func() {
			It("returns a string representation", func() {
				Expect(endpointPair).To(HaveStringRepr("<link -> zelda>"))
			})
		})
	})
})

func beOrdered() types.GomegaMatcher {
	return WithTransform(
		func(endpointPair graph.EndpointPair[string]) bool {
			return endpointPair.IsOrdered()
		},
		BeTrue())
}

func haveSource(source string) types.GomegaMatcher {
	return WithTransform(
		func(endpointPair graph.EndpointPair[string]) string {
			return endpointPair.Source()
		},
		Equal(source))
}

func haveTarget(target string) types.GomegaMatcher {
	return WithTransform(
		func(endpointPair graph.EndpointPair[string]) string {
			return endpointPair.Target()
		},
		Equal(target))
}

func haveUnavailableSource() types.GomegaMatcher {
	return WithTransform(
		func(endpointPair graph.EndpointPair[string]) func() {
			return func() { endpointPair.Source() }
		},
		PanicWith(notAvailableOnUndirected))
}

func haveUnavailableTarget() types.GomegaMatcher {
	return WithTransform(
		func(endpointPair graph.EndpointPair[string]) func() {
			return func() { endpointPair.Target() }
		},
		PanicWith(notAvailableOnUndirected))
}

const notAvailableOnUndirected = "cannot call Source()/Target() on an EndpointPair from an " +
	"undirected graph; consider calling AdjacentNode(node) if you already have a node, or " +
	"NodeU()/NodeV() if you don't"

func haveNodeU(node string) types.GomegaMatcher {
	return WithTransform(
		func(endpointPair graph.EndpointPair[string]) string {
			return endpointPair.NodeU()
		},
		Equal(node))
}

func haveNodeV(node string) types.GomegaMatcher {
	return WithTransform(
		func(endpointPair graph.EndpointPair[string]) string {
			return endpointPair.NodeV()
		},
		Equal(node))
}

// TODO: Guava's EndpointPairTest.java has some tests that use EndpointPair but check Graph.edges()
//       and Network.asGraph().edges(). Adapt these tests for our own Graph and Network types.
