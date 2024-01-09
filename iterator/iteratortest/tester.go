package iteratortest

import (
	"fmt"
	"reflect"

	slices2 "github.com/jbduncan/go-containers/internal/slices"
	"github.com/jbduncan/go-containers/iterator"
	"golang.org/x/exp/slices"
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
		expectedElements: slices2.CopyToNonNilSlice(expectedElements),
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
		expectedElements: slices2.CopyToNonNilSlice(expectedElements),
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
	opSequences := slices2.CartesianProduct(slices2.Repeat(uniqueOps, steps))

	for _, opSequence := range opSequences {
		actualIter := t.newIterator()
		remainingExpected := slices2.CopyToNonNilSlice(t.expectedElements)
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
		return slices.Delete(remainingExpected, 0, 1), nil

	case unknownOrder:
		i := slices.IndexFunc(
			remainingExpected,
			func(t T) bool { return equal(t, value) },
		)
		if i == -1 {
			return remainingExpected, t.misbehavingIteratorError(opSequence)
		}
		return slices.Delete(remainingExpected, i, i+1), nil
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
