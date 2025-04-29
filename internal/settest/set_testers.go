package settest

import (
	"fmt"
	"iter"
	"slices"
	"strings"
	"testing"

	"github.com/jbduncan/go-containers/internal/orderagnostic"
	"github.com/jbduncan/go-containers/internal/stringsx"
)

type set[T comparable] interface {
	Contains(element T) bool
	Len() int
	All() iter.Seq[T]
	String() string
}

type mutableSet[T comparable] interface {
	set[T]
	Add(element T, others ...T) bool
	Remove(element T, others ...T) bool
}

func Len[T comparable](
	t *testing.T,
	setName string,
	s set[T],
	expectedLen int,
) {
	t.Helper()

	var prefix string
	if len(setName) > 0 {
		prefix = setName + ": "
	}

	if got, want := s.Len(), expectedLen; got != want {
		t.Errorf(
			"%sgot Set.Len of %d, want %d",
			prefix,
			got,
			want,
		)
	}
}

func All[T comparable](
	t *testing.T,
	setName string,
	s set[T],
	expectedElements []T,
) {
	t.Helper()

	var prefix string
	if len(setName) > 0 {
		prefix = setName + ": "
	}

	got, want := slices.Collect(s.All()), expectedElements
	if diff := orderagnostic.Diff(got, want); diff != "" {
		t.Errorf("%sSet.All mismatch (-want +got):\n%s", prefix, diff)
	}
}

func Contains[T comparable](
	t *testing.T,
	setName string,
	s set[T],
	contains []T,
) {
	t.Helper()

	var prefix string
	if len(setName) > 0 {
		prefix = setName + ": "
	}

	for _, element := range contains {
		if !s.Contains(element) {
			t.Errorf(
				"%sgot Set.Contains(%v) == false, want true",
				prefix,
				element,
			)
		}
	}
}

func DoesNotContain[T comparable](
	t *testing.T,
	setName string,
	s set[T],
	doesNotContain []T,
) {
	t.Helper()

	var prefix string
	if len(setName) > 0 {
		prefix = setName + ": "
	}

	for _, element := range doesNotContain {
		if s.Contains(element) {
			t.Errorf(
				"%sgot Set.Contains(%v) == true, want false",
				prefix,
				element,
			)
		}
	}
}

func String[T comparable](
	t *testing.T,
	setName string,
	s set[T],
	expectedElements []T,
) {
	t.Helper()

	var prefix string
	if len(setName) > 0 {
		prefix = setName + ": "
	}

	str := s.String()
	trimmed, prefixFound := strings.CutPrefix(str, "[")
	if !prefixFound {
		t.Errorf(
			`%sgot Set.String of %q, want to have prefix "["`,
			prefix,
			str,
		)
		return
	}
	trimmed, suffixFound := strings.CutSuffix(trimmed, "]")
	if !suffixFound {
		t.Errorf(
			`%sgot Set.String of %q, want to have suffix "]"`,
			prefix,
			str,
		)
		return
	}

	want := make([]string, 0, len(expectedElements))
	for _, v := range expectedElements {
		want = append(want, fmt.Sprintf("%v", v))
	}
	got := stringsx.SplitByComma(trimmed)

	if diff := orderagnostic.Diff(got, want); diff != "" {
		t.Errorf(
			"%sSet.String of %q: elements mismatch: (-want +got):\n%s",
			prefix,
			str,
			diff,
		)
	}
}

func IsMutable[T comparable](t *testing.T, setName string, s set[T]) {
	t.Helper()

	var prefix string
	if len(setName) > 0 {
		prefix = setName + ": "
	}

	if _, mutable := s.(mutableSet[T]); mutable {
		t.Errorf(
			"%sgot a mutable set: %v, want just a set",
			prefix,
			s,
		)
	}
}
