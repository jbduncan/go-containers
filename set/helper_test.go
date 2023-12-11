package set_test

import "github.com/jbduncan/go-containers/set"

func oneElementSet() set.Set[string] {
	return setOf("link")
}

func twoElementSet() set.Set[string] {
	return setOf("link", "zelda")
}

func setOf(first string, others ...string) set.Set[string] {
	s := set.New[string]()
	s.Add(first)
	for _, elem := range others {
		s.Add(elem)
	}
	return s
}
