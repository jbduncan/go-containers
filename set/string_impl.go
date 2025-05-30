package set

import (
	"fmt"
	"iter"
	"strings"
)

// StringImpl is a helper function for implementors of Set.String.
//
// This function is implemented in terms of Set.All, so the order of the
// elements in the returned string is undefined; it may even change from one
// call of StringImpl to the next.
//
// Note: Go needs the generic type to be defined explicitly, like:
//
//	a := set.Of(1)
//	b := set.Of(2)
//	s := set.StringImpl[int](a, b)
//	                   ^^^^^
func StringImpl[T comparable](s interface {
	All() iter.Seq[T]
},
) string {
	var builder strings.Builder

	builder.WriteRune('[')
	index := 0
	for element := range s.All() {
		if index > 0 {
			builder.WriteString(", ")
		}

		builder.WriteString(fmt.Sprintf("%v", element))
		index++
	}

	builder.WriteRune(']')
	return builder.String()
}
