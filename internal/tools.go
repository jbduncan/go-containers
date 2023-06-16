//go:build tools
// +build tools

package internal

// Import the tools we use so that the version numbers in go.mod survives
// a 'go mod tidy', in turn ensuring that the versions we use are
// deterministic.
import (
	_ "honnef.co/go/tools/cmd/staticcheck"
	_ "mvdan.cc/gofumpt"
)
