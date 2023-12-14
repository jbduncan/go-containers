package set_test

import (
	"fmt"

	"github.com/jbduncan/go-containers/set"
)

func ExampleNew() {
	// Create a new mutable set and put some strings in it.
	exampleSet := set.New[string]()
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

	// Check if it has the same elements as another set in any order.
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
	// 2
	// false
	// [link]
	// link
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

func ExampleEqual() {
	a := set.New[string]()
	a.Add("link")

	b := set.New[string]()
	b.Add("link")
	fmt.Println(set.Equal[string](a, b)) // true

	c := set.New[string]()
	c.Add("ganondorf")
	fmt.Println(set.Equal[string](a, c)) // false

	// Output:
	// true
	// false
}
