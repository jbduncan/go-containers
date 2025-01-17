package settest

import (
	"fmt"
	"slices"
	"testing"

	"github.com/jbduncan/go-containers/internal/orderagnostic"
	"github.com/jbduncan/go-containers/set"
)

// TestingT is an interface for the parts of *testing.T that settest.Set needs
// to run. Whenever you see an argument of this type, pass in an instance of
// *testing.T or your unit testing framework's equivalent.
type TestingT interface {
	Helper()
	Run(name string, f func(t *testing.T)) bool
}

func Set(t TestingT, setBuilder func(elems []string) set.Set[string]) {
	tt := newTester(t, setBuilder)

	tt.emptySetHasLengthOfZero()

	tt.emptySetContainsNothing()

	tt.emptySetIterationDoesNothing()

	tt.emptySetHasEmptyStringRepr()

	tt.emptySetRemoveDoesNothing()

	tt.oneElementSetHasLengthOfOne()
	tt.emptySetPlusOneHasLengthOfOne()

	tt.oneElementSetContainsPresentElement()
	tt.emptySetPlusOneContainsPresentElement()

	tt.oneElementSetDoesNotContainAbsentElement()
	tt.emptySetPlusOneDoesNotContainAbsentElement()

	tt.emptySetPlusOneMinusOneDoesNotContainAnything()

	tt.oneElementSetReturnsElementOnIteration()
	tt.emptySetPlusOneReturnsElementOnIteration()

	tt.oneElementSetHasOneElementStringRepr()
	tt.emptySetPlusOneHasOneElementStringRepr()

	tt.twoElementSetHasLengthOfTwo()
	tt.emptySetPlusTwoHasLengthOfTwo()

	tt.twoElementSetContainsBothElements()
	tt.emptySetPlusTwoContainsBothElements()

	tt.twoElementSetReturnsBothElementsOnIteration()
	tt.emptySetPlusTwoReturnsBothElementsOnIteration()
	tt.emptySetPlusVarargsReturnsBothElementsOnIteration()

	tt.twoElementSetHasTwoElementStringRepr()
	tt.emptySetPlusTwoReturnsTwoElementStringRepr()

	tt.emptySetPlusTwoMinusOneHasLengthOfOne()

	tt.emptySetPlusTwoMinusVarargsHasLengthOfZero()

	tt.threeElementSetContainsAllThreeElements()
	tt.emptySetPlusThreeContainsAllThreeElements()

	tt.threeElementSetHasThreeElementStringRepr()
	tt.emptySetPlusThreeHasThreeElementStringRepr()

	tt.setInitializedFromTwoOfSameElementHasLengthOfOne()
	tt.emptySetPlusSameElementTwiceHasLengthOfOne()

	tt.setInitializedFromTwoOfSameElementReturnsOneElementOnIteration()
	tt.emptySetPlusSameElementTwiceReturnsOneElementOnIteration()

	tt.emptySetPlusOneReturnsTrue()

	tt.emptySetPlusSameElementTwiceReturnsFalse()

	tt.emptySetPlusSameElementTwiceThenDifferentOnceReturnsTrue()

	tt.emptySetPlusOnePlusVarargsReturnsTrue()

	tt.emptySetMinusOneReturnsFalse()

	tt.emptySetPlusOneMinusSameElementReturnsTrue()

	tt.emptySetPlusOneMinusSameElementTwiceReturnsFalse()

	tt.emptySetPlusOneMinusVarargsReturnsTrue()
}

type tester struct {
	t          TestingT
	setBuilder func(elems []string) set.Set[string]
	mutable    bool
}

func newTester(
	t TestingT,
	setBuilder func(elems []string) set.Set[string],
) *tester {
	_, mutable := setBuilder(empty()).(set.MutableSet[string])
	return &tester{
		t:          t,
		setBuilder: setBuilder,
		mutable:    mutable,
	}
}

func (tt tester) runEmptyIfMutable(
	name string,
	f func(t *testing.T, s set.MutableSet[string]),
) {
	tt.t.Helper()

	if s, mutable := tt.setBuilder(empty()).(set.MutableSet[string]); mutable {
		tt.t.Run(name, func(t *testing.T) {
			f(t, s)
		})
	}
}

func (tt tester) emptySetHasLengthOfZero() {
	tt.t.Helper()

	tt.t.Run(
		"empty set: has length of 0",
		func(t *testing.T) {
			t.Helper()
			s := tt.setBuilder(empty())

			if got, want := s.Len(), 0; got != want {
				t.Fatalf("got Set.Len of %d, want %d", got, want)
			}
		})
}

func (tt tester) emptySetContainsNothing() {
	tt.t.Helper()

	tt.t.Run(
		"empty set: contains nothing",
		func(t *testing.T) {
			s := tt.setBuilder(empty())

			if s.Contains(a) {
				t.Fatalf("got Set.Contains(%q) == true, want false", a)
			}
		})
}

func (tt tester) emptySetIterationDoesNothing() {
	tt.t.Helper()

	tt.t.Run(
		"empty set: iteration does nothing",
		func(t *testing.T) {
			s := tt.setBuilder(empty())

			if got, want := slices.Collect(s.All()), empty();
				!slices.Equal(got, want) {
				t.Fatalf("got Set.All of %q, want empty", got)
			}
		})
}

func (tt tester) emptySetHasEmptyStringRepr() {
	tt.t.Helper()

	tt.t.Run(
		"empty set: has empty string representation",
		func(t *testing.T) {
			s := tt.setBuilder(empty())

			if got, want := s.String(), "[]"; got != want {
				t.Fatalf("got Set.String of %q, want %q", got, want)
			}
		})
}

func (tt tester) emptySetRemoveDoesNothing() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: remove does nothing",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()

			s.Remove(a)

			if got, want := s.Len(), 0; got != want {
				t.Fatalf("got Set.Len of %v, want %v", got, want)
			}
		})
}

func (tt tester) oneElementSetHasLengthOfOne() {
	tt.t.Helper()

	tt.t.Run(
		"one element set: has length of 1",
		func(t *testing.T) {
			s := tt.setBuilder(oneElement())

			if got, want := s.Len(), 1; got != want {
				t.Fatalf("got Set.Len of %v, want %v", got, want)
			}
		})
}

func (tt tester) emptySetPlusOneHasLengthOfOne() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add: has length of 1",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()

			s.Add(a)

			if got, want := s.Len(), 1; got != want {
				t.Fatalf("got Set.Len of %v, want %v", got, want)
			}
		})
}

func (tt tester) oneElementSetContainsPresentElement() {
	tt.t.Helper()

	tt.t.Run(
		"one element set: contains present element",
		func(t *testing.T) {
			s := tt.setBuilder(oneElement())

			if !s.Contains(a) {
				t.Fatalf("got Set.Contains(%q) == false, want true", a)
			}
		})
}

func (tt tester) emptySetPlusOneContainsPresentElement() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add: contains present element",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()

			s.Add(a)

			if !s.Contains(a) {
				t.Fatalf("got Set.Contains(%q) == false, want true", a)
			}
		})
}

func (tt tester) oneElementSetDoesNotContainAbsentElement() {
	tt.t.Helper()

	tt.t.Run(
		"one element set: does not contain absent element",
		func(t *testing.T) {
			s := tt.setBuilder(oneElement())

			if s.Contains(b) {
				t.Fatalf("got Set.Contains(%q) == true, want false", b)
			}
		})
}

func (tt tester) emptySetPlusOneDoesNotContainAbsentElement() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add: does not contain absent element",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()

			s.Add(a)

			if s.Contains(b) {
				t.Fatalf("got Set.Contains(%q) == true, want false", b)
			}
		})
}

func (tt tester) oneElementSetReturnsElementOnIteration() {
	tt.t.Helper()

	tt.t.Run(
		"one element set: returns element on iteration",
		func(t *testing.T) {
			s := tt.setBuilder(oneElement())

			if got, want := slices.Collect(s.All()), oneElement();
				!slices.Equal(got, want) {
				t.Fatalf("got Set.All of %q, want %q", got, want)
			}
		})
}

func (tt tester) emptySetPlusOneReturnsElementOnIteration() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add: returns element on iteration",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()

			s.Add(a)

			if got, want := slices.Collect(s.All()), oneElement();
				!slices.Equal(got, want) {
				t.Fatalf("got Set.All of %q, want %q", got, want)
			}
		})
}

func (tt tester) oneElementSetHasOneElementStringRepr() {
	tt.t.Helper()

	tt.t.Run(
		"one element set: has one-element string representation",
		func(t *testing.T) {
			s := tt.setBuilder(oneElement())

			if got, want := s.String(), aString(); got != want {
				t.Fatalf("got Set.String of %q, want %q", got, want)
			}
		})
}

func (tt tester) emptySetPlusOneHasOneElementStringRepr() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add: has one-element string representation",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()

			s.Add(a)

			if got, want := s.String(), aString(); got != want {
				t.Fatalf("got Set.String of %q, want %q", got, want)
			}
		})
}

func (tt tester) emptySetPlusOneMinusOneDoesNotContainAnything() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add: remove: does not contain anything",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()

			s.Add(a)
			s.Remove(a)

			if s.Contains(a) {
				t.Fatalf("got Set.Contains(%q) == true, want false", a)
			}
		})
}

func (tt tester) twoElementSetHasLengthOfTwo() {
	tt.t.Helper()

	tt.t.Run(
		"two element set: has length of 2",
		func(t *testing.T) {
			s := tt.setBuilder(twoElements())

			if got, want := s.Len(), 2; got != want {
				t.Fatalf("got Set.Len of %v, want %v", got, want)
			}
		})
}

func (tt tester) emptySetPlusTwoHasLengthOfTwo() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add x2: has length of 2",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()

			s.Add(a)
			s.Add(b)

			if got, want := s.Len(), 2; got != want {
				t.Fatalf("got Set.Len of %v, want %v", got, want)
			}
		})
}

func (tt tester) twoElementSetContainsBothElements() {
	tt.t.Helper()

	tt.t.Run(
		"two element set: contains both elements",
		func(t *testing.T) {
			s := tt.setBuilder(twoElements())

			for _, element := range twoElements() {
				if !s.Contains(element) {
					t.Fatalf(
						"got Set.Contains(%q) == false, want true",
						element,
					)
				}
			}
		})
}

func (tt tester) emptySetPlusTwoContainsBothElements() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add x2: contains both elements",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()

			s.Add(a)
			s.Add(b)

			for _, element := range twoElements() {
				if !s.Contains(element) {
					t.Fatalf(
						"got Set.Contains(%q) == false, want true",
						element,
					)
				}
			}
		})
}

func (tt tester) twoElementSetReturnsBothElementsOnIteration() {
	tt.t.Helper()

	tt.t.Run(
		"two element set: returns both elements on iteration",
		func(t *testing.T) {
			s := tt.setBuilder(twoElements())

			got, want := slices.Collect(s.All()), twoElements()
			if diff := orderagnostic.Diff(got, want); diff != "" {
				t.Errorf("Set.All mismatch (-want +got):\n%s", diff)
			}
		})
}

func (tt tester) emptySetPlusTwoReturnsBothElementsOnIteration() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add x2: returns both elements on iteration",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()

			s.Add(a)
			s.Add(b)

			got, want := slices.Collect(s.All()), twoElements()
			if diff := orderagnostic.Diff(got, want); diff != "" {
				t.Errorf("Set.All mismatch (-want +got):\n%s", diff)
			}
		})
}

func (tt tester) emptySetPlusVarargsReturnsBothElementsOnIteration() {
	tt.runEmptyIfMutable(
		"empty set: add varargs: returns all elements on iteration",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()

			s.Add(a, b)

			got, want := slices.Collect(s.All()), twoElements()
			if diff := orderagnostic.Diff(got, want); diff != "" {
				t.Errorf("Set.All mismatch (-want +got):\n%s", diff)
			}
		})
}

func (tt tester) twoElementSetHasTwoElementStringRepr() {
	tt.t.Helper()

	tt.t.Run(
		"two element set: has two-element string representation",
		func(t *testing.T) {
			s := tt.setBuilder(twoElements())

			if got, wantAny := s.String(), abStringCombinations();
				!slices.Contains(wantAny, s.String()) {
				t.Fatalf("got Set.String of %v, want any of %q", got, wantAny)
			}
		})
}

func (tt tester) emptySetPlusTwoReturnsTwoElementStringRepr() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add x2: has two-element string representation",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()

			s.Add(a)
			s.Add(b)

			if got, wantAny := s.String(), abStringCombinations();
				!slices.Contains(wantAny, s.String()) {
				t.Fatalf("got Set.String of %v, want any of %q", got, wantAny)
			}
		})
}

func (tt tester) emptySetPlusTwoMinusOneHasLengthOfOne() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add x2: remove x1: has length of 1",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()

			s.Add(a)
			s.Add(b)
			s.Remove(a)

			if got, want := s.Len(), 1; got != want {
				t.Fatalf("got Set.Len of %v, want %v", got, want)
			}
		})
}

func (tt tester) emptySetPlusTwoMinusVarargsHasLengthOfZero() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add x2: remove varargs: has length of 0",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()

			s.Add(a)
			s.Add(b)
			s.Remove(a, b)

			if got, want := s.Len(), 0; got != want {
				t.Fatalf("got Set.Len of %v, want %v", got, want)
			}
		})
}

func (tt tester) emptySetPlusThreeContainsAllThreeElements() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add x3: contains all three elements",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()

			s.Add(a)
			s.Add(b)
			s.Add(c)

			for _, element := range threeElements() {
				if !s.Contains(element) {
					t.Fatalf(
						"got Set.Contains(%q) == false, want true",
						element,
					)
				}
			}
		})
}

func (tt tester) threeElementSetContainsAllThreeElements() {
	tt.t.Helper()

	tt.t.Run(
		"three element set: contains all three elements",
		func(t *testing.T) {
			s := tt.setBuilder(threeElements())

			for _, element := range threeElements() {
				if !s.Contains(element) {
					t.Fatalf("got Set.Contains(%q) == false, want true", element)
				}
			}
		},
	)
}

func (tt tester) emptySetPlusThreeHasThreeElementStringRepr() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add x3: has three-element string representation",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()

			s.Add(a)
			s.Add(b)
			s.Add(c)

			if got, wantAny := s.String(), abcStringCombinations();
				!slices.Contains(wantAny, s.String()) {
				t.Fatalf(
					"got Set.String of %v, want any of %q",
					got,
					wantAny,
				)
			}
		})
}

func (tt tester) threeElementSetHasThreeElementStringRepr() {
	tt.t.Helper()

	tt.t.Run("three element set: has three-element string representation",
		func(t *testing.T) {
			s := tt.setBuilder(threeElements())

			if got, wantAny := s.String(), abcStringCombinations();
				!slices.Contains(wantAny, s.String()) {
				t.Fatalf("got Set.String of %v, want any of %q", got, wantAny)
			}
		})
}

func (tt tester) setInitializedFromTwoOfSameElementHasLengthOfOne() {
	tt.t.Helper()

	tt.t.Run("set initialized from two of same element: has length of 1",
		func(t *testing.T) {
			s := tt.setBuilder(twoSameElements())

			if got, want := s.Len(), 1; got != want {
				t.Fatalf("got Set.Len of %v, want %v", got, want)
			}
		})
}

func (tt tester) emptySetPlusSameElementTwiceHasLengthOfOne() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add same element x2: has length of 1",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()

			s.Add(a)
			s.Add(a)

			if got, want := s.Len(), 1; got != want {
				t.Fatalf("got Set.Len of %v, want %v", got, want)
			}
		})
}

func (tt tester) setInitializedFromTwoOfSameElementReturnsOneElementOnIteration() {
	tt.t.Helper()

	tt.t.Run("set initialized from two of same element: returns one element on iteration",
		func(t *testing.T) {
			s := tt.setBuilder(twoSameElements())

			if got, want := slices.Collect(s.All()), oneElement();
				!slices.Equal(got, want) {
				t.Fatalf("got Set.All of %q, want %q", got, want)
			}
		})
}

func (tt tester) emptySetPlusSameElementTwiceReturnsOneElementOnIteration() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add same element x2: returns one element on iteration",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()

			s.Add(a)
			s.Add(a)

			if got, want := slices.Collect(s.All()), oneElement();
				!slices.Equal(got, want) {
				t.Fatalf("got Set.All of %q, want %q", got, want)
			}
		})
}

func (tt tester) emptySetPlusOneReturnsTrue() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add: returns true",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()

			got := s.Add(a)

			if !got {
				t.Fatalf("got Set.Add(%q) == false, want true", a)
			}
		})
}

func (tt tester) emptySetPlusSameElementTwiceReturnsFalse() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add same element x2: returns true",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()

			s.Add(a)
			got := s.Add(a)

			if got {
				t.Fatalf("got Set.Add(%q) == true, want false", a)
			}
		})
}

func (tt tester) emptySetPlusSameElementTwiceThenDifferentOnceReturnsTrue() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add same element x2: add different element: returns true",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()

			s.Add(a)
			s.Add(a)
			got := s.Add(b)

			if !got {
				t.Fatalf("got Set.Add(%q) == false, want true", b)
			}
		})
}

func (tt tester) emptySetPlusOnePlusVarargsReturnsTrue() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add x1: add varargs: returns true",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()

			s.Add(a)
			got := s.Add(b, a)

			if !got {
				t.Fatalf("got Set.Add(%q, %q) == false, want true", b, a)
			}
		})
}

func (tt tester) emptySetMinusOneReturnsFalse() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: remove: returns false",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()

			got := s.Remove(a)

			if got {
				t.Fatalf(
					"got Set.Remove(%q) == true, want false",
					a,
				)
			}
		})
}

func (tt tester) emptySetPlusOneMinusSameElementReturnsTrue() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add: remove same element: returns true",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()

			s.Add(a)
			got := s.Remove(a)

			if !got {
				t.Fatalf("got Set.Remove(%q) == false, want true", a)
			}
		})
}

func (tt tester) emptySetPlusOneMinusSameElementTwiceReturnsFalse() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add: remove same element x2: returns false",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()

			s.Add(a)
			s.Remove(a)
			got := s.Remove(a)

			if got {
				t.Fatalf("got Set.Remove(%q) == true, want false", a)
			}
		})
}

func (tt tester) emptySetPlusOneMinusVarargsReturnsTrue() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add: remove varargs: returns true",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()

			s.Add(a)
			got := s.Remove(b, a)

			if !got {
				t.Fatalf("got Set.Remove(%q) == false, want true", a)
			}
		})
}

const (
	a = "link"
	b = "zelda"
	c = "ganondorf"
)

func empty() []string {
	return nil
}

func oneElement() []string {
	return []string{a}
}

func twoElements() []string {
	return []string{a, b}
}

func threeElements() []string {
	return []string{a, b, c}
}

func twoSameElements() []string {
	return []string{a, a}
}

func aString() string {
	return fmt.Sprintf("[%s]", a)
}

func abStringCombinations() []string {
	return []string{
		fmt.Sprintf("[%s, %s]", a, b),
		fmt.Sprintf("[%s, %s]", b, a),
	}
}

func abcStringCombinations() []string {
	template := "[%s, %s, %s]"
	return []string{
		fmt.Sprintf(template, a, b, c),
		fmt.Sprintf(template, a, c, b),
		fmt.Sprintf(template, b, a, c),
		fmt.Sprintf(template, b, c, a),
		fmt.Sprintf(template, c, a, b),
		fmt.Sprintf(template, c, b, a),
	}
}
