package graph_test

import (
	"testing"

	"github.com/jbduncan/go-containers/graph"
	. "github.com/jbduncan/go-containers/internal/matchers"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gcustom"
	"github.com/onsi/gomega/types"
)

func TestEndpointPair(t *testing.T) {
	t.Run("EndpointPair.Source() returns the first node", func(t *testing.T) {
		g := NewWithT(t)
		g.Expect(endpointPair()).To(haveSource("link"))
	})

	t.Run("EndpointPair.Target() returns the second node", func(t *testing.T) {
		g := NewWithT(t)
		g.Expect(endpointPair()).To(haveTarget("zelda"))
	})

	t.Run(
		"EndpointPair.AdjacentNode(EndpointPair.Source()) returns the target",
		func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(endpointPair().AdjacentNode(endpointPair().Source())).
				To(Equal(endpointPair().Target()))
		},
	)

	t.Run(
		"EndpointPair.AdjacentNode(EndpointPair.Target()) returns the source",
		func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(endpointPair().AdjacentNode(endpointPair().Target())).
				To(Equal(endpointPair().Source()))
		},
	)

	t.Run(
		"EndpointPair.AdjacentNode(nonAdjacentNode) panics",
		func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(func() { endpointPair().AdjacentNode("ganondorf") }).
				To(PanicWith("EndpointPair <link -> zelda> does not contain node ganondorf"))
		},
	)

	t.Run(
		"EndpointPair.String() returns a string representation",
		func(t *testing.T) {
			g := NewWithT(t)
			g.Expect(endpointPair()).To(HaveStringRepr("<link -> zelda>"))
		},
	)
}

func endpointPair() graph.EndpointPair[string] {
	return graph.EndpointPairOf("link", "zelda")
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
