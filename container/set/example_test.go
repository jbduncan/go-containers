package set_test

import (
	"fmt"
	"go-containers/container/set"
)

func Example() {
	// Create a new set and put some strings in it.
	s := set.New[string]()
	s.Add("link")
	s.Add("zelda")

	fmt.Println(s.Contains("link"))  // true
	fmt.Println(s.Contains("zelda")) // true
	fmt.Println(s.Contains("epona")) // false

	// Remove a string from the set.
	s.Remove("zelda")
	fmt.Println(s.Contains("zelda")) // false

	// Make an unmodifiable set that wraps s.
	u := set.Unmodifiable[string](s)
	fmt.Println(u.Contains("link"))  // true
	fmt.Println(u.Contains("zelda")) // false

	// Add an element back to s; this also adds it to u
	s.Add("zelda")
	fmt.Println(u.Contains("link"))  // true
	fmt.Println(u.Contains("zelda")) // true

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
