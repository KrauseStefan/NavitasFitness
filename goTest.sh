#!/usr/bin/env bash

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

goapp fmt ./...
goapp get -d -v ./...
goapp test -v ./...
