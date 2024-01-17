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
	Describe("given a new endpoint pair", func() {
		var endpointPair graph.EndpointPair[string]

		BeforeEach(func() {
			endpointPair = graph.EndpointPairOf("link", "zelda")
		})

		Context("when calling .Source()", func() {
			It("returns the first node", func() {
				Expect(endpointPair).To(haveSource("link"))
			})
		})

		Context("when calling .Target()", func() {
			It("returns the second node", func() {
				Expect(endpointPair).To(haveTarget("zelda"))
			})
		})

		Context("when calling .AdjacentNode() with Source", func() {
			It("returns Target", func() {
				Expect(endpointPair.AdjacentNode(endpointPair.Source())).
					To(Equal(endpointPair.Target()))
			})
		})

		Context("when calling .AdjacentNode() with Target", func() {
			It("returns Source", func() {
				Expect(endpointPair.AdjacentNode(endpointPair.Target())).
					To(Equal(endpointPair.Source()))
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
