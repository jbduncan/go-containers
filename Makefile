SHELL := /bin/bash

ifeq (, $(shell which go))
	$(error "No go in $(PATH)")
endif

.PHONY: fmt
fmt: # Uses version in internal/tools.go
	go run mvdan.cc/gofumpt -l -w .

.PHONY: fmt_check
fmt_check:
	scripts/fmt_check.sh

.PHONY: vet
vet:
	go vet ./...

.PHONY: staticcheck
staticcheck: # Uses version in internal/tools.go
	go run honnef.co/go/tools/cmd/staticcheck ./...

.PHONY: test
test:
	go test -shuffle=on -race ./...

.PHONY: check
check: fmt_check vet staticcheck test
