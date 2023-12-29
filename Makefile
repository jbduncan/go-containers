SHELL := /bin/bash

ifeq (, $(shell which go))
	$(error "No go in $(PATH)")
endif

.PHONY: build
build:
	go build ./...

.PHONY: lint
lint:
	# Uses version imported by internal/tools.go, in turn using version in go.mod
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run

.PHONY: lint_fix
lint_fix:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run --fix
.PHONY: test
test:
	go test -shuffle=on -race ./...

.PHONY: check
# Build early to catch compiler errors sooner
check: build lint test

.PHONY: update_versions
update_versions:
	go get -u ./... && go get -u -t ./... && go mod tidy && go mod verify && go mod download

# TODO: Adopt https://github.com/uber-go/nilaway
# TODO: Add a 'go mod tidy' lint:
# https://github.com/uber-go/nilaway/blob/6b5d588e97aa719fc89271cda1c8aa7a804874bf/Makefile#L36-L41
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
#   - Examples in https://errorprone.info/docs/refaster
#   - Examples in https://github.com/PicnicSupermarket/error-prone-support
#   - Any other examples that tools like revive can catch but can't auto-fix

# TODO: Add eg templates for Graph.Equal and Set.Equal
# TODO: Add eg template for:
#   - `string == ""` or `string == ``` to `len(string) == 0`
