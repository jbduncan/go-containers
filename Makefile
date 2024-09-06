SHELL := /usr/bin/env bash

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
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.3 run
	go run go.uber.org/nilaway/cmd/nilaway@v0.0.0-20240821220108-c91e71c080b7 \
		-include-pkgs github.com/jbduncan/go-containers ./...
	./scripts/eg_lint.sh
	find . -name 'depaware.txt' | \
		xargs -n1 dirname | \
		xargs go run github.com/tailscale/depaware@v0.0.0-20210622194025-720c4b409502 --check

.PHONY: fix
fix:
	go mod tidy
	go mod download
	go run golang.org/x/tools/cmd/eg@v0.24.0 -t eg/fmt_errorf_to_errors_new.template -w ./...
	go run golang.org/x/tools/cmd/eg@v0.24.0 -t eg/rwmutex_lock_to_rlock.template -w ./...
	go run golang.org/x/tools/cmd/eg@v0.24.0 -t eg/time_now_sub_to_since.template -w ./...
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.3 run --fix
	go run github.com/tailscale/depaware@v0.0.0-20210622194025-720c4b409502 ./graph > ./graph/depaware.txt
	go run github.com/tailscale/depaware@v0.0.0-20210622194025-720c4b409502 ./set > ./set/depaware.txt
	go run github.com/tailscale/depaware@v0.0.0-20210622194025-720c4b409502 ./set/settest > ./set/settest/depaware.txt

.PHONY: test
test:
	go test -shuffle=on -race ./...

.PHONY: check
check: build lint test

.PHONY: update_versions
update_versions:
	go get -u -t ./... && go mod tidy && go mod verify && go mod download
	@echo "Make sure to update golangci-lint, eg, nilaway and depaware in the Makefile and scripts, too."
