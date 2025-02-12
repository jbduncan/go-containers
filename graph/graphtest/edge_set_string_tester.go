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
	tt                   *testing.T
	setName              string
	directedOrUndirected DirectionMode
	edges                string
	expectedEdges        []graph.EndpointPair[int]
}

func newEdgeSetStringTester(
	tt *testing.T,
	setName string,
	directedOrUndirected DirectionMode,
	edges set.Set[graph.EndpointPair[int]],
	expectedEdges []graph.EndpointPair[int],
) *edgeSetStringTester {
	return &edgeSetStringTester{
		tt:                   tt,
		setName:              setName,
		directedOrUndirected: directedOrUndirected,
		edges:                edges.String(),
		expectedEdges:        expectedEdges,
	}
}

func (t *edgeSetStringTester) Test() {
	t.tt.Helper()

	t.tt.Run("Set.String", func(ttt *testing.T) {
		trimmed, prefixFound := strings.CutPrefix(t.edges, "[")
		if !prefixFound {
			ttt.Fatalf(
				`%s: got Set.String of %q, want to have prefix "["`,
				t.setName,
				t.edges,
			)
		}
		trimmed, suffixFound := strings.CutSuffix(trimmed, "]")
		if !suffixFound {
			ttt.Fatalf(
				`%s: got Set.String of %q, want to have suffix "]"`,
				t.setName,
				t.edges,
			)
		}

		elems := splitByComma(trimmed)
		want := make([]graph.EndpointPair[int], 0, len(elems))
		for _, elemStr := range elems {
			want = append(want, t.toEndpointPair(ttt, elemStr))
		}

		switch t.directedOrUndirected {
		case Directed:
			if diff := orderagnostic.Diff(t.expectedEdges, want); diff != "" {
				t.report(ttt)
			}
		case Undirected:
			if diff := undirectedEndpointPairsDiff(
				t.expectedEdges,
				want,
			); diff != "" {
				t.report(ttt)
			}
		default:
			panic("unreachable")
		}
	})
}

var endpointPairStringRegex = regexp.MustCompile(`<(\d+) -> (\d)+>`)

func (t *edgeSetStringTester) toEndpointPair(
	ttt *testing.T,
	s string,
) graph.EndpointPair[int] {
	ttt.Helper()

	matches := endpointPairStringRegex.FindStringSubmatch(s)
	if len(matches) != 3 {
		t.report(ttt)
	}
	source, err := strconv.Atoi(matches[1])
	if err != nil {
		t.report(ttt)
	}
	target, err := strconv.Atoi(matches[2])
	if err != nil {
		t.report(ttt)
	}
	return graph.EndpointPairOf(source, target)
}

func (t *edgeSetStringTester) report(ttt *testing.T) {
	ttt.Helper()

	var msg strings.Builder
	switch {
	case len(t.expectedEdges) == 0:
		msg.WriteString(`%s: got Set.String of %q, want "[]"`)
	case t.directedOrUndirected == Directed:
		msg.WriteString(
			"%s: got Set.String of %q, want to contain substrings:\n")
		for _, edge := range t.expectedEdges {
			msg.WriteString("    ")
			msg.WriteString(edge.String())
		}
	case t.directedOrUndirected == Undirected:
		msg.WriteString(
			"%s: got Set.String of %q, want to contain substrings:\n")
		for _, edge := range t.expectedEdges {
			msg.WriteString("    ")
			msg.WriteString(edge.String())
			msg.WriteString(" or ")
			msg.WriteString(reverseOf(edge).String())
		}
	default:
		panic("unreachable")
	}
	ttt.Fatalf(msg.String(), t.setName, t.edges)
}
