package graph_test

import (
	"github.com/jbduncan/go-containers/graph"
	. "github.com/jbduncan/go-containers/internal/matchers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gcustom"
	"github.com/onsi/gomega/types"
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
					To(PanicWith(MatchRegexp(`^EndpointPair (\[link, zelda\]|\[zelda, link\]) does not contain node ganondorf$`)))
			})
		})

		Context("when calling .Equal() with an equivalent unordered endpoint pair", func() {
			It("returns true", func() {
				other := graph.NewUnorderedEndpointPair("link", "zelda")
				Expect(endpointPair).To(BeEquivalentToUsingEqualMethod(other))
			})
		})

		Context("when calling .Equal() with an ordered endpoint pair", func() {
			It("returns false", func() {
				other := graph.NewOrderedEndpointPair("link", "zelda")
				Expect(endpointPair).ToNot(BeEquivalentToUsingEqualMethod(other))
			})
		})

		Context("when calling .Equal() with an unordered endpoint pair with a different NodeU", func() {
			It("returns false", func() {
				other := graph.NewUnorderedEndpointPair("ganon", "zelda")
				Expect(endpointPair).ToNot(BeEquivalentToUsingEqualMethod(other))
			})
		})

		Context("when calling .Equal() with an unordered endpoint pair with a different NodeV", func() {
			It("returns false", func() {
				other := graph.NewUnorderedEndpointPair("link", "ganon")
				Expect(endpointPair).ToNot(BeEquivalentToUsingEqualMethod(other))
			})
		})

		Context("when calling .Equal() with a reversed unordered endpoint pair", func() {
			It("returns true", func() {
				other := graph.NewUnorderedEndpointPair("zelda", "link")
				Expect(endpointPair).To(BeEquivalentToUsingEqualMethod(other))
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

		Context("when calling .Equal() with an equivalent ordered endpoint pair", func() {
			It("returns true", func() {
				other := graph.NewOrderedEndpointPair("link", "zelda")
				Expect(endpointPair).To(BeEquivalentToUsingEqualMethod(other))
			})
		})

		Context("when calling .Equal() with an unordered endpoint pair", func() {
			It("returns false", func() {
				other := graph.NewUnorderedEndpointPair("link", "zelda")
				Expect(endpointPair).ToNot(BeEquivalentToUsingEqualMethod(other))
			})
		})

		Context("when calling .Equal() with an ordered endpoint pair with a different source", func() {
			It("returns false", func() {
				other := graph.NewOrderedEndpointPair("ganon", "zelda")
				Expect(endpointPair).ToNot(BeEquivalentToUsingEqualMethod(other))
			})
		})

		Context("when calling .Equal() with an ordered endpoint pair with a different target", func() {
			It("returns false", func() {
				other := graph.NewOrderedEndpointPair("link", "ganon")
				Expect(endpointPair).ToNot(BeEquivalentToUsingEqualMethod(other))
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
	return gcustom.MakeMatcher(
		func(endpointPair graph.EndpointPair[string]) (bool, error) {
			return endpointPair.IsOrdered(), nil
		}).
		WithTemplate("Expected\n{{.FormattedActual}}\n{{.To}} be ordered")
}

func haveSource(source string) types.GomegaMatcher {
	return gcustom.MakeMatcher(
		func(endpointPair graph.EndpointPair[string]) (bool, error) {
			return endpointPair.Source() == source, nil
		}).
		WithTemplate("Expected\n{{.FormattedActual}}\n{{.To}} to have a source equal to\n{{format .Data 1}}").
		WithTemplateData(source)
}

func haveTarget(target string) types.GomegaMatcher {
	return gcustom.MakeMatcher(
		func(endpointPair graph.EndpointPair[string]) (bool, error) {
			return endpointPair.Target() == target, nil
		}).
		WithTemplate("Expected\n{{.FormattedActual}}\n{{.To}} to have a target equal to\n{{format .Data 1}}").
		WithTemplateData(target)
}

func haveUnavailableSource() types.GomegaMatcher {
	return gcustom.MakeMatcher(
		func(endpointPair graph.EndpointPair[string]) (bool, error) {
			sourceMethod := func() { endpointPair.Source() }
			return PanicWith(notAvailableOnUndirected).Match(sourceMethod)
		}).
		WithTemplate("Expected Source() of\n{{.FormattedActual}}\n{{.To}} to panic with\n{{format .Data 1}}").
		WithTemplateData(notAvailableOnUndirected)
}

func haveUnavailableTarget() types.GomegaMatcher {
	return gcustom.MakeMatcher(
		func(endpointPair graph.EndpointPair[string]) (bool, error) {
			sourceMethod := func() { endpointPair.Target() }
			return PanicWith(notAvailableOnUndirected).Match(sourceMethod)
		}).
		WithTemplate("Expected Target() of\n{{.FormattedActual}}\n{{.To}} to panic with\n{{format .Data 1}}").
		WithTemplateData(notAvailableOnUndirected)
}

const notAvailableOnUndirected = "cannot call Source()/Target() on an EndpointPair from an " +
	"undirected graph; consider calling AdjacentNode(node) if you already have a node, or " +
	"NodeU()/NodeV() if you don't"

func haveNodeU(node string) types.GomegaMatcher {
	return gcustom.MakeMatcher(
		func(endpointPair graph.EndpointPair[string]) (bool, error) {
			return endpointPair.NodeU() == node, nil
		}).
		WithTemplate("Expected\n{{.FormattedActual}}\n{{.To}} to have a nodeU equal to\n{{format .Data 1}}").
		WithTemplateData(node)
}

func haveNodeV(node string) types.GomegaMatcher {
	return gcustom.MakeMatcher(
		func(endpointPair graph.EndpointPair[string]) (bool, error) {
			return endpointPair.NodeV() == node, nil
		}).
		WithTemplate("Expected\n{{.FormattedActual}}\n{{.To}} to have a nodeV equal to\n{{format .Data 1}}").
		WithTemplateData(node)
}

// TODO: Guava's EndpointPairTest.java has some tests that use EndpointPair but check Graph.edges()
//       and Network.asGraph().edges() for both directed and undirected graphs. Adapt these tests
//       for our own graph types.
