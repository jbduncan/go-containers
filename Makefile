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

# TODO: Add another target for static analysis with staticcheck:
#       https://github.com/dominikh/go-tools
# TODO: Add static analysis target to prerequisites of "check"
