package container

func zeroValue[T any]() T {
	var result T
	return result
}
