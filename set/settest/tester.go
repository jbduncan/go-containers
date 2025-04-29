package settest

import (
	"iter"
	"testing"

	internalsettest "github.com/jbduncan/go-containers/internal/settest"
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

func TestReadOnly(
	t *testing.T,
	sliceToSet func(elements []int) Set[int],
) {
	tt := newTester(t, sliceToSet)

	tt.emptySetHasLengthOfZero()

	tt.emptySetContainsNothing()

	tt.emptySetIterationDoesNothing()

	tt.emptySetHasEmptyStringRepr()

	tt.oneElementSetHasLengthOfOne()

	tt.oneElementSetContainsPresentElement()

	tt.oneElementSetDoesNotContainAbsentElement()

	tt.oneElementSetReturnsElementOnIteration()

	tt.oneElementSetHasOneElementStringRepr()

	tt.twoElementSetHasLengthOfTwo()

	tt.twoElementSetContainsBothElements()

	tt.twoElementSetReturnsBothElementsOnIteration()

	tt.twoElementSetHasTwoElementStringRepr()

	tt.threeElementSetContainsAllThreeElements()

	tt.threeElementSetHasThreeElementStringRepr()

	tt.setInitializedFromTwoOfSameElementHasLengthOfOne()

	tt.setInitializedFromTwoOfSameElementReturnsOneElementOnIteration()
}

func TestMutable(
	t *testing.T,
	sliceToSet func(elements []int) MutableSet[int],
) {
	tt := newTester(t, func(elements []int) Set[int] {
		return sliceToSet(elements)
	})
	ttt := newMutableTester(t, sliceToSet)

	tt.emptySetHasLengthOfZero()

	tt.emptySetContainsNothing()

	tt.emptySetIterationDoesNothing()

	tt.emptySetHasEmptyStringRepr()

	ttt.emptySetRemoveDoesNothing()

	tt.oneElementSetHasLengthOfOne()
	ttt.emptySetPlusOneHasLengthOfOne()

	tt.oneElementSetContainsPresentElement()
	ttt.emptySetPlusOneContainsPresentElement()

	tt.oneElementSetDoesNotContainAbsentElement()
	ttt.emptySetPlusOneDoesNotContainAbsentElement()

	ttt.emptySetPlusOneMinusOneDoesNotContainAnything()

	tt.oneElementSetReturnsElementOnIteration()
	ttt.emptySetPlusOneReturnsElementOnIteration()

	tt.oneElementSetHasOneElementStringRepr()
	ttt.emptySetPlusOneHasOneElementStringRepr()

	tt.twoElementSetHasLengthOfTwo()
	ttt.emptySetPlusTwoHasLengthOfTwo()

	tt.twoElementSetContainsBothElements()
	ttt.emptySetPlusTwoContainsBothElements()

	tt.twoElementSetReturnsBothElementsOnIteration()
	ttt.emptySetPlusTwoReturnsBothElementsOnIteration()
	ttt.emptySetPlusVarargsReturnsBothElementsOnIteration()

	tt.twoElementSetHasTwoElementStringRepr()
	ttt.emptySetPlusTwoReturnsTwoElementStringRepr()

	ttt.emptySetPlusTwoMinusOneHasLengthOfOne()

	ttt.emptySetPlusTwoMinusVarargsHasLengthOfZero()

	tt.threeElementSetContainsAllThreeElements()
	ttt.emptySetPlusThreeContainsAllThreeElements()

	tt.threeElementSetHasThreeElementStringRepr()
	ttt.emptySetPlusThreeHasThreeElementStringRepr()

	tt.setInitializedFromTwoOfSameElementHasLengthOfOne()
	ttt.emptySetPlusSameElementTwiceHasLengthOfOne()

	tt.setInitializedFromTwoOfSameElementReturnsOneElementOnIteration()
	ttt.emptySetPlusSameElementTwiceReturnsOneElementOnIteration()

	ttt.emptySetPlusOneReturnsTrue()

	ttt.emptySetPlusSameElementTwiceReturnsFalse()

	ttt.emptySetPlusSameElementTwiceThenDifferentOnceReturnsTrue()

	ttt.emptySetPlusOnePlusVarargsReturnsTrue()

	ttt.emptySetMinusOneReturnsFalse()

	ttt.emptySetPlusOneMinusSameElementReturnsTrue()

	ttt.emptySetPlusOneMinusSameElementTwiceReturnsFalse()

	ttt.emptySetPlusOneMinusVarargsReturnsTrue()
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

			testString(t, s, empty())
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

func (tt tester) oneElementSetContainsPresentElement() {
	tt.t.Run(
		"one element set: contains present element",
		func(t *testing.T) {
			s := tt.sliceToSet(oneElement())

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

func (tt tester) oneElementSetReturnsElementOnIteration() {
	tt.t.Run(
		"one element set: returns element on iteration",
		func(t *testing.T) {
			s := tt.sliceToSet(oneElement())

			testAll(t, s, oneElement())
		})
}

func (tt tester) oneElementSetHasOneElementStringRepr() {
	tt.t.Run(
		"one element set: has one-element string representation",
		func(t *testing.T) {
			s := tt.sliceToSet(oneElement())

			testString(t, s, oneElement())
		})
}

func (tt tester) twoElementSetHasLengthOfTwo() {
	tt.t.Run(
		"two element set: has length of 2",
		func(t *testing.T) {
			s := tt.sliceToSet(twoElements())

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

func (tt tester) twoElementSetReturnsBothElementsOnIteration() {
	tt.t.Run(
		"two element set: returns both elements on iteration",
		func(t *testing.T) {
			s := tt.sliceToSet(twoElements())

			testAll(t, s, twoElements())
		})
}

func (tt tester) twoElementSetHasTwoElementStringRepr() {
	tt.t.Run(
		"two element set: has two-element string representation",
		func(t *testing.T) {
			s := tt.sliceToSet(twoElements())

			testString(t, s, twoElements())
		})
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

func (tt tester) threeElementSetHasThreeElementStringRepr() {
	tt.t.Run("three element set: has three-element string representation",
		func(t *testing.T) {
			s := tt.sliceToSet(threeElements())

			testString(t, s, threeElements())
		})
}

func (tt tester) setInitializedFromTwoOfSameElementHasLengthOfOne() {
	tt.t.Run("set initialized from two of same element: has length of 1",
		func(t *testing.T) {
			s := tt.sliceToSet(twoSameElements())

			testLen(t, s, 1)
		})
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

type mutableTester struct {
	t          *testing.T
	sliceToSet func(elements []int) MutableSet[int]
}

func newMutableTester(
	t *testing.T,
	sliceToSet func(elements []int) MutableSet[int],
) *mutableTester {
	return &mutableTester{
		t:          t,
		sliceToSet: sliceToSet,
	}
}

func (tt mutableTester) emptySetRemoveDoesNothing() {
	tt.t.Run("empty set: remove does nothing", func(t *testing.T) {
		s := tt.sliceToSet(empty())

		s.Remove(a)

		testLen(t, s, 0)
	})
}

func (tt mutableTester) emptySetPlusOneHasLengthOfOne() {
	tt.t.Run("empty set: add: has length of 1", func(t *testing.T) {
		s := tt.sliceToSet(empty())

		s.Add(a)

		testLen(t, s, 1)
	})
}

func (tt mutableTester) emptySetPlusOneContainsPresentElement() {
	tt.t.Run("empty set: add: contains present element", func(t *testing.T) {
		s := tt.sliceToSet(empty())

		s.Add(a)

		testContains(t, s, a)
	})
}

func (tt mutableTester) emptySetPlusOneDoesNotContainAbsentElement() {
	tt.t.Run(
		"empty set: add: does not contain absent element",
		func(t *testing.T) {
			s := tt.sliceToSet(empty())

			s.Add(a)

			testDoesNotContain(t, s, b)
		},
	)
}

func (tt mutableTester) emptySetPlusOneReturnsElementOnIteration() {
	tt.t.Run(
		"empty set: add: returns element on iteration",
		func(t *testing.T) {
			s := tt.sliceToSet(empty())

			s.Add(a)

			testAll(t, s, oneElement())
		},
	)
}

func (tt mutableTester) emptySetPlusOneHasOneElementStringRepr() {
	tt.t.Run(
		"empty set: add: has one-element string representation",
		func(t *testing.T) {
			s := tt.sliceToSet(empty())

			s.Add(a)

			testString(t, s, oneElement())
		},
	)
}

func (tt mutableTester) emptySetPlusOneMinusOneDoesNotContainAnything() {
	tt.t.Run(
		"empty set: add: remove: does not contain anything",
		func(t *testing.T) {
			s := tt.sliceToSet(empty())

			s.Add(a)
			s.Remove(a)

			testDoesNotContain(t, s, a)
		},
	)
}

func (tt mutableTester) emptySetPlusTwoHasLengthOfTwo() {
	tt.t.Run("empty set: add x2: has length of 2", func(t *testing.T) {
		s := tt.sliceToSet(empty())

		s.Add(a)
		s.Add(b)

		testLen(t, s, 2)
	})
}

func (tt mutableTester) emptySetPlusTwoContainsBothElements() {
	tt.t.Run("empty set: add x2: contains both elements", func(t *testing.T) {
		s := tt.sliceToSet(empty())

		s.Add(a)
		s.Add(b)

		for _, element := range twoElements() {
			testContains(t, s, element)
		}
	})
}

func (tt mutableTester) emptySetPlusTwoReturnsBothElementsOnIteration() {
	tt.t.Run(
		"empty set: add x2: returns both elements on iteration",
		func(t *testing.T) {
			s := tt.sliceToSet(empty())

			s.Add(a)
			s.Add(b)

			testAll(t, s, twoElements())
		},
	)
}

func (tt mutableTester) emptySetPlusVarargsReturnsBothElementsOnIteration() {
	tt.t.Run(
		"empty set: add varargs: returns all elements on iteration",
		func(t *testing.T) {
			s := tt.sliceToSet(empty())

			s.Add(a, b)

			testAll(t, s, twoElements())
		},
	)
}

func (tt mutableTester) emptySetPlusTwoReturnsTwoElementStringRepr() {
	tt.t.Run(
		"empty set: add x2: has two-element string representation",
		func(t *testing.T) {
			s := tt.sliceToSet(empty())

			s.Add(a)
			s.Add(b)

			testString(t, s, twoElements())
		},
	)
}

func (tt mutableTester) emptySetPlusTwoMinusOneHasLengthOfOne() {
	tt.t.Run(
		"empty set: add x2: remove x1: has length of 1",
		func(t *testing.T) {
			s := tt.sliceToSet(empty())

			s.Add(a)
			s.Add(b)
			s.Remove(a)

			testLen(t, s, 1)
		},
	)
}

func (tt mutableTester) emptySetPlusTwoMinusVarargsHasLengthOfZero() {
	tt.t.Run(
		"empty set: add x2: remove varargs: has length of 0",
		func(t *testing.T) {
			s := tt.sliceToSet(empty())

			s.Add(a)
			s.Add(b)
			s.Remove(a, b)

			testLen(t, s, 0)
		},
	)
}

func (tt mutableTester) emptySetPlusThreeContainsAllThreeElements() {
	tt.t.Run(
		"empty set: add x3: contains all three elements",
		func(t *testing.T) {
			s := tt.sliceToSet(empty())

			s.Add(a)
			s.Add(b)
			s.Add(c)

			for _, element := range threeElements() {
				testContains(t, s, element)
			}
		},
	)
}

func (tt mutableTester) emptySetPlusThreeHasThreeElementStringRepr() {
	tt.t.Run(
		"empty set: add x3: has three-element string representation",
		func(t *testing.T) {
			s := tt.sliceToSet(empty())

			s.Add(a)
			s.Add(b)
			s.Add(c)

			testString(t, s, threeElements())
		},
	)
}

func (tt mutableTester) emptySetPlusSameElementTwiceHasLengthOfOne() {
	tt.t.Run(
		"empty set: add same element x2: has length of 1",
		func(t *testing.T) {
			s := tt.sliceToSet(empty())

			s.Add(a)
			s.Add(a)

			testLen(t, s, 1)
		},
	)
}

func (tt mutableTester) emptySetPlusSameElementTwiceReturnsOneElementOnIteration() {
	tt.t.Run(
		"empty set: add same element x2: returns one element on iteration",
		func(t *testing.T) {
			s := tt.sliceToSet(empty())

			s.Add(a)
			s.Add(a)

			testAll(t, s, oneElement())
		},
	)
}

func (tt mutableTester) emptySetPlusOneReturnsTrue() {
	tt.t.Run("empty set: add: returns true", func(t *testing.T) {
		s := tt.sliceToSet(empty())

		got := s.Add(a)

		if !got {
			t.Fatalf("got Set.Add(%d) == false, want true", a)
		}
	})
}

func (tt mutableTester) emptySetPlusSameElementTwiceReturnsFalse() {
	tt.t.Run("empty set: add same element x2: returns true", func(t *testing.T) {
		s := tt.sliceToSet(empty())

		s.Add(a)
		got := s.Add(a)

		if got {
			t.Fatalf("got Set.Add(%d) == true, want false", a)
		}
	})
}

func (tt mutableTester) emptySetPlusSameElementTwiceThenDifferentOnceReturnsTrue() {
	tt.t.Run(
		"empty set: add same element x2: add different element: returns true",
		func(t *testing.T) {
			s := tt.sliceToSet(empty())

			s.Add(a)
			s.Add(a)
			got := s.Add(b)

			if !got {
				t.Fatalf("got Set.Add(%d) == false, want true", b)
			}
		},
	)
}

func (tt mutableTester) emptySetPlusOnePlusVarargsReturnsTrue() {
	tt.t.Run(
		"empty set: add x1: add varargs: returns true",
		func(t *testing.T) {
			s := tt.sliceToSet(empty())

			s.Add(a)
			got := s.Add(b, a)

			if !got {
				t.Fatalf("got Set.Add(%d, %d) == false, want true", b, a)
			}
		},
	)
}

func (tt mutableTester) emptySetMinusOneReturnsFalse() {
	tt.t.Run("empty set: remove: returns false", func(t *testing.T) {
		s := tt.sliceToSet(empty())

		got := s.Remove(a)

		if got {
			t.Fatalf(
				"got Set.Remove(%d) == true, want false",
				a,
			)
		}
	})
}

func (tt mutableTester) emptySetPlusOneMinusSameElementReturnsTrue() {
	tt.t.Run(
		"empty set: add: remove same element: returns true",
		func(t *testing.T) {
			s := tt.sliceToSet(empty())

			s.Add(a)
			got := s.Remove(a)

			if !got {
				t.Fatalf("got Set.Remove(%d) == false, want true", a)
			}
		},
	)
}

func (tt mutableTester) emptySetPlusOneMinusSameElementTwiceReturnsFalse() {
	tt.t.Run(
		"empty set: add: remove same element x2: returns false",
		func(t *testing.T) {
			s := tt.sliceToSet(empty())

			s.Add(a)
			s.Remove(a)
			got := s.Remove(a)

			if got {
				t.Fatalf("got Set.Remove(%d) == true, want false", a)
			}
		},
	)
}

func (tt mutableTester) emptySetPlusOneMinusVarargsReturnsTrue() {
	tt.t.Run(
		"empty set: add: remove varargs: returns true",
		func(t *testing.T) {
			s := tt.sliceToSet(empty())

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

	internalsettest.Len(t, "", s, want)
}

func testContains(t *testing.T, s Set[int], want int) {
	t.Helper()

	internalsettest.Contains(t, "", s, []int{want})
}

func testDoesNotContain(t *testing.T, s Set[int], want int) {
	t.Helper()

	internalsettest.DoesNotContain(t, "", s, []int{want})
}

func testAll(t *testing.T, s Set[int], want []int) {
	t.Helper()

	internalsettest.All(t, "", s, want)
}

func testString(t *testing.T, s Set[int], want []int) {
	t.Helper()

	internalsettest.String(t, "", s, want)
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
