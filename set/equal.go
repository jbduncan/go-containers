package set

import "iter"

// Equal returns true if set a has the same elements as set b in any order. Otherwise, it returns false.
//
// This method should be used over ==, the behaviour of which is undefined.
//
// Equal follows these rules:
//   - Reflexive: for any potentially-nil set a, Equal(a, a) returns true.
//   - Symmetric: for any potentially-nil sets a and b, Equal(a, b) and Equal(b, a) have the same results.
//   - Transitive: for any potentially-nil sets a, b and c, if Equal(a, b) and Equal(b, c), then Equal(a, c) is true.
//   - Consistent: for any potentially-nil sets a and b, multiple calls to Equal(a, b) consistently returns true or
//     consistently returns false, as long as the sets do not change.
//
// Note: Go needs the generic type to be defined explicitly, like:
//
//	a := set.Of(1)
//	b := set.Of(2)
//	result := set.Equal[int](a, b)
//	                   ^^^^^
func Equal[T comparable](a, b interface {
	Contains(element T) bool
	Len() int
	All() iter.Seq[T]
},
) bool {
	if a == nil || b == nil {
		return a == b
	}

	if a.Len() != b.Len() {
		return false
	}

	for element := range b.All() {
		if !a.Contains(element) {
			return false
		}
	}
	return true
}
