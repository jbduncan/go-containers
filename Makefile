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
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2 run
	go run go.uber.org/nilaway/cmd/nilaway@v0.0.0-20231204220708-2f6a74d7c0e2 \
		-include-pkgs github.com/jbduncan/go-containers ./...
	./scripts/eg_lint.sh

.PHONY: lint_fix
lint_fix:
	go run golang.org/x/tools/cmd/eg@v0.16.1 -t eg/fmt_errorf_to_errors_new.template -w ./...
	go run golang.org/x/tools/cmd/eg@v0.16.1 -t eg/rwmutex_lock_to_rlock.template -w ./...
	go run golang.org/x/tools/cmd/eg@v0.16.1 -t eg/time_now_sub_to_since.template -w ./...
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2 run --fix

.PHONY: test
test:
	go test -shuffle=on -race ./...

.PHONY: check
# Build early to catch compiler errors sooner
check: build lint test

.PHONY: update_versions
update_versions:
	go get -u ./... && go get -u -t ./... && go mod tidy && go mod verify && go mod download
	@echo "Make sure to update the golangci-lint version in Makefile, too."
