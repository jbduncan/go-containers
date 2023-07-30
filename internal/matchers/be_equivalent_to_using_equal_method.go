package matchers

import (
	"fmt"
	"reflect"

	"github.com/onsi/gomega/gcustom"
	"github.com/onsi/gomega/types"
)

const noEqualMethodMessageFormat = "<%T> is expected to have an Equal " +
	"method with a single parameter of type <%T> and a single return value " +
	"of type <bool>"

func BeEquivalentToUsingEqualMethod(expected any) types.GomegaMatcher {
	return gcustom.MakeMatcher(
		func(actual any) (bool, error) {
			typ := reflect.TypeOf(actual)

			equalMethod, ok := typ.MethodByName("Equal")
			if !ok {
				return errNoEqualMethod(expected)
			}

			if !hasReceiverAndParamOfSameType(equalMethod.Type) {
				return errNoEqualMethod(expected)
			}

			if equalMethod.Type.NumOut() != 1 {
				return errNoEqualMethod(expected)
			}

			if !equalMethod.Type.Out(0).AssignableTo(reflect.TypeOf(true)) {
				return errNoEqualMethod(expected)
			}

			params := []reflect.Value{
				reflect.ValueOf(actual),
				reflect.ValueOf(expected),
			}
			result := equalMethod.Func.Call(params)[0].Bool()
			return result, nil
		}).
		WithTemplate(
			"Expected:\n{{.FormattedActual}}\n{{.To}} be equivalent to "+
				"according to .Equal():\n{{format .Data 1}}",
			expected)
}

func errNoEqualMethod(expected any) (bool, error) {
	result := fmt.Errorf(noEqualMethodMessageFormat, expected, expected)
	return false, result
}

func hasReceiverAndParamOfSameType(methodType reflect.Type) bool {
	return methodType.NumIn() == 2 && methodType.In(0) == methodType.In(1)
}
