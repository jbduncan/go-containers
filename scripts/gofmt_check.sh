#!/usr/bin/env bash

# This script is only required because gofmt doesn't return a different exit
# code if unformatted files were found.

set -ex

if ! test -z "$(go run cmd/gofmt -r 'interface{} -> any' -l .)"; then
    echo "gofmt found at least one file that uses interface{} instead of any."
    echo "Run 'make lint_fix' to fix them."
    exit 1
fi

if ! test -z "$(go run cmd/gofmt -r '(a) -> a' -l .)"; then
    echo "gofmt found at least one file that has unnecessary parentheses."
    echo "Run 'make lint_fix' to fix them."
    exit 1
fi

if ! test -z "$(go run cmd/gofmt -r 'a[b:len(a)] -> a[b:]' -l .)"; then
    echo "gofmt found at least one file that uses a[b:len(a)] instead of a[b:]."
    echo "Run 'make lint_fix' to fix them."
    exit 1
fi

if ! test -z "$(go run cmd/gofmt -s -l .)"; then
    echo "gofmt found at least one file that can be simplified."
    echo "Run 'make lint_fix' to fix them."
    exit 1
fi
