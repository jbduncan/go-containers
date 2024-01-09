package iteratortest

import (
	"fmt"
	"reflect"

	"github.com/jbduncan/go-containers/internal/slices"
	"github.com/jbduncan/go-containers/iterator"
)

type IteratorOrder int

const (
	knownOrder IteratorOrder = iota
	unknownOrder
)

type Tester[T any] struct {
	iteratorName     string
	expectedElements []T
	newIterator      func() iterator.Iterator[T]
	knownOrder       IteratorOrder
}

// TODO: Change to a style similar to
// https://www.calhoun.io/more-effective-ddd-with-interface-test-suites/
func ForIteratorWithKnownOrder[T any](
	iteratorName string,
	newIterator func() iterator.Iterator[T],
	expectedElements []T,
) *Tester[T] {
	return &Tester[T]{
		iteratorName:     iteratorName,
		expectedElements: slices.CopyToNonNilSlice(expectedElements),
		newIterator:      newIterator,
		knownOrder:       knownOrder,
	}
}

func ForIteratorWithUnknownOrder[T any](
	iteratorName string,
	newIterator func() iterator.Iterator[T],
	expectedElements []T,
) *Tester[T] {
	return &Tester[T]{
		iteratorName:     iteratorName,
		expectedElements: slices.CopyToNonNilSlice(expectedElements),
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
	opSequences := slices.CartesianProduct(slices.Repeat(uniqueOps, steps))

	for _, opSequence := range opSequences {
		actualIter := t.newIterator()
		remainingExpected := slices.CopyToNonNilSlice(t.expectedElements)
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

// At time of writing, we target Go 1.18 which doesn't have access to the
// builtin "max", because it is only available in Go 1.21+.
//
//goland:noinspection GoReservedWordUsedAsName
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
