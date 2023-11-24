package set

import (
	"fmt"
	"strings"
)

// StringImpl is a helper function for implementors of Set.String.
//
// This function is implemented in terms of Set.ForEach, so the order of the
// elements in the returned string is undefined; it may even change from one
// call of StringImpl to the next.
func StringImpl[T comparable](s Set[T]) string {
	var builder strings.Builder

	builder.WriteRune('[')
	index := 0
	s.ForEach(func(elem T) {
		if index > 0 {
			builder.WriteString(", ")
		}

		builder.WriteString(fmt.Sprintf("%v", elem))
		index++
	})

	builder.WriteRune(']')
	return builder.String()
}
