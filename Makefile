SHELL := /bin/bash

ifeq (, $(shell which go))
	$(error "No go in $(PATH)")
endif

.PHONY: build
build:
	go build ./...

.PHONY: lint
lint:
	@echo "Linting 'go mod tidy'..."
	@go mod tidy && \
		git diff --exit-code -- go.mod go.sum || \
		(echo "'go mod tidy' changed files" && false)
	@echo "Linting 'go mod verify'..."
	@go mod verify
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.1 run
	go run go.uber.org/nilaway/cmd/nilaway@v0.0.0-20240208191803-c4313bf3c543 \
		-include-pkgs github.com/jbduncan/go-containers ./...
	./scripts/eg_lint.sh

.PHONY: lint_fix
lint_fix:
	go mod tidy
	go mod download
	go run golang.org/x/tools/cmd/eg@v0.17.0 -t eg/fmt_errorf_to_errors_new.template -w ./...
	go run golang.org/x/tools/cmd/eg@v0.17.0 -t eg/rwmutex_lock_to_rlock.template -w ./...
	go run golang.org/x/tools/cmd/eg@v0.17.0 -t eg/time_now_sub_to_since.template -w ./...
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.1 run --fix

.PHONY: test
test:
	go test -shuffle=on -race ./...

.PHONY: check
check: build lint test

.PHONY: update_versions
update_versions:
	go get -u ./... && go get -u -t ./... && go mod tidy && go mod verify && go mod download
	@echo "Make sure to update golangci-lint, eg and nilaway in the Makefile, too."
