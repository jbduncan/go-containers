package container_test

import (
	"github.com/onsi/gomega/types"
	"go-containers/container"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Graph", func() {
	Describe("given a new unordered endpoint pair", func() {
		Context("when calling .NodeU()", func() {
			It("returns the first endpoint", func() {
				endpointPair := container.NewUnorderedEndpointPair("link", "zelda")

				Expect(endpointPair).To(HaveNodeU("link"))
			})
		})

		Context("when calling .NodeV()", func() {
			It("returns the second endpoint", func() {
				endpointPair := container.NewUnorderedEndpointPair("link", "zelda")

				Expect(endpointPair).To(HaveNodeV("zelda"))
			})
		})
	})
})

func HaveNodeU(node string) types.GomegaMatcher {
	return WithTransform(
		func(endpointPair container.EndpointPair[string]) string {
			return endpointPair.NodeU()
		},
		Equal(node))
}

func HaveNodeV(node string) types.GomegaMatcher {
	return WithTransform(
		func(endpointPair container.EndpointPair[string]) string {
			return endpointPair.NodeV()
		},
		Equal(node))
}
