//go:build tools
// +build tools

package main

import _ "github.com/onsi/ginkgo/v2/ginkgo"

// Running ginkgo through go:generate ensures it uses the
// version of ginkgo specified in go.mod.
//go:generate go run github.com/onsi/ginkgo/v2/ginkgo -r

// TODO: Improve above instruction to follow Ginkgo best
//       practices for CI:
//       https://onsi.github.io/ginkgo/#recommended-continuous-integration-configuration
