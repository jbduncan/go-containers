package set_test

import (
	"fmt"

	"github.com/jbduncan/go-containers/set"
)

func ExampleSet() {
	// Create a new mutable set and put some strings in it.
	exampleSet := set.New[string]()
	added := exampleSet.Add("link")
	fmt.Println(added) // true
	addedAgain := exampleSet.Add("link")
	fmt.Println(addedAgain) // false
	exampleSet.Add("zelda")

	// Check that the set contains everything added to it...
	fmt.Println(exampleSet.Contains("link"))  // true
	fmt.Println(exampleSet.Contains("zelda")) // true
	// ...and that it doesn't contain anything that wasn't added to it.
	fmt.Println(exampleSet.Contains("ganondorf")) // false

	// Remove a string from the set.
	exampleSet.Remove("zelda")
	fmt.Println(exampleSet.Contains("zelda")) // false

	// Print the contents of the set
	fmt.Println(exampleSet.String()) // [link]

	// Check if the contents are equal to another set's contents
	anotherSet := set.New[string]()
	anotherSet.Add("link")
	fmt.Println(set.Equal[string](exampleSet, anotherSet)) // true

	yetAnotherSet := set.New[string]()
	yetAnotherSet.Add("ganondorf")
	fmt.Println(set.Equal[string](exampleSet, yetAnotherSet)) // false

	// Output:
	// true
	// false
	// true
	// true
	// false
	// false
	// [link]
	// true
	// false
}

func ExampleUnmodifiable() {
	underlyingSet := set.New[string]()
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
