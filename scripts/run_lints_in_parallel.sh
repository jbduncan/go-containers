#!/usr/bin/env bash

set -ex

(trap 'kill 0' SIGINT; \
		go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2 run & \
 		go run go.uber.org/nilaway/cmd/nilaway@v0.0.0-20231204220708-2f6a74d7c0e2 \
			-include-pkgs github.com/jbduncan/go-containers ./...  & \
		./scripts/eg_lint.sh & \
		wait)
