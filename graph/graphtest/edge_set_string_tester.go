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
	tt            *testing.T
	setName       string
	directionMode DirectionMode
	edges         string
	expectedEdges []graph.EndpointPair[int]
}

func newEdgeSetStringTester(
	tt *testing.T,
	setName string,
	directionMode DirectionMode,
	edges set.Set[graph.EndpointPair[int]],
	expectedEdges []graph.EndpointPair[int],
) *edgeSetStringTester {
	return &edgeSetStringTester{
		tt:            tt,
		setName:       setName,
		directionMode: directionMode,
		edges:         edges.String(),
		expectedEdges: expectedEdges,
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

		var diff string
		if t.directionMode == Directed {
			diff = orderagnostic.Diff(t.expectedEdges, want)
		} else {
			diff = undirectedEndpointPairsDiff(t.expectedEdges, want)
		}
		if diff != "" {
			t.report(ttt)
		}
	})
}

var endpointPairStringRegex = regexp.MustCompile(`<(\d+) -> (\d)+>`)

func (t *edgeSetStringTester) toEndpointPair(
	tt *testing.T,
	s string,
) graph.EndpointPair[int] {
	tt.Helper()

	matches := endpointPairStringRegex.FindStringSubmatch(s)
	if len(matches) != 3 {
		t.report(tt)
	}

	// The regex guarantees that the matches are integers.
	source, _ := strconv.Atoi(matches[1])
	target, _ := strconv.Atoi(matches[2])
	return graph.EndpointPairOf(source, target)
}

func (t *edgeSetStringTester) report(tt *testing.T) {
	tt.Helper()

	if len(t.expectedEdges) == 0 {
		tt.Fatalf(`%s: got Set.String of %q, want "[]"`, t.setName, t.edges)
	}

	var msg strings.Builder
	msg.WriteString("%s: got Set.String of %q, want to contain substrings:\n")
	for _, edge := range t.expectedEdges {
		msg.WriteString("    ")
		msg.WriteString(edge.String())
		if t.directionMode == Undirected {
			msg.WriteString(" or ")
			msg.WriteString(reverseOf(edge).String())
		}
		msg.WriteString("\n")
	}
	tt.Fatalf(msg.String(), t.setName, t.edges)
}
