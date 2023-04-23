//go:build tools
// +build tools

package go_containers

import _ "mvdan.cc/gofumpt"

// Running gofumpt through go:generate ensures it uses the
// version of gofumpt specified in go.mod.
//go:generate go run mvdan.cc/gofumpt -l -w .
