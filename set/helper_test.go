package set_test

import "github.com/jbduncan/go-containers/set"

func oneElementSet() set.Set[string] {
	s := set.New[string]()
	s.Add("link")
	return s
}

func twoElementSet() set.Set[string] {
	s := set.New[string]()
	s.Add("link")
	s.Add("zelda")
	return s
}
