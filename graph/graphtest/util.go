package graphtest

import (
	"github.com/google/go-cmp/cmp"
	"github.com/jbduncan/go-containers/graph"
	"github.com/jbduncan/go-containers/internal/orderagnostic"
)

func reverseOf(endpointPair graph.EndpointPair[int]) graph.EndpointPair[int] {
	return graph.EndpointPairOf(endpointPair.Target(), endpointPair.Source())
}

func reversesOf(edges []graph.EndpointPair[int]) []graph.EndpointPair[int] {
	result := make([]graph.EndpointPair[int], 0, len(edges))
	for _, edge := range edges {
		result = append(result, reverseOf(edge))
	}
	return result
}

func undirectedEndpointPairsDiff(
	got []graph.EndpointPair[int],
	want []graph.EndpointPair[int],
) string {
	return cmp.Diff(
		want,
		got,
		cmp.Comparer(
			func(
				a []graph.EndpointPair[int],
				b []graph.EndpointPair[int],
			) bool {
				aCopy, bCopy := deepCopyAndNormalise(a), deepCopyAndNormalise(b)
				return orderagnostic.SlicesEqual(aCopy, bCopy)
			},
		),
	)
}

func deepCopyAndNormalise(
	s []graph.EndpointPair[int],
) []graph.EndpointPair[int] {
	var result []graph.EndpointPair[int]
	for _, edge := range s {
		var newEdge graph.EndpointPair[int]
		if edge.Source() < edge.Target() {
			newEdge = graph.EndpointPairOf(edge.Source(), edge.Target())
		} else {
			newEdge = graph.EndpointPairOf(edge.Target(), edge.Source())
		}
		result = append(result, newEdge)
	}
	return result
}
