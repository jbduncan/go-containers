package graphtest

import (
	"strings"

	"github.com/jbduncan/go-containers/graph"
)

func reverseOf(endpointPair graph.EndpointPair[int]) graph.EndpointPair[int] {
	return graph.EndpointPairOf(endpointPair.Target(), endpointPair.Source())
}

func splitByComma(s string) []string {
	if len(s) == 0 {
		return make([]string, 0)
	}
	return strings.Split(s, ", ")
}
