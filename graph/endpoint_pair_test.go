package graph_test

import (
	"testing"

	"github.com/jbduncan/go-containers/graph"
)

func TestEndpointPair(t *testing.T) {
	t.Parallel()

	t.Run("EndpointPair.Source() returns the first node", func(t *testing.T) {
		t.Parallel()
		if got, want := endpointPair().Source(), "link"; got != want {
			t.Errorf("EndpointPair.Source: got %q, want %q", got, want)
		}
	})

	t.Run("EndpointPair.Target() returns the second node", func(t *testing.T) {
		t.Parallel()
		if got, want := endpointPair().Target(), "zelda"; got != want {
			t.Errorf("EndpointPair.Source: got %q, want %q", got, want)
		}
	})

	t.Run(
		"EndpointPair.AdjacentNode(EndpointPair.Source()) returns the target",
		func(t *testing.T) {
			t.Parallel()
			got := endpointPair().AdjacentNode(endpointPair().Source())
			want := endpointPair().Target()
			if got != want {
				t.Errorf(
					"EndpointPair.AdjacentNode(%q): got %q, want %q",
					endpointPair().Source(),
					got,
					want,
				)
			}
		},
	)

	t.Run(
		"EndpointPair.AdjacentNode(EndpointPair.Target()) returns the source",
		func(t *testing.T) {
			t.Parallel()
			got := endpointPair().AdjacentNode(endpointPair().Target())
			want := endpointPair().Source()
			if got != want {
				t.Errorf(
					"EndpointPair.AdjacentNode(%q): got %q, want %q",
					endpointPair().Target(),
					got,
					want,
				)
			}
		},
	)

	t.Run(
		"EndpointPair.AdjacentNode(nonAdjacentNode) panics",
		func(t *testing.T) {
			t.Parallel()
			defer func() { _ = recover() }()
			endpointPair().AdjacentNode("ganondorf")
			t.Errorf(`EndpointPair.AdjacentNode("ganondorf"): should have panicked`)
		},
	)

	t.Run(
		"EndpointPair.String() returns a string representation",
		func(t *testing.T) {
			t.Parallel()
			if got, want := endpointPair().String(), "<link -> zelda>"; got != want {
				t.Errorf("EndpointPair.String: got %s, want %s", got, want)
			}
		},
	)
}

func endpointPair() graph.EndpointPair[string] {
	return graph.EndpointPairOf("link", "zelda")
}
