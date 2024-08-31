package set_test

import (
	"fmt"

	"github.com/jbduncan/go-containers/set"
)

func ExampleOf() {
	// Create a new set and put some strings in it.
	exampleSet := set.Of[string]()
	added := exampleSet.Add("link")
	fmt.Println(added) // true
	addedAgain := exampleSet.Add("link")
	fmt.Println(addedAgain) // false
	exampleSet.Add("zelda", "ganondorf")

	// Check that the set contains everything added to it.
	fmt.Println(exampleSet.Contains("link"))      // true
	fmt.Println(exampleSet.Contains("zelda"))     // true
	fmt.Println(exampleSet.Contains("ganondorf")) // true
	fmt.Println(exampleSet.Contains("mario"))     // false
	fmt.Println(exampleSet.Len())                 // 3

	// Remove strings from the set.
	exampleSet.Remove("zelda")
	exampleSet.Remove("ganondorf")
	fmt.Println(exampleSet.Contains("zelda"))     // false
	fmt.Println(exampleSet.Contains("ganondorf")) // false
	fmt.Println(exampleSet.String())              // [link]

	// Loop over all elements in the set.
	for elem := range exampleSet.All() {
		fmt.Println(elem) // link
	}

	// Output:
	// true
	// false
	// true
	// true
	// true
	// false
	// 3
	// false
	// false
	// [link]
	// link
}

func ExampleUnmodifiable() {
	underlyingSet := set.Of[string]()
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
	a := set.Of[string]()
	a.Add("link")
	b := set.Of("link")
	fmt.Println(set.Equal[string](a, b)) // true

	c := set.Of("ganondorf")
	fmt.Println(set.Equal[string](a, c)) // false

	// Output:
	// true
	// false
}
