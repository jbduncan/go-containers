package graph

import "github.com/jbduncan/go-containers/set"

// Equal returns true if graph a and graph b meet the following rules, otherwise, it returns false:
//   - a.IsDirected() == b.IsDirected()
//   - a.AllowsSelfLoops() == b.AllowsSelfLoops()
//   - a.Nodes() and b.Nodes() are equal according to set.Equal
//   - a.Edges() and b.Edges() are equal according to set.Equal
//
// This method should be used over ==, the behaviour of which is undefined.
//
// Equal itself follows these rules:
//   - Reflexive: for any potentially-nil graph a, Equal(a, a) returns true.
//   - Symmetric: for any potentially-nil graphs a and b, Equal(a, b) and Equal(b, a) have the same results.
//   - Transitive: for any potentially-nil graphs a, b and c, if Equal(a, b) and Equal(b, c), then Equal(a, c) is true.
//   - Consistent: for any potentially-nil graphs a and b, multiple calls to Equal(a, b) consistently returns true or
//     consistently returns false, as long as the graphs do not change.
func Equal[T comparable](a, b Graph[T]) bool {
	if a == nil || b == nil {
		return a == b
	}

	return a.IsDirected() == b.IsDirected() &&
		a.AllowsSelfLoops() == b.AllowsSelfLoops() &&
		set.Equal(a.Nodes(), b.Nodes()) &&
		set.Equal(a.Edges(), b.Edges())
}
