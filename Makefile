SHELL := /bin/bash

ifeq (, $(shell which go))
	$(error "No go in $(PATH)")
endif

.PHONY: format_apply
format_apply:
	go generate gofumpt_apply.go

.PHONY: format_check
format_check:
	scripts/format_check.sh

.PHONY: test
test:
	go generate test.go

.PHONY: check
check: format_check test

# TODO: Add another target for static analysis with golangci-lint:
#       https://golangci-lint.run/
#       Note: golangci-lint has gofumpt support, which should make
#       scripts/format_check.sh redundant.
#       Furthermore, it has support for:
#         - staticcheck (enabled by default)
#         - go vet (enabled by default)
#         - ginkgolinter
#         - bidichk (for checking dangerous Unicode characters in source code)
#         - gosec (for checking security problems)
# TODO: Add static analysis target to prerequisites of "check"
