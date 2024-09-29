#!/usr/bin/env bash

set -eo pipefail

command -v go >/dev/null 2>&1 || { echo >&2 "go is not on the PATH. Aborting."; exit 1; }

cd x
go build .
./x "$1"
