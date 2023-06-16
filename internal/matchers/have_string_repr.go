package matchers

import (
	"fmt"

	"github.com/onsi/gomega/types"

	//lint:ignore ST1001 dot importing gomega matchers is best practice and
	// this package is used by test code only
	. "github.com/onsi/gomega"
)

func HaveStringRepr(valueOrMatcher any) types.GomegaMatcher {
	value, ok := valueOrMatcher.(string)
	if ok {
		// TODO: Use gcustom.MakeMatcher to improve error message
		return WithTransform(
			func(stringer fmt.Stringer) string {
				return stringer.String()
			},
			Equal(value))
	}

	matcher, ok := valueOrMatcher.(types.GomegaMatcher)
	if ok {
		// TODO: Use gcustom.MakeMatcher to improve error message
		return WithTransform(
			func(stringer fmt.Stringer) string {
				return stringer.String()
			},
			matcher)
	}

	panic(fmt.Sprintf(
		"valueOrMatcher is neither a string nor a types.GomegaMatcher: %v",
		valueOrMatcher))
}
