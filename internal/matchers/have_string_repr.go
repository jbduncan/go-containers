package matchers

import (
	"fmt"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gcustom"
	"github.com/onsi/gomega/types"
)

func HaveStringRepr(value string) types.GomegaMatcher {
	return gcustom.MakeMatcher(
		func(stringer fmt.Stringer) (bool, error) {
			actual := stringer.String()
			return value == actual, nil
		}).
		WithTemplate("Expected String() of\n{{.FormattedActual}}\n{{.To}} equal {{.Data}}").
		WithTemplateData(value)
}

func HaveStringReprThatIsAnyOf(elems ...any) types.GomegaMatcher {
	return gcustom.MakeMatcher(
		func(stringer fmt.Stringer) (bool, error) {
			actual := stringer.String()
			return gomega.BeElementOf(elems).Match(actual)
		}).
		WithTemplate("Expected String() of\n{{.FormattedActual}}\n{{.To}} be any of\n{{format .Data 1}}").
		WithTemplateData(elems)
}
