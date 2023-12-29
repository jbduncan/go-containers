#!/usr/bin/env bash

set -ex

go get -u ./... && go get -u -t ./... && go mod tidy && go mod download
