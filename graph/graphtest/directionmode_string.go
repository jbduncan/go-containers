// Code generated by "stringer -type=DirectionMode"; DO NOT EDIT.

package graphtest

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Directed-0]
	_ = x[Undirected-1]
}

const _DirectionMode_name = "DirectedUndirected"

var _DirectionMode_index = [...]uint8{0, 8, 18}

func (i DirectionMode) String() string {
	if i < 0 || i >= DirectionMode(len(_DirectionMode_index)-1) {
		return "DirectionMode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _DirectionMode_name[_DirectionMode_index[i]:_DirectionMode_index[i+1]]
}
