//go:build tools
// +build tools

package tools

// Import the tools we use so that the version numbers in go.mod survives
// a 'go mod tidy', in turn ensuring that the versions we use are
// deterministic.
import (
	_ "github.com/onsi/ginkgo/v2/ginkgo"
	_ "mvdan.cc/gofumpt"
)
