package iteratortest_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// TODO: Consider migrating from Ginkgo to Go's testing framework to reduce
//       dependencies.

func TestIteratorTester(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IteratorTest Suite")
}
