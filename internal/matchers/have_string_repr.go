package matchers

import (
	"fmt"

	"github.com/onsi/gomega/gcustom"
	"github.com/onsi/gomega/types"
)

func HaveStringRepr(valueOrMatcher any) types.GomegaMatcher {
	if value, ok := valueOrMatcher.(string); ok {
		return gcustom.MakeMatcher(
			func(stringer fmt.Stringer) (bool, error) {
				actual := stringer.String()
				return value == actual, nil
			}).
			WithTemplate("Expected String() of\n{{.FormattedActual}}\n{{.To}} equal {{.Data}}").
			WithTemplateData(value)
	}

	if matcher, ok := valueOrMatcher.(types.GomegaMatcher); ok {
		return gcustom.MakeMatcher(
			func(stringer fmt.Stringer) (bool, error) {
				actual := stringer.String()
				return matcher.Match(actual)
			}).
			WithTemplate("Expected String() of\n{{.FormattedActual}}\n{{.To}} satisfy matcher\n{{format .Data 1}}").
			WithTemplateData(matcher)
	}

	panic(fmt.Sprintf(
		"valueOrMatcher is neither a string nor a types.GomegaMatcher: %v",
		valueOrMatcher))
}
