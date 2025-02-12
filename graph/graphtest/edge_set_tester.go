package graphtest

import (
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/jbduncan/go-containers/graph"
	"github.com/jbduncan/go-containers/internal/orderagnostic"
	"github.com/jbduncan/go-containers/set"
)

type edgeSetStringTester struct {
	tt              *testing.T
	setName         string
	graphIsDirected bool
	edges           string
	expectedEdges   []graph.EndpointPair[int]
}

func newEdgeSetStringTester(
	tt *testing.T,
	setName string,
	graphIsDirected bool,
	edges set.Set[graph.EndpointPair[int]],
	expectedEdges []graph.EndpointPair[int],
) *edgeSetStringTester {
	return &edgeSetStringTester{
		tt:              tt,
		setName:         setName,
		graphIsDirected: graphIsDirected,
		edges:           edges.String(),
		expectedEdges:   expectedEdges,
	}
}

func (t *edgeSetStringTester) Test() {
	t.tt.Helper()

	t.tt.Run("Set.String", func(tt *testing.T) {
		trimmed, prefixFound := strings.CutPrefix(t.edges, "[")
		if !prefixFound {
			tt.Fatalf(
				`%s: got Set.String of %q, want to have prefix "["`,
				t.setName,
				t.edges,
			)
		}
		trimmed, suffixFound := strings.CutSuffix(trimmed, "]")
		if !suffixFound {
			tt.Fatalf(
				`%s: got Set.String of %q, want to have suffix "]"`,
				t.setName,
				t.edges,
			)
		}

		elems := splitByComma(trimmed)
		want := make([]graph.EndpointPair[int], 0, len(elems))
		for _, elemStr := range elems {
			want = append(want, t.toEndpointPair(tt, elemStr))
		}

		if t.graphIsDirected {
			if diff := orderagnostic.Diff(t.expectedEdges, want); diff != "" {
				t.report(tt)
			}
		} else {
			if diff := undirectedEndpointPairsDiff(
				t.expectedEdges,
				want,
			); diff != "" {
				t.report(tt)
			}
		}
	})
}

func (t *edgeSetStringTester) toEndpointPair(
	tt *testing.T,
	s string,
) graph.EndpointPair[int] {
	tt.Helper()

	// TODO: Extract into global variable
	endpointPairStringRegex := regexp.MustCompile(`<(\d+) -> (\d)+>`)
	matches := endpointPairStringRegex.FindStringSubmatch(s)
	if len(matches) != 3 {
		t.report(tt)
	}
	source, err := strconv.Atoi(matches[1])
	if err != nil {
		t.report(tt)
	}
	target, err := strconv.Atoi(matches[2])
	if err != nil {
		t.report(tt)
	}
	return graph.EndpointPairOf(source, target)
}

func (t *edgeSetStringTester) report(tt *testing.T) {
	tt.Helper()

	var msg strings.Builder
	if len(t.expectedEdges) == 0 {
		msg.WriteString(`%s: got Set.String of %q, want "[]"`)
	} else if t.graphIsDirected {
		msg.WriteString(
			"%s: got Set.String of %q, want to contain substrings:\n")
		for _, edge := range t.expectedEdges {
			msg.WriteString("    ")
			msg.WriteString(edge.String())
		}
	} else {
		msg.WriteString(
			"%s: got Set.String of %q, want to contain substrings:\n")
		for _, edge := range t.expectedEdges {
			msg.WriteString("    ")
			msg.WriteString(edge.String())
			msg.WriteString(" or ")
			msg.WriteString(reverseOf(edge).String())
		}
	}
	tt.Fatalf(msg.String(), t.setName, t.edges)
}
