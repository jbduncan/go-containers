package graphtest

import (
	"slices"
	"testing"

	"github.com/jbduncan/go-containers/graph"
	"github.com/jbduncan/go-containers/internal/orderagnostic"
	"github.com/jbduncan/go-containers/internal/slicesx"
	"github.com/jbduncan/go-containers/set"
)

type edgeSetTester struct {
	t             *testing.T
	setName       string
	edges         set.Set[graph.EndpointPair[int]]
	directionMode DirectionMode
	expectedEdges []graph.EndpointPair[int]
}

func (tt edgeSetTester) test() {
	tt.t.Helper()

	var contains []graph.EndpointPair[int]
	var doesNotContain []graph.EndpointPair[int]
	if tt.directionMode == Directed {
		contains = tt.expectedEdges
		doesNotContain = slicesx.AllOf(
			graph.EndpointPairOf(nodeNotInGraph, nodeNotInGraph),
			reversesOf(tt.expectedEdges),
		)
	} else {
		contains = slices.Concat(
			tt.expectedEdges,
			reversesOf(tt.expectedEdges),
		)
		doesNotContain = []graph.EndpointPair[int]{
			graph.EndpointPairOf(nodeNotInGraph, nodeNotInGraph),
		}
	}

	testSetLen(tt.t, tt.setName, tt.edges, len(tt.expectedEdges))
	tt.testEdgeSetAll(tt.t, tt.setName, tt.edges, tt.expectedEdges)
	testSetContains(tt.t, tt.setName, tt.edges, contains, doesNotContain)
	newEdgeSetStringTester(
		tt.t,
		tt.setName,
		tt.directionMode,
		tt.edges,
		tt.expectedEdges,
	).test()
}

func (tt edgeSetTester) testEdgeSetAll(
	t *testing.T,
	setName string,
	edges set.Set[graph.EndpointPair[int]],
	expectedEdges []graph.EndpointPair[int],
) {
	t.Helper()

	got, want := slices.Collect(edges.All()), expectedEdges
	if tt.directionMode == Directed {
		if diff := orderagnostic.Diff(got, want); diff != "" {
			t.Errorf(
				"%s: Set.All mismatch (-want +got):\n%s",
				setName,
				diff,
			)
		}
	} else {
		if diff := undirectedEndpointPairsDiff(got, want); diff != "" {
			t.Errorf(
				"%s: Set.All mismatch (-want +got):\n%s",
				setName,
				diff,
			)
		}
	}
}
