#!/usr/bin/env bash

set -ex

if ! test -z "$(go run golang.org/x/tools/cmd/eg@v0.24.0 -t eg/fmt_errorf_to_errors_new.template ./...)"; then
		echo "eg found at least one non-conforming file."
		echo "Run 'make lint_fix' to fix them."
		exit 1
fi

if ! test -z "$(go run golang.org/x/tools/cmd/eg@v0.24.0 -t eg/rwmutex_lock_to_rlock.template ./...)"; then
		echo "eg found at least one non-conforming file."
		echo "Run 'make lint_fix' to fix them."
		exit 1
fi

if ! test -z "$(go run golang.org/x/tools/cmd/eg@v0.24.0 -t eg/time_now_sub_to_since.template ./...)"; then
		echo "eg found at least one non-conforming file."
		echo "Run 'make lint_fix' to fix them."
		exit 1
fi
