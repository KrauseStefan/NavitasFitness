#!/usr/bin/env bash

ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
export GOPATH="$GOPATH:$ROOT"

cd $ROOT/src
echo "$(date +%H:%M:%S.%N) go fmt ./..."
go fmt ./...
echo "$(date +%H:%M:%S.%N) go vet ./..."
go vet ./...
echo "$(date +%H:%M:%S.%N) go get -d -v ./..."
go get -d -v ./...
echo "$(date +%H:%M:%S.%N) go test ./..."
go test ./...
echo "$(date +%H:%M:%S.%N) Finished"
cd -