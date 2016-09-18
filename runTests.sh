#!/usr/bin/env bash

GOPATH="$HOME/TEST-GOPATH"

rm $GOPATH -rf

mkdir -p $GOPATH
mkdir -p "$GOPATH/src"

for dir in $(ls -d */); do
    if [[ $dir =~ ^[A-Z].* ]]; then
        ln -s "$HOME/git/NavitasFitness/${dir::-1}" "$HOME/TEST-GOPATH/src/"
    fi;
done

goapp get -v ./...

#ln -s "$GOPATH/src" "."


goapp test -v ./...
