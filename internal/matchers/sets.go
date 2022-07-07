package matchers

import (
	"fmt"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
	"go-containers/container/set"
	"reflect"

	. "github.com/onsi/gomega"
)

func HaveLenOf(len int) types.GomegaMatcher {
	return And(
		WithTransform(
			func(value any) (int, error) {
				errNotLenMethod := fmt.Errorf(format.Message(value, "to have a Len method with a single return value of type <int>"))

				typ := reflect.TypeOf(value)
				lenMethod, ok := typ.MethodByName("Len")
				if !ok {
					return 0, errNotLenMethod
				}

				if hasReceiverAndNoParams(lenMethod) {
					return 0, errNotLenMethod
				}

				if lenMethod.Type.NumOut() != 1 {
					return 0, errNotLenMethod
				}

				if !lenMethod.Type.Out(0).AssignableTo(reflect.TypeOf(0)) {
					return 0, errNotLenMethod
				}

				result := lenMethod.Func.Call([]reflect.Value{reflect.ValueOf(value)})[0].Int()
				return int(result), nil
			},
			Equal(len)))
}

func hasReceiverAndNoParams(method reflect.Method) bool {
	return method.Type.NumIn() != 1
}

func Contain[T comparable](elem T) types.GomegaMatcher {
	return WithTransform(
		func(set set.Set[T]) bool {
			return set.Contains(elem)
		},
		BeTrue())
}

func ForEachToSlice[T comparable](s set.Set[T]) []T {
	var result []T

	s.ForEach(func(elem T) {
		result = append(result, elem)
	})

	return result
}
