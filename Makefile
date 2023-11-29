SHELL := /bin/bash

ifeq (, $(shell which go))
	$(error "No go in $(PATH)")
endif

.PHONY: fmt
fmt: # Uses version imported by internal/tools.go, in turn using version in go.mod
	go run mvdan.cc/gofumpt -l -w .

.PHONY: fmt_check
fmt_check:
	scripts/fmt_check.sh

.PHONY: vet
vet:
	go vet ./...

# TODO: Consider replacing with golangci-lint, even if it has to be installed manually:
# https://github.com/uber-go/nilaway/blob/6b5d588e97aa719fc89271cda1c8aa7a804874bf/Makefile#L26-L34
# Alternatively, do it as tailscale have:
# https://github.com/tailscale/tailscale/commit/280255acae604796a1113861f5a84e6fa2dc6121
.PHONY: staticcheck
staticcheck: # Uses version imported by internal/tools.go, in turn using version in go.mod
	go run honnef.co/go/tools/cmd/staticcheck ./...

.PHONY: test
test:
	go test -shuffle=on -race ./...

.PHONY: check
check: fmt_check vet staticcheck test

# TODO: Add https://github.com/uber-go/nilaway
# TODO: Add a 'go mod tidy' lint:
# https://github.com/uber-go/nilaway/blob/6b5d588e97aa719fc89271cda1c8aa7a804874bf/Makefile#L36-L41
# TODO: Adopt eg (see https://rakyll.org/eg/), with refactoring templates for:
#   - Examples in https://github.com/golang/tools/tree/master/refactor/eg/testdata
#   - Examples in https://rakyll.org/eg/
#   - time package:
#     - == to Equal
#     - time.Add(time.Duration(duration)) to time.Add(duration * time.Nanosecond)
#   - bool expression simplifications
#     - !(a >= b) to a < b
#     - !(a > b) to a <= b
#     - !(a <= b) to a > b
#     - !(a < b) to a >= b
#     - !(a != b) to a == b
#     - !(a == b) to a != b
#     - !!a to a
#   - Examples in https://errorprone.info/docs/refaster
#   - Examples in https://github.com/PicnicSupermarket/error-prone-support

# TODO: Add eg templates for Graph.Equal and Set.Equal
# TODO: Add eg template for:
#   - `string == ""` or `string == ``` to `len(string) == 0`

# TODO: Extract common eg templates into its own Git repo