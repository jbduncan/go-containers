SHELL := /bin/bash

ifeq (, $(shell which go))
	$(error "No go in $(PATH)")
endif

.PHONY: build
build:
	go build ./...

# TODO: Experiment with managing Go and golangci-lint with Nix.

.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2 run
	go run go.uber.org/nilaway/cmd/nilaway@v0.0.0-20231204220708-2f6a74d7c0e2 ./...
	@echo "Linting 'go mod tidy'..."
	@go mod tidy && \
		git diff --exit-code -- go.mod go.sum || \
		(echo "'go mod tidy' changed files" && false)
	@echo "Linting 'go mod verify'..."
	@go mod verify

.PHONY: lint_fix
lint_fix:
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

# TODO: Adopt eg (see https://rakyll.org/eg/), with refactoring templates for the following:
#   - Examples in https://github.com/golang/tools/tree/master/refactor/eg/testdata
#   - Examples in https://rakyll.org/eg/
#   - time package:
#     - == to Equal (even if revive catches this, being able to auto-fix it is valuable)
#     - time.Add(time.Duration(duration)) to time.Add(duration * time.Nanosecond)
#   - bool expression simplifications
#     - !(a >= b) to a < b
#     - !(a > b) to a <= b
#     - !(a <= b) to a > b
#     - !(a < b) to a >= b
#     - !(a != b) to a == b
#     - !(a == b) to a != b
#     - !!a to a
#   - `string == ""` or `string == ``` to `len(string) == 0`
#   - Examples in https://errorprone.info/docs/refaster
#   - Examples in https://github.com/PicnicSupermarket/error-prone-support
#   - Any other examples that tools like revive can catch but can't auto-fix

# TODO: Add eg templates for Graph.Equal and Set.Equal
