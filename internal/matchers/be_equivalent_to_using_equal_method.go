package matchers

import (
	"fmt"
	"reflect"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

func BeEquivalentToUsingEqualMethod(value any) types.GomegaMatcher {
	return WithTransform(
		func(this any) (bool, error) {
			typ := reflect.TypeOf(this)
			noEqualMethodMessage := "to have an Equal method with a single parameter of type <" +
				typ.Name() +
				"> and a single return value of type <bool>"
			errNoEqualMethod := fmt.Errorf(format.Message(this, noEqualMethodMessage))

			equalMethod, ok := typ.MethodByName("Equal")
			if !ok {
				return false, errNoEqualMethod
			}

			if !hasReceiverAndParamOfSameType(equalMethod.Type) {
				return false, errNoEqualMethod
			}

			if equalMethod.Type.NumOut() != 1 {
				return false, errNoEqualMethod
			}

			if !equalMethod.Type.Out(0).AssignableTo(reflect.TypeOf(true)) {
				return false, errNoEqualMethod
			}

			params := []reflect.Value{
				reflect.ValueOf(this),
				reflect.ValueOf(value),
			}
			result := equalMethod.Func.Call(params)[0].Bool()
			return result, nil
		},
		BeTrue())
}

func hasReceiverAndParamOfSameType(methodType reflect.Type) bool {
	return methodType.NumIn() == 2 && methodType.In(0) == methodType.In(1)
}
