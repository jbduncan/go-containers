package set_test

import (
	"fmt"

	"github.com/jbduncan/go-containers/set"
)

func Example() {
	// Create a new mutable set and put some strings in it.
	exampleSet := set.New[string]()
	exampleSet.Add("link")
	exampleSet.Add("zelda")

	fmt.Println(exampleSet.Contains("link"))      // true
	fmt.Println(exampleSet.Contains("zelda"))     // true
	fmt.Println(exampleSet.Contains("ganondorf")) // false

	// Remove a string from the set.
	exampleSet.Remove("zelda")
	fmt.Println(exampleSet.Contains("zelda")) // false

	// Make an unmodifiable set that wraps set.
	unmodifiableSet := set.Unmodifiable[string](exampleSet)
	fmt.Println(unmodifiableSet.Contains("link"))  // true
	fmt.Println(unmodifiableSet.Contains("zelda")) // false

	// Add an element back to set; this also adds it to unmodifiable set.
	exampleSet.Add("zelda")
	fmt.Println(unmodifiableSet.Contains("link"))  // true
	fmt.Println(unmodifiableSet.Contains("zelda")) // true

	// Output:
	// true
	// true
	// false
	// false
	// true
	// false
	// true
	// true
}
