#!/usr/bin/env bash

# This script is only required because gofumpt doesn't return a different exit
# code if unformatted files were found.
# Consider inlining it into Makefile when this golang issue is fixed:
# https://github.com/golang/go/issues/46289

set -ex

if ! test -z "$(go run mvdan.cc/gofumpt -l .)"; then
    echo "gofumpt found at least one unformatted file"
    exit 1
fi
