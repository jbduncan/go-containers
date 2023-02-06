package iter_test

import (
	. "github.com/onsi/ginkgo/v2"
)

var _ = Describe("Iter", func() {
	// TODO:
	//  IteratorTester:
	//  - Accept parameters for:
	//    - A function returning a new iterator to be tested
	//    - A slice of expected elements (to be copied for defensive
	//      programming reasons)
	//  - Recursively call all permutations of
	//    max(5, len(expectedElements) + 1) iterator "operations". E.g.:
	//    - "Next", "Next", "Next", "Next", "Next"
	//    - "Next", "Next", "Next", "Next", "Value"
	//    - "Next", "Next", "Next", "Value", "Next"
	//    - "Next", "Next", "Next", "Value", "Value"
	//    - ...
	//    - "Value", "Value", "Value", "Value", "Next"
	//    - "Value", "Value", "Value", "Value", "Value"
	//  - For each permutation:
	//    - Ask the function for a new iterator
	//    - For each "operation":
	//      - Keep track of a list of "remaining" elements.
	//      - Initialise "remaining" with the slice of expected elements.
	//      - Run the "operation" on the iterator.
	//      - If the operation was "Next", assert that it returns true
	//        if "remaining" isn't empty yet, otherwise false.
	//      - If the operation was "Value":
	//          - If "remaining" isn't empty, assert that calling
	//            iterator.Value() returns an element from "remaining",
	//            then remove that element from "remaining". Otherwise,
	//            fail the test with a message containing the current
	//            permutation of operations for future debugging.
	//          - If "remaining" is empty, assert that iterator.Value()
	//            panics.
	//          - (We need to do all this because iterators aren't
	//            guaranteed to have a known order, such as when iterating
	//            over sets backed by Go maps.)
})
