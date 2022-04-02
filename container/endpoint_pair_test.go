package container_test

import (
	"github.com/onsi/gomega/types"
	"go-containers/container"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "go-containers/internal/matchers"
)

var _ = Describe("EndpointPair", func() {
	Describe("given a new unordered endpoint pair", func() {
		var endpointPair container.EndpointPair[string]

		BeforeEach(func() {
			endpointPair = container.NewUnorderedEndpointPair("link", "zelda")
		})

		Context("when calling .IsOrdered()", func() {
			It("returns false", func() {
				Expect(endpointPair).ToNot(BeOrdered())
			})
		})

		Context("when calling .Source()", func() {
			It("panics", func() {
				Expect(endpointPair).To(HaveUnavailableSource())
			})
		})

		Context("when calling .Target()", func() {
			It("panics", func() {
				Expect(endpointPair).To(HaveUnavailableTarget())
			})
		})

		Context("when calling .NodeU()", func() {
			It("returns the first node", func() {
				Expect(endpointPair).To(HaveNodeU("link"))
			})
		})

		Context("when calling .NodeV()", func() {
			It("returns the second node", func() {
				Expect(endpointPair).To(HaveNodeV("zelda"))
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
		var endpointPair container.EndpointPair[string]

		BeforeEach(func() {
			endpointPair = container.NewOrderedEndpointPair("link", "zelda")
		})

		Context("when calling .IsOrdered()", func() {
			It("returns true", func() {
				Expect(endpointPair).To(BeOrdered())
			})
		})

		Context("when calling .Source()", func() {
			It("returns the first node", func() {
				Expect(endpointPair).To(HaveSource("link"))
			})
		})

		Context("when calling .Target()", func() {
			It("panics", func() {
				Expect(endpointPair).To(HaveTarget("zelda"))
			})
		})

		Context("when calling .NodeU()", func() {
			It("returns the first node", func() {
				Expect(endpointPair).To(HaveNodeU("link"))
			})
		})

		Context("when calling .NodeV()", func() {
			It("returns the second node", func() {
				Expect(endpointPair).To(HaveNodeV("zelda"))
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

func BeOrdered() types.GomegaMatcher {
	return WithTransform(
		func(endpointPair container.EndpointPair[string]) bool {
			return endpointPair.IsOrdered()
		},
		BeTrue())
}

func HaveSource(source string) types.GomegaMatcher {
	return WithTransform(
		func(endpointPair container.EndpointPair[string]) string {
			return endpointPair.Source()
		},
		Equal(source))
}

func HaveTarget(target string) types.GomegaMatcher {
	return WithTransform(
		func(endpointPair container.EndpointPair[string]) string {
			return endpointPair.Target()
		},
		Equal(target))
}

func HaveUnavailableSource() types.GomegaMatcher {
	return WithTransform(
		func(endpointPair container.EndpointPair[string]) func() {
			return func() { endpointPair.Source() }
		},
		PanicWith(notAvailableOnUndirected))
}

func HaveUnavailableTarget() types.GomegaMatcher {
	return WithTransform(
		func(endpointPair container.EndpointPair[string]) func() {
			return func() { endpointPair.Target() }
		},
		PanicWith(notAvailableOnUndirected))
}

const notAvailableOnUndirected = "cannot call Source()/Target() on an EndpointPair from an " +
	"undirected graph; consider calling AdjacentNode(node) if you already have a node, or " +
	"NodeU()/NodeV() if you don't"

func HaveNodeU(node string) types.GomegaMatcher {
	return WithTransform(
		func(endpointPair container.EndpointPair[string]) string {
			return endpointPair.NodeU()
		},
		Equal(node))
}

func HaveNodeV(node string) types.GomegaMatcher {
	return WithTransform(
		func(endpointPair container.EndpointPair[string]) string {
			return endpointPair.NodeV()
		},
		Equal(node))
}
