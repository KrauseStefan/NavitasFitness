#!/usr/bin/env bash

root="$(git rev-parse --show-toplevel)"

export GOPATH="$GOPATH:$root"