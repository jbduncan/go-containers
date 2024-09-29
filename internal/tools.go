//go:build tools
// +build tools

package internal

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/tailscale/depaware"
	_ "go.uber.org/nilaway/cmd/nilaway"
	_ "golang.org/x/tools/cmd/eg"
)
