#!/usr/bin/env bash

# This script is only required because gofumpt doesn't return
# a different exit code if unformatted files were found.
# Consider removing it and inlining "go generate gofumpt.go"
# into Makefile when this golang issue is fixed:
# https://github.com/golang/go/issues/46289

set -ex

if ! test -z "$(go generate gofumpt.go)"; then
    echo "gofumpt found at least one unformatted file"
    exit 1
fi
