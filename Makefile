SHELL := /bin/bash

ifeq (, $(shell which go))
	$(error "No go in $(PATH)")
endif

.PHONY: fmt
fmt:
	go run mvdan.cc/gofumpt -l -w .

.PHONY: fmt_check
fmt_check:
	scripts/fmt_check.sh

.PHONY: test
test:
	go test ./...
# TODO: Run tests in parallel. Refer to:
#     - https://pkg.go.dev/cmd/go#hdr-Testing_flags
#     - https://stackoverflow.com/a/58987438/2252930
# TODO: Turn on the race detector
# TODO: Consider turning on test code coverage
# TODO: Randomly shuffle the order of tests everytime (try "go test -shuffle=on")

.PHONY: check
check: fmt_check test

# TODO: Add another target for static analysis with golangci-lint:
#       https://golangci-lint.run/
#       Note: golangci-lint has gofumpt support, which should make
#       scripts/fmt_check.sh redundant.
#       Furthermore, it has support for:
#         - staticcheck (enabled by default)
#         - go vet (enabled by default)
#         - ginkgolinter
#         - bidichk (for checking dangerous Unicode characters in source code)
#         - gosec (for checking security problems)
# TODO: Add static analysis target to prerequisites of "check"
