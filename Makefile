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
	go run github.com/onsi/ginkgo/v2/ginkgo -r
# TODO: Improve above instruction to follow Ginkgo best practices for CI:
#       https://onsi.github.io/ginkgo/#recommended-continuous-integration-configuration


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
