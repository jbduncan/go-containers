package set_test

import (
	"fmt"

	"github.com/jbduncan/go-containers/set"
)

func ExampleNewMutable() {
	// Create a new mutable set and put some strings in it.
	exampleSet := set.NewMutable[string]()
	added := exampleSet.Add("link")
	fmt.Println(added) // true
	addedAgain := exampleSet.Add("link")
	fmt.Println(addedAgain) // false
	exampleSet.Add("zelda")

	// Check that the set contains everything added to it.
	fmt.Println(exampleSet.Contains("link"))      // true
	fmt.Println(exampleSet.Contains("zelda"))     // true
	fmt.Println(exampleSet.Contains("ganondorf")) // false
	fmt.Println(exampleSet.Len())                 // 2

	// Remove a string from the set.
	exampleSet.Remove("zelda")
	fmt.Println(exampleSet.Contains("zelda")) // false
	fmt.Println(exampleSet.String())          // [link]

	// Loop over all elements in the set.
	exampleSet.ForEach(func(elem string) {
		fmt.Println(elem) // link
	})

	// Output:
	// true
	// false
	// true
	// true
	// false
	// 2
	// false
	// [link]
	// link
}

func ExampleOf() {
	// Create a new immutable set with some strings.
	// Note that it doesn't implement set.MutableSet, so new elements cannot be added.
	exampleSet := set.Of("link", "zelda")

	fmt.Println(exampleSet.Contains("link"))      // true
	fmt.Println(exampleSet.Contains("ganondorf")) // false
}

func ExampleUnmodifiable() {
	underlyingSet := set.NewMutable[string]()
	underlyingSet.Add("link")

	// Make an unmodifiable set that wraps underlyingSet.
	unmodifiableSet := set.Unmodifiable[string](underlyingSet)
	fmt.Println(unmodifiableSet.Contains("link"))  // true
	fmt.Println(unmodifiableSet.Contains("zelda")) // false

	// Add an element back to underlyingSet.
	// This also adds it to unmodifiable set.
	underlyingSet.Add("zelda")
	fmt.Println(unmodifiableSet.Contains("link"))  // true
	fmt.Println(unmodifiableSet.Contains("zelda")) // true

	// Output:
	// true
	// false
	// true
	// true
}

func ExampleEqual() {
	// Check if these two sets have the same elements in any order.
	a := set.NewMutable[string]()
	a.Add("link")

	b := set.Of("link")
	fmt.Println(set.Equal[string](a, b)) // true

	c := set.Of("ganondorf")
	fmt.Println(set.Equal[string](a, c)) // false

	// Output:
	// true
	// false
}
