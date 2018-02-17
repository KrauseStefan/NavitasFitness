#!/usr/bin/env bash

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
export GOPATH="$GOPATH:$ROOT"

code $ROOT/src