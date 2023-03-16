package iteratortester

import (
	"fmt"
	"go-containers/container/iterator"
	"reflect"
)

type IteratorOrder int

const (
	knownOrder = iota
	unknownOrder
)

type Tester[T any] struct {
	iteratorName     string
	expectedElements []T
	newIterator      func() iterator.Iterator[T]
	knownOrder       IteratorOrder
}

func ForIteratorWithKnownOrder[T any](
	iteratorName string,
	newIterator func() iterator.Iterator[T],
	expectedElements []T) *Tester[T] {

	return &Tester[T]{
		iteratorName:     iteratorName,
		expectedElements: copyToNonNilSlice(expectedElements),
		newIterator:      newIterator,
		knownOrder:       knownOrder,
	}
}

func ForIteratorWithUnknownOrder[T any](
	iteratorName string,
	newIterator func() iterator.Iterator[T],
	expectedElements []T) *Tester[T] {

	return &Tester[T]{
		iteratorName:     iteratorName,
		expectedElements: copyToNonNilSlice(expectedElements),
		newIterator:      newIterator,
		knownOrder:       unknownOrder,
	}
}

type valueOp struct{}

func (v valueOp) String() string {
	return "Value()"
}

type nextOp struct{}

func (n nextOp) String() string {
	return "Next()"
}

func (t Tester[T]) Test() error {
	steps := max(5, len(t.expectedElements)+1)
	uniqueOps := []any{valueOp{}, nextOp{}}
	opSequences := cartesianProduct(repeat(uniqueOps, steps))

	for _, opSequence := range opSequences {
		actualIter := t.newIterator()
		remainingExpected := copyToNonNilSlice(t.expectedElements)
		for _, op := range opSequence {
			switch op.(type) {
			case valueOp:
				var err error
				remainingExpected, err = t.doValueOpAndCheck(actualIter, remainingExpected, opSequence)
				if err != nil {
					return err
				}

			case nextOp:
				err := t.doNextOpAndCheck(actualIter, remainingExpected, opSequence)
				if err != nil {
					return err
				}

			default:
				panic(fmt.Sprintf("unrecognised op: %v", op))
			}
		}
	}

	return nil
}

func (t Tester[T]) doValueOpAndCheck(actualIter iterator.Iterator[T], remainingExpected []T, opSequence []any) ([]T, error) {
	if len(remainingExpected) == 0 {
		if !panics(func() { actualIter.Value() }) {
			return remainingExpected, t.misbehavingIteratorError(opSequence)
		}

		return remainingExpected, nil
	}

	value := actualIter.Value()

	switch t.knownOrder {
	case knownOrder:
		if !equal(value, remainingExpected[0]) {
			return remainingExpected, t.misbehavingIteratorError(opSequence)
		}
		return deleteInPlace(remainingExpected, 0), nil

	case unknownOrder:
		i := indexOf(remainingExpected, value)
		if i == -1 {
			return remainingExpected, t.misbehavingIteratorError(opSequence)
		}
		return deleteInPlace(remainingExpected, i), nil
	}

	panic(fmt.Sprintf("unrecognised order: %v", t.knownOrder))
}

func panics(f func()) (panicked bool) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
		}
	}()

	f()

	panicked = false

	return
}

func equal(x any, y any) bool {
	return reflect.DeepEqual(x, y)
}

func (t Tester[T]) doNextOpAndCheck(actualIter iterator.Iterator[T], remainingExpected []T, opSequence []any) error {
	if actualIter.Next() {
		if len(remainingExpected) == 0 {
			return t.misbehavingIteratorError(opSequence)
		}
		return nil
	}

	if len(remainingExpected) != 0 {
		return t.misbehavingIteratorError(opSequence)
	}
	return nil
}

func (t Tester[T]) misbehavingIteratorError(opSequence []any) error {
	return fmt.Errorf("iterator '%s' misbehaves when running operations %v", t.iteratorName, opSequence)
}

func deleteInPlace[T any](values []T, index int) []T {
	return append(values[:index], values[index+1:]...)
}

func indexOf[T any](haystack []T, needle T) int {
	for i, x := range haystack {
		if equal(x, needle) {
			return i
		}
	}
	return -1
}

func repeat[T any](value T, times int) []T {
	result := make([]T, times)
	for i := range result {
		result[i] = value
	}
	return result
}

func cartesianProduct[T any](values [][]T) [][]T {
	result := [][]T{{}}
	for _, innerValues := range values {
		var newResult [][]T
		for _, rest := range result {
			for _, tail := range innerValues {
				newResult = append(newResult, copyToNonNilSlice(append(rest, tail)))
			}
		}
		result = newResult
	}
	return result
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func copyToNonNilSlice[T any](values []T) []T {
	result := make([]T, len(values))
	copy(result, values)
	return result
}
