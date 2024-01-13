package iteratortest_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestIteratorTester(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IteratorTest Suite")
}
