#!/usr/bin/env bash

# Hints
# sudo ln -s ~/Apps/google-cloud-sdk/platform/google_appengine/goroot/src/appengine $GOROOT/src/
# sudo ln -s ~/Apps/google-cloud-sdk/platform/google_appengine/goroot/src/appengine_internal $GOROOT/src/

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd ${DIR}

export GOPATH="$HOME/TEST-GOPATH"

rm $GOPATH -rf

mkdir -p $GOPATH
mkdir -p "$GOPATH/src"

for dir in $(ls -d */); do
    if [[ $dir =~ ^[A-Z].* ]]; then
        ln -s "$HOME/git/NavitasFitness/${dir::-1}" "$HOME/TEST-GOPATH/src/"
    fi;
done

go fmt ./...
go get -d -v ./...
go test -v ./...
