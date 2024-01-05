package settest

import (
	"testing"

	// dot importing gomega matchers is best practice and this package is used by test code only
	. "github.com/jbduncan/go-containers/internal/matchers" //nolint:stylecheck
	"github.com/jbduncan/go-containers/set"
	. "github.com/onsi/gomega" //nolint:stylecheck
)

// TestingT is an interface for the parts of *testing.T that settest.Set needs
// to run. Whenever you see this interface being used, pass in an instance of
// *testing.T or your unit testing framework's equivalent.
type TestingT interface {
	Helper()
	Run(name string, f func(t *testing.T)) bool
}

// TODO: Document; make a note of the fact that the returned set can have mutation methods...
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

func newTester(t TestingT, setBuilder func(elems []string) set.Set[string]) *tester {
	_, mutable := setBuilder(empty()).(set.MutableSet[string])
	return &tester{
		t:          t,
		setBuilder: setBuilder,
		mutable:    mutable,
	}
}

func (tt tester) runEmptyIfMutable(name string, f func(t *testing.T, s set.MutableSet[string])) {
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
			g := NewWithT(t)
			s := tt.setBuilder(empty())

			g.Expect(s).To(HaveLenOfZero())
		})
}

func (tt tester) emptySetContainsNothing() {
	tt.t.Helper()

	tt.t.Run(
		"empty set: contains nothing",
		func(t *testing.T) {
			g := NewWithT(t)
			s := tt.setBuilder(empty())

			g.Expect(s).ToNot(Contain("link"))
		})
}

func (tt tester) emptySetIterationDoesNothing() {
	tt.t.Helper()

	tt.t.Run(
		"empty set: iteration does nothing",
		func(t *testing.T) {
			g := NewWithT(t)
			s := tt.setBuilder(empty())

			g.Expect(s).To(HaveForEachThatEmitsNothing[string]())
		})
}

func (tt tester) emptySetHasEmptyStringRepr() {
	tt.t.Helper()

	tt.t.Run(
		"empty set: has empty string representation",
		func(t *testing.T) {
			g := NewWithT(t)
			s := tt.setBuilder(empty())

			g.Expect(s).To(HaveStringRepr("[]"))
		})
}

func (tt tester) emptySetRemoveDoesNothing() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: remove does nothing",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()
			g := NewWithT(t)

			s.Remove("link")

			g.Expect(s).To(HaveLenOfZero())
		})
}

func (tt tester) oneElementSetHasLengthOfOne() {
	tt.t.Helper()

	tt.t.Run(
		"one element set: has length of 1",
		func(t *testing.T) {
			g := NewWithT(t)
			s := tt.setBuilder(oneElement())

			g.Expect(s).To(HaveLenOf(1))
		})
}

func (tt tester) emptySetPlusOneHasLengthOfOne() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add: has length of 1",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()
			g := NewWithT(t)

			s.Add("link")

			g.Expect(s).To(HaveLenOf(1))
		})
}

func (tt tester) oneElementSetContainsPresentElement() {
	tt.t.Helper()

	tt.t.Run(
		"one element set: contains present element",
		func(t *testing.T) {
			g := NewWithT(t)
			s := tt.setBuilder(oneElement())

			g.Expect(s).To(Contain("link"))
		})
}

func (tt tester) emptySetPlusOneContainsPresentElement() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add: contains present element",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()
			g := NewWithT(t)

			s.Add("link")

			g.Expect(s).To(Contain("link"))
		})
}

func (tt tester) oneElementSetDoesNotContainAbsentElement() {
	tt.t.Helper()

	tt.t.Run(
		"one element set: does not contain absent element",
		func(t *testing.T) {
			g := NewWithT(t)
			s := tt.setBuilder(oneElement())

			g.Expect(s).ToNot(Contain("zelda"))
		})
}

func (tt tester) emptySetPlusOneDoesNotContainAbsentElement() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add: does not contain absent element",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()
			g := NewWithT(t)

			s.Add("link")

			g.Expect(s).ToNot(Contain("zelda"))
		})
}

func (tt tester) oneElementSetReturnsElementOnIteration() {
	tt.t.Helper()

	tt.t.Run(
		"one element set: returns element on iteration",
		func(t *testing.T) {
			g := NewWithT(t)
			s := tt.setBuilder(oneElement())

			g.Expect(s).To(HaveForEachThatConsistsOf[string]("link"))
		})
}

func (tt tester) emptySetPlusOneReturnsElementOnIteration() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add: returns element on iteration",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()
			g := NewWithT(t)

			s.Add("link")

			g.Expect(s).To(HaveForEachThatConsistsOf[string]("link"))
		})
}

func (tt tester) oneElementSetHasOneElementStringRepr() {
	tt.t.Helper()

	tt.t.Run(
		"one element set: has one-element string representation",
		func(t *testing.T) {
			g := NewWithT(t)
			s := tt.setBuilder(oneElement())

			g.Expect(s).To(HaveStringRepr("[link]"))
		})
}

func (tt tester) emptySetPlusOneHasOneElementStringRepr() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add: has one-element string representation",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()
			g := NewWithT(t)

			s.Add("link")

			g.Expect(s).To(HaveStringRepr("[link]"))
		})
}

func (tt tester) emptySetPlusOneMinusOneDoesNotContainAnything() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add: remove: does not contain anything",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()
			g := NewWithT(t)

			s.Add("link")
			s.Remove("link")

			g.Expect(s).ToNot(Contain("link"))
		})
}

func (tt tester) twoElementSetHasLengthOfTwo() {
	tt.t.Helper()

	tt.t.Run(
		"two element set: has length of 2",
		func(t *testing.T) {
			g := NewWithT(t)
			s := tt.setBuilder(twoElements())

			g.Expect(s).To(HaveLenOf(2))
		})
}

func (tt tester) emptySetPlusTwoHasLengthOfTwo() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add x2: has length of 2",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()
			g := NewWithT(t)

			s.Add("link")
			s.Add("zelda")

			g.Expect(s).To(HaveLenOf(2))
		})
}

func (tt tester) twoElementSetContainsBothElements() {
	tt.t.Helper()

	tt.t.Run(
		"two element set: contains both elements",
		func(t *testing.T) {
			g := NewWithT(t)
			s := tt.setBuilder(twoElements())

			g.Expect(s).To(ContainAtLeast("link", "zelda"))
		})
}

func (tt tester) emptySetPlusTwoContainsBothElements() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add x2: contains both elements",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()
			g := NewWithT(t)

			s.Add("link")
			s.Add("zelda")

			g.Expect(s).To(ContainAtLeast("link", "zelda"))
		})
}

func (tt tester) twoElementSetReturnsBothElementsOnIteration() {
	tt.t.Helper()

	tt.t.Run(
		"two element set: returns both elements on iteration",
		func(t *testing.T) {
			g := NewWithT(t)
			s := tt.setBuilder(twoElements())

			g.Expect(s).To(HaveForEachThatConsistsOf[string]("link", "zelda"))
		})
}

func (tt tester) emptySetPlusTwoReturnsBothElementsOnIteration() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add x2: returns both elements on iteration",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()
			g := NewWithT(t)

			s.Add("link")
			s.Add("zelda")

			g.Expect(s).To(
				HaveForEachThatConsistsOf[string]("link", "zelda"))
		})
}

func (tt tester) emptySetPlusVarargsReturnsBothElementsOnIteration() {
	tt.runEmptyIfMutable(
		"empty set: add varargs: returns all elements on iteration",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()
			g := NewWithT(t)

			s.Add("link", "zelda")

			g.Expect(s).To(
				HaveForEachThatConsistsOf[string]("link", "zelda"))
		})
}

func (tt tester) twoElementSetHasTwoElementStringRepr() {
	tt.t.Helper()

	tt.t.Run(
		"two element set: has two-element string representation",
		func(t *testing.T) {
			g := NewWithT(t)
			s := tt.setBuilder(twoElements())

			g.Expect(s).To(
				HaveStringReprThatIsAnyOf("[link, zelda]", "[zelda, link]"))
		})
}

func (tt tester) emptySetPlusTwoReturnsTwoElementStringRepr() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add x2: has two-element string representation",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()
			g := NewWithT(t)

			s.Add("link")
			s.Add("zelda")

			g.Expect(s).To(
				HaveStringReprThatIsAnyOf("[link, zelda]", "[zelda, link]"))
		})
}

func (tt tester) emptySetPlusTwoMinusOneHasLengthOfOne() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add x2: remove x1: has length of 1",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()
			g := NewWithT(t)

			s.Add("link")
			s.Add("zelda")
			s.Remove("link")

			g.Expect(s).To(HaveLenOf(1))
		})
}

func (tt tester) emptySetPlusTwoMinusVarargsHasLengthOfZero() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add x2: remove varargs: has length of 0",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()
			g := NewWithT(t)

			s.Add("link")
			s.Add("zelda")
			s.Remove("link", "zelda")

			g.Expect(s).To(HaveLenOf(0))
		})
}

func (tt tester) emptySetPlusThreeContainsAllThreeElements() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add x3: contains all three elements",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()
			g := NewWithT(t)

			s.Add("link")
			s.Add("zelda")
			s.Add("ganondorf")

			g.Expect(s).To(
				ContainAtLeast("link", "zelda", "ganondorf"))
		})
}

func (tt tester) threeElementSetContainsAllThreeElements() {
	tt.t.Helper()

	tt.t.Run("three element set: contains all three elements", func(t *testing.T) {
		g := NewWithT(t)
		s := tt.setBuilder(threeElements())

		g.Expect(s).To(
			ContainAtLeast("link", "zelda", "ganondorf"))
	})
}

func (tt tester) emptySetPlusThreeHasThreeElementStringRepr() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add x3: has three-element string representation",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()
			g := NewWithT(t)

			s.Add("link")
			s.Add("zelda")
			s.Add("ganondorf")

			g.Expect(s).To(
				HaveStringReprThatIsAnyOf(
					"[link, zelda, ganondorf]",
					"[link, ganondorf, zelda]",
					"[zelda, link, ganondorf]",
					"[zelda, ganondorf, link]",
					"[ganondorf, link, zelda]",
					"[ganondorf, zelda, link]"))
		})
}

func (tt tester) threeElementSetHasThreeElementStringRepr() {
	tt.t.Helper()

	tt.t.Run("three element set: has three-element string representation",
		func(t *testing.T) {
			g := NewWithT(t)
			s := tt.setBuilder(threeElements())

			g.Expect(s).To(
				HaveStringReprThatIsAnyOf(
					"[link, zelda, ganondorf]",
					"[link, ganondorf, zelda]",
					"[zelda, link, ganondorf]",
					"[zelda, ganondorf, link]",
					"[ganondorf, link, zelda]",
					"[ganondorf, zelda, link]"))
		})
}

func (tt tester) setInitializedFromTwoOfSameElementHasLengthOfOne() {
	tt.t.Helper()

	tt.t.Run("set initialized from two of same element: has length of 1",
		func(t *testing.T) {
			g := NewWithT(t)
			s := tt.setBuilder(twoSameElements())

			g.Expect(s).To(HaveLenOf(1))
		})
}

func (tt tester) emptySetPlusSameElementTwiceHasLengthOfOne() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add same element x2: has length of 1",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()
			g := NewWithT(t)

			s.Add("link")
			s.Add("link")

			g.Expect(s).To(HaveLenOf(1))
		})
}

func (tt tester) emptySetPlusOneReturnsTrue() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add: returns true",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()
			g := NewWithT(t)

			result := s.Add("link")

			g.Expect(result).To(BeTrue())
		})
}

func (tt tester) emptySetPlusSameElementTwiceReturnsFalse() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add same element x2: returns true",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()
			g := NewWithT(t)

			s.Add("link")
			result := s.Add("link")

			g.Expect(result).To(BeFalse())
		})
}

func (tt tester) emptySetPlusSameElementTwiceThenDifferentOnceReturnsTrue() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add same element x2: add different element: returns true",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()
			g := NewWithT(t)

			s.Add("link")
			s.Add("link")
			result := s.Add("zelda")

			g.Expect(result).To(BeTrue())
		})
}

func (tt tester) emptySetPlusOnePlusVarargsReturnsTrue() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add x1: add varargs: returns true",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()
			g := NewWithT(t)

			s.Add("link")
			result := s.Add("zelda", "link")

			g.Expect(result).To(BeTrue())
		})
}

func (tt tester) emptySetMinusOneReturnsFalse() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: remove: returns false",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()
			g := NewWithT(t)

			result := s.Remove("link")

			g.Expect(result).To(BeFalse())
		})
}

func (tt tester) emptySetPlusOneMinusSameElementReturnsTrue() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add: remove same element: returns true",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()
			g := NewWithT(t)

			s.Add("link")
			result := s.Remove("link")

			g.Expect(result).To(BeTrue())
		})
}

func (tt tester) emptySetPlusOneMinusSameElementTwiceReturnsFalse() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add: remove same element x2: returns false",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()
			g := NewWithT(t)

			s.Add("link")
			s.Remove("link")
			result := s.Remove("link")

			g.Expect(result).To(BeFalse())
		})
}

func (tt tester) emptySetPlusOneMinusVarargsReturnsTrue() {
	tt.t.Helper()

	tt.runEmptyIfMutable(
		"empty set: add: remove varargs: returns true",
		func(t *testing.T, s set.MutableSet[string]) {
			t.Helper()
			g := NewWithT(t)

			s.Add("link")
			result := s.Remove("zelda", "link")

			g.Expect(result).To(BeTrue())
		})
}

func empty() []string {
	return make([]string, 0)
}

func oneElement() []string {
	return []string{"link"}
}

func twoElements() []string {
	return []string{"link", "zelda"}
}

func threeElements() []string {
	return []string{"link", "zelda", "ganondorf"}
}

func twoSameElements() []string {
	return []string{"link", "link"}
}
