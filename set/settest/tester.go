package settest

import (
	"fmt"
	"iter"
	"slices"
	"testing"

	"github.com/jbduncan/go-containers/internal/orderagnostic"
)

// Set is a generic, unordered collection of unique elements.
type Set[T comparable] interface {
	// Contains returns true if this set contains the given element, otherwise it returns false.
	Contains(element T) bool

	// Len returns the number of elements in this set.
	Len() int

	// All returns an iter.Seq that returns each and every element in this set.
	//
	// The iteration order is undefined; it may even change from one call to the next.
	All() iter.Seq[T]

	// String returns a string representation of all the elements in this set.
	//
	// The format of this string is a single "[" followed by a comma-separated list (", ") of this set's elements in
	// the same order as All (which is undefined and may change from one call to the next), followed by a single
	// "]".
	//
	// This method satisfies fmt.Stringer.
	String() string
}

// MutableSet is a Set with additional methods for adding and removing elements.
//
// An instance of MutableSet can be made with set.Of.
type MutableSet[T comparable] interface {
	Set[T]

	// Add adds the given element(s) to this set. If any of the elements are already present, the set will not add
	// those elements again. Returns true if this set changed as a result of this call, otherwise false.
	Add(element T, others ...T) bool

	// Remove removes the given element(s) from this set. If any of the elements are already absent, the set will not
	// attempt to remove those elements. Returns true if this set changed as a result of this call, otherwise false.
	Remove(element T, others ...T) bool
}

// TODO: split into TestSet and TestMutableSet

func TestSet(t *testing.T, sliceToSet func(elements []int) Set[int]) {
	tt := newTester(t, sliceToSet)

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
	t          *testing.T
	sliceToSet func(elements []int) Set[int]
}

func newTester(
	t *testing.T,
	sliceToSet func(elements []int) Set[int],
) *tester {
	return &tester{
		t:          t,
		sliceToSet: sliceToSet,
	}
}

func (tt tester) emptySetHasLengthOfZero() {
	tt.t.Run(
		"empty set: has length of 0",
		func(t *testing.T) {
			s := tt.sliceToSet(empty())

			testLen(t, s, 0)
		})
}

func (tt tester) emptySetContainsNothing() {
	tt.t.Run(
		"empty set: contains nothing",
		func(t *testing.T) {
			s := tt.sliceToSet(empty())

			testDoesNotContain(t, s, a)
		})
}

func (tt tester) emptySetIterationDoesNothing() {
	tt.t.Run(
		"empty set: iteration does nothing",
		func(t *testing.T) {
			s := tt.sliceToSet(empty())

			testAll(t, s, empty())
		})
}

func (tt tester) emptySetHasEmptyStringRepr() {
	tt.t.Run(
		"empty set: has empty string representation",
		func(t *testing.T) {
			s := tt.sliceToSet(empty())

			testString(t, s, "[]")
		})
}

func (tt tester) emptySetRemoveDoesNothing() {
	s, mutable := tt.sliceToSet(empty()).(MutableSet[int])
	if !mutable {
		return
	}

	tt.t.Run("empty set: remove does nothing", func(t *testing.T) {
		s.Remove(a)

		testLen(t, s, 0)
	})
}

func (tt tester) oneElementSetHasLengthOfOne() {
	tt.t.Run(
		"one element set: has length of 1",
		func(t *testing.T) {
			s := tt.sliceToSet(oneElement())

			testLen(t, s, 1)
		})
}

func (tt tester) emptySetPlusOneHasLengthOfOne() {
	s, mutable := tt.sliceToSet(empty()).(MutableSet[int])
	if !mutable {
		return
	}

	tt.t.Run("empty set: add: has length of 1", func(t *testing.T) {
		s.Add(a)

		testLen(t, s, 1)
	})
}

func (tt tester) oneElementSetContainsPresentElement() {
	tt.t.Run(
		"one element set: contains present element",
		func(t *testing.T) {
			s := tt.sliceToSet(oneElement())

			testContains(t, s, a)
		})
}

func (tt tester) emptySetPlusOneContainsPresentElement() {
	s, mutable := tt.sliceToSet(empty()).(MutableSet[int])
	if !mutable {
		return
	}

	tt.t.Run("empty set: add: contains present element", func(t *testing.T) {
		s.Add(a)

		testContains(t, s, a)
	})
}

func (tt tester) oneElementSetDoesNotContainAbsentElement() {
	tt.t.Run(
		"one element set: does not contain absent element",
		func(t *testing.T) {
			s := tt.sliceToSet(oneElement())

			testDoesNotContain(t, s, b)
		})
}

func (tt tester) emptySetPlusOneDoesNotContainAbsentElement() {
	s, mutable := tt.sliceToSet(empty()).(MutableSet[int])
	if !mutable {
		return
	}

	tt.t.Run(
		"empty set: add: does not contain absent element",
		func(t *testing.T) {
			s.Add(a)

			testDoesNotContain(t, s, b)
		},
	)
}

func (tt tester) oneElementSetReturnsElementOnIteration() {
	tt.t.Run(
		"one element set: returns element on iteration",
		func(t *testing.T) {
			s := tt.sliceToSet(oneElement())

			testAll(t, s, oneElement())
		})
}

func (tt tester) emptySetPlusOneReturnsElementOnIteration() {
	s, mutable := tt.sliceToSet(empty()).(MutableSet[int])
	if !mutable {
		return
	}

	tt.t.Run(
		"empty set: add: returns element on iteration",
		func(t *testing.T) {
			s.Add(a)

			testAll(t, s, oneElement())
		},
	)
}

func (tt tester) oneElementSetHasOneElementStringRepr() {
	tt.t.Run(
		"one element set: has one-element string representation",
		func(t *testing.T) {
			s := tt.sliceToSet(oneElement())

			testString(t, s, aString())
		})
}

func (tt tester) emptySetPlusOneHasOneElementStringRepr() {
	s, mutable := tt.sliceToSet(empty()).(MutableSet[int])
	if !mutable {
		return
	}

	tt.t.Run(
		"empty set: add: has one-element string representation",
		func(t *testing.T) {
			s.Add(a)

			testString(t, s, aString())
		},
	)
}

func (tt tester) emptySetPlusOneMinusOneDoesNotContainAnything() {
	s, mutable := tt.sliceToSet(empty()).(MutableSet[int])
	if !mutable {
		return
	}

	tt.t.Run(
		"empty set: add: remove: does not contain anything",
		func(t *testing.T) {
			s.Add(a)
			s.Remove(a)

			testDoesNotContain(t, s, a)
		},
	)
}

func (tt tester) twoElementSetHasLengthOfTwo() {
	tt.t.Run(
		"two element set: has length of 2",
		func(t *testing.T) {
			s := tt.sliceToSet(twoElements())

			testLen(t, s, 2)
		})
}

func (tt tester) emptySetPlusTwoHasLengthOfTwo() {
	s, mutable := tt.sliceToSet(empty()).(MutableSet[int])
	if !mutable {
		return
	}

	tt.t.Run("empty set: add x2: has length of 2", func(t *testing.T) {
		s.Add(a)
		s.Add(b)

		testLen(t, s, 2)
	})
}

func (tt tester) twoElementSetContainsBothElements() {
	tt.t.Run(
		"two element set: contains both elements",
		func(t *testing.T) {
			s := tt.sliceToSet(twoElements())

			for _, element := range twoElements() {
				testContains(t, s, element)
			}
		})
}

func (tt tester) emptySetPlusTwoContainsBothElements() {
	s, mutable := tt.sliceToSet(empty()).(MutableSet[int])
	if !mutable {
		return
	}

	tt.t.Run("empty set: add x2: contains both elements", func(t *testing.T) {
		s.Add(a)
		s.Add(b)

		for _, element := range twoElements() {
			testContains(t, s, element)
		}
	})
}

func (tt tester) twoElementSetReturnsBothElementsOnIteration() {
	tt.t.Run(
		"two element set: returns both elements on iteration",
		func(t *testing.T) {
			s := tt.sliceToSet(twoElements())

			testAll(t, s, twoElements())
		})
}

func (tt tester) emptySetPlusTwoReturnsBothElementsOnIteration() {
	s, mutable := tt.sliceToSet(empty()).(MutableSet[int])
	if !mutable {
		return
	}

	tt.t.Run(
		"empty set: add x2: returns both elements on iteration",
		func(t *testing.T) {
			s.Add(a)
			s.Add(b)

			testAll(t, s, twoElements())
		},
	)
}

func (tt tester) emptySetPlusVarargsReturnsBothElementsOnIteration() {
	s, mutable := tt.sliceToSet(empty()).(MutableSet[int])
	if !mutable {
		return
	}

	tt.t.Run(
		"empty set: add varargs: returns all elements on iteration",
		func(t *testing.T) {
			s.Add(a, b)

			testAll(t, s, twoElements())
		},
	)
}

func (tt tester) twoElementSetHasTwoElementStringRepr() {
	tt.t.Run(
		"two element set: has two-element string representation",
		func(t *testing.T) {
			s := tt.sliceToSet(twoElements())

			testStringAnyOf(t, s, abStringCombinations())
		})
}

func (tt tester) emptySetPlusTwoReturnsTwoElementStringRepr() {
	s, mutable := tt.sliceToSet(empty()).(MutableSet[int])
	if !mutable {
		return
	}

	tt.t.Run(
		"empty set: add x2: has two-element string representation",
		func(t *testing.T) {
			s.Add(a)
			s.Add(b)

			testStringAnyOf(t, s, abStringCombinations())
		},
	)
}

func (tt tester) emptySetPlusTwoMinusOneHasLengthOfOne() {
	s, mutable := tt.sliceToSet(empty()).(MutableSet[int])
	if !mutable {
		return
	}

	tt.t.Run(
		"empty set: add x2: remove x1: has length of 1",
		func(t *testing.T) {
			s.Add(a)
			s.Add(b)
			s.Remove(a)

			testLen(t, s, 1)
		},
	)
}

func (tt tester) emptySetPlusTwoMinusVarargsHasLengthOfZero() {
	s, mutable := tt.sliceToSet(empty()).(MutableSet[int])
	if !mutable {
		return
	}

	tt.t.Run(
		"empty set: add x2: remove varargs: has length of 0",
		func(t *testing.T) {
			s.Add(a)
			s.Add(b)
			s.Remove(a, b)

			testLen(t, s, 0)
		},
	)
}

func (tt tester) emptySetPlusThreeContainsAllThreeElements() {
	s, mutable := tt.sliceToSet(empty()).(MutableSet[int])
	if !mutable {
		return
	}

	tt.t.Run(
		"empty set: add x3: contains all three elements",
		func(t *testing.T) {
			s.Add(a)
			s.Add(b)
			s.Add(c)

			for _, element := range threeElements() {
				testContains(t, s, element)
			}
		},
	)
}

func (tt tester) threeElementSetContainsAllThreeElements() {
	tt.t.Run(
		"three element set: contains all three elements",
		func(t *testing.T) {
			s := tt.sliceToSet(threeElements())

			for _, element := range threeElements() {
				testContains(t, s, element)
			}
		},
	)
}

func (tt tester) emptySetPlusThreeHasThreeElementStringRepr() {
	s, mutable := tt.sliceToSet(empty()).(MutableSet[int])
	if !mutable {
		return
	}

	tt.t.Run(
		"empty set: add x3: has three-element string representation",
		func(t *testing.T) {
			s.Add(a)
			s.Add(b)
			s.Add(c)

			testStringAnyOf(t, s, abcStringCombinations())
		},
	)
}

func (tt tester) threeElementSetHasThreeElementStringRepr() {
	tt.t.Run("three element set: has three-element string representation",
		func(t *testing.T) {
			s := tt.sliceToSet(threeElements())

			testStringAnyOf(t, s, abcStringCombinations())
		})
}

func (tt tester) setInitializedFromTwoOfSameElementHasLengthOfOne() {
	tt.t.Run("set initialized from two of same element: has length of 1",
		func(t *testing.T) {
			s := tt.sliceToSet(twoSameElements())

			testLen(t, s, 1)
		})
}

func (tt tester) emptySetPlusSameElementTwiceHasLengthOfOne() {
	s, mutable := tt.sliceToSet(empty()).(MutableSet[int])
	if !mutable {
		return
	}

	tt.t.Run(
		"empty set: add same element x2: has length of 1",
		func(t *testing.T) {
			s.Add(a)
			s.Add(a)

			testLen(t, s, 1)
		},
	)
}

func (tt tester) setInitializedFromTwoOfSameElementReturnsOneElementOnIteration() {
	tt.t.Run(
		"set initialized from two of same element: returns one element on iteration",
		func(t *testing.T) {
			s := tt.sliceToSet(twoSameElements())

			testAll(t, s, oneElement())
		},
	)
}

func (tt tester) emptySetPlusSameElementTwiceReturnsOneElementOnIteration() {
	s, mutable := tt.sliceToSet(empty()).(MutableSet[int])
	if !mutable {
		return
	}

	tt.t.Run(
		"empty set: add same element x2: returns one element on iteration",
		func(t *testing.T) {
			s.Add(a)
			s.Add(a)

			testAll(t, s, oneElement())
		},
	)
}

// TODO: extract testX functions for Set.Add/Remove

func (tt tester) emptySetPlusOneReturnsTrue() {
	s, mutable := tt.sliceToSet(empty()).(MutableSet[int])
	if !mutable {
		return
	}

	tt.t.Run("empty set: add: returns true", func(t *testing.T) {
		got := s.Add(a)

		if !got {
			t.Fatalf("got Set.Add(%d) == false, want true", a)
		}
	})
}

func (tt tester) emptySetPlusSameElementTwiceReturnsFalse() {
	s, mutable := tt.sliceToSet(empty()).(MutableSet[int])
	if !mutable {
		return
	}

	tt.t.Run("empty set: add same element x2: returns true", func(t *testing.T) {
		s.Add(a)
		got := s.Add(a)

		if got {
			t.Fatalf("got Set.Add(%d) == true, want false", a)
		}
	})
}

func (tt tester) emptySetPlusSameElementTwiceThenDifferentOnceReturnsTrue() {
	s, mutable := tt.sliceToSet(empty()).(MutableSet[int])
	if !mutable {
		return
	}

	tt.t.Run(
		"empty set: add same element x2: add different element: returns true",
		func(t *testing.T) {
			s.Add(a)
			s.Add(a)
			got := s.Add(b)

			if !got {
				t.Fatalf("got Set.Add(%d) == false, want true", b)
			}
		},
	)
}

func (tt tester) emptySetPlusOnePlusVarargsReturnsTrue() {
	s, mutable := tt.sliceToSet(empty()).(MutableSet[int])
	if !mutable {
		return
	}

	tt.t.Run(
		"empty set: add x1: add varargs: returns true",
		func(t *testing.T) {
			s.Add(a)
			got := s.Add(b, a)

			if !got {
				t.Fatalf("got Set.Add(%d, %d) == false, want true", b, a)
			}
		},
	)
}

func (tt tester) emptySetMinusOneReturnsFalse() {
	s, mutable := tt.sliceToSet(empty()).(MutableSet[int])
	if !mutable {
		return
	}

	tt.t.Run("empty set: remove: returns false", func(t *testing.T) {
		got := s.Remove(a)

		if got {
			t.Fatalf(
				"got Set.Remove(%d) == true, want false",
				a,
			)
		}
	})
}

func (tt tester) emptySetPlusOneMinusSameElementReturnsTrue() {
	s, mutable := tt.sliceToSet(empty()).(MutableSet[int])
	if !mutable {
		return
	}

	tt.t.Run(
		"empty set: add: remove same element: returns true",
		func(t *testing.T) {
			s.Add(a)
			got := s.Remove(a)

			if !got {
				t.Fatalf("got Set.Remove(%d) == false, want true", a)
			}
		},
	)
}

func (tt tester) emptySetPlusOneMinusSameElementTwiceReturnsFalse() {
	s, mutable := tt.sliceToSet(empty()).(MutableSet[int])
	if !mutable {
		return
	}

	tt.t.Run(
		"empty set: add: remove same element x2: returns false",
		func(t *testing.T) {
			s.Add(a)
			s.Remove(a)
			got := s.Remove(a)

			if got {
				t.Fatalf("got Set.Remove(%d) == true, want false", a)
			}
		},
	)
}

func (tt tester) emptySetPlusOneMinusVarargsReturnsTrue() {
	s, mutable := tt.sliceToSet(empty()).(MutableSet[int])
	if !mutable {
		return
	}

	tt.t.Run(
		"empty set: add: remove varargs: returns true",
		func(t *testing.T) {
			s.Add(a)
			got := s.Remove(b, a)

			if !got {
				t.Fatalf("got Set.Remove(%d) == false, want true", a)
			}
		},
	)
}

func testLen(t *testing.T, s Set[int], want int) {
	t.Helper()

	if got := s.Len(); got != want {
		t.Errorf("got Set.Len of %d, want %d", got, want)
	}
}

func testContains(t *testing.T, s Set[int], want int) {
	t.Helper()

	if !s.Contains(want) {
		t.Errorf("got Set.Contains(%d) == false, want true", want)
	}
}

func testDoesNotContain(t *testing.T, s Set[int], want int) {
	t.Helper()

	if s.Contains(want) {
		t.Errorf("got Set.Contains(%d) == true, want false", want)
	}
}

func testAll(t *testing.T, s Set[int], want []int) {
	t.Helper()

	got := slices.Collect(s.All())
	if diff := orderagnostic.Diff(got, want); diff != "" {
		t.Errorf("Set.All mismatch (-want +got):\n%s", diff)
	}
}

func testString(t *testing.T, s Set[int], want string) {
	t.Helper()

	if got := s.String(); got != want {
		t.Errorf("got Set.String of %q, want %q", got, want)
	}
}

func testStringAnyOf(t *testing.T, s Set[int], wantAny []string) {
	t.Helper()

	if got := s.String(); !slices.Contains(wantAny, s.String()) {
		t.Errorf("got Set.String of %q, want any of %q", got, wantAny)
	}
}

const (
	a = 1
	b = 2
	c = 3
)

func empty() []int {
	return nil
}

func oneElement() []int {
	return []int{a}
}

func twoElements() []int {
	return []int{a, b}
}

func threeElements() []int {
	return []int{a, b, c}
}

func twoSameElements() []int {
	return []int{a, a}
}

func aString() string {
	return fmt.Sprintf("[%d]", a)
}

func abStringCombinations() []string {
	return []string{
		fmt.Sprintf("[%d, %d]", a, b),
		fmt.Sprintf("[%d, %d]", b, a),
	}
}

func abcStringCombinations() []string {
	template := "[%d, %d, %d]"
	return []string{
		fmt.Sprintf(template, a, b, c),
		fmt.Sprintf(template, a, c, b),
		fmt.Sprintf(template, b, a, c),
		fmt.Sprintf(template, b, c, a),
		fmt.Sprintf(template, c, a, b),
		fmt.Sprintf(template, c, b, a),
	}
}
