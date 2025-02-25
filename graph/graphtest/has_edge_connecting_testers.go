package graphtest

import (
	"testing"

	"github.com/jbduncan/go-containers/graph"
)

func testHasEdgeConnecting(
	t *testing.T,
	g graph.Graph[int],
	source, target int,
) {
	if got := g.HasEdgeConnecting(source, target); !got {
		t.Errorf("Graph.HasEdgeConnecting: got false, want true")
	}
	if got := g.HasEdgeConnectingEndpoints(
		graph.EndpointPairOf(source, target),
	); !got {
		t.Errorf("Graph.HasEdgeConnectingEndpoints: got false, want true")
	}
}

func testHasNoEdgeConnecting(
	t *testing.T,
	g graph.Graph[int],
	source, target int,
) {
	if got := g.HasEdgeConnecting(source, target); got {
		t.Errorf("Graph.HasEdgeConnecting: got true, want false")
	}
	if got := g.HasEdgeConnectingEndpoints(
		graph.EndpointPairOf(source, target),
	); got {
		t.Errorf("Graph.HasEdgeConnectingEndpoints: got true, want false")
	}
}
