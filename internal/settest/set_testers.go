package settest

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/jbduncan/go-containers/internal/orderagnostic"
	"github.com/jbduncan/go-containers/internal/stringsx"
	"github.com/jbduncan/go-containers/set"
)

func SetLen[T comparable](
	t *testing.T,
	setName string,
	s set.Set[T],
	expectedLen int,
) {
	t.Helper()

	if got, want := s.Len(), expectedLen; got != want {
		t.Errorf(
			"%s: got Set.Len of %d, want %d",
			setName,
			got,
			want,
		)
	}
}

func SetAll[T comparable](
	t *testing.T,
	setName string,
	s set.Set[T],
	expectedValues []T,
) {
	t.Helper()

	got, want := slices.Collect(s.All()), expectedValues
	if diff := orderagnostic.Diff(got, want); diff != "" {
		t.Errorf("%s: Set.All mismatch (-want +got):\n%s", setName, diff)
	}
}

func SetContains[T comparable](
	t *testing.T,
	setName string,
	s set.Set[T],
	contains []T,
	doesNotContain []T,
) {
	t.Helper()

	for _, value := range contains {
		if !s.Contains(value) {
			t.Errorf(
				"%s: got Set.Contains(%v) == false, want true",
				setName,
				value,
			)
		}
	}
	for _, value := range doesNotContain {
		if s.Contains(value) {
			t.Errorf(
				"%s: got Set.Contains(%v) == true, want false",
				setName,
				value,
			)
		}
	}
}

func SetString[T comparable](
	t *testing.T,
	setName string,
	s set.Set[T],
	expectedValues []T,
) {
	t.Helper()

	str := s.String()
	trimmed, prefixFound := strings.CutPrefix(str, "[")
	if !prefixFound {
		t.Errorf(
			`%s: got Set.String of %q, want to have prefix "["`,
			setName,
			str,
		)
		return
	}
	trimmed, suffixFound := strings.CutSuffix(trimmed, "]")
	if !suffixFound {
		t.Errorf(
			`%s: got Set.String of %q, want to have suffix "]"`,
			setName,
			str,
		)
		return
	}

	want := make([]string, 0, len(expectedValues))
	for _, v := range expectedValues {
		want = append(want, fmt.Sprintf("%v", v))
	}
	got := stringsx.SplitByComma(trimmed)

	if diff := orderagnostic.Diff(got, want); diff != "" {
		t.Errorf(
			"%s: Set.String of %q: elements mismatch: (-want +got):\n%s",
			setName,
			str,
			diff,
		)
	}
}
