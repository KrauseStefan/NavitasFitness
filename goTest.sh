#!/usr/bin/env bash

# Hints
# sudo ln -s ~/Apps/google-cloud-sdk/platform/google_appengine/goroot/src/appengine $GOROOT/src/
# sudo ln -s ~/Apps/google-cloud-sdk/platform/google_appengine/goroot/src/appengine_internal $GOROOT/src/

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd ${DIR}

GO_TEST_PATH="$HOME/TEST-GOPATH"

#rm $GO_TEST_PATH -rf

mkdir -p $GO_TEST_PATH
mkdir -p "$GO_TEST_PATH/src"

for dir in $(ls -d */); do
    if [[ $dir =~ ^[A-Z].* ]]; then
        ln -sf "$HOME/git/NavitasFitness/${dir::-1}" "$HOME/TEST-GOPATH/src/"
    fi;
done

export GOPATH="$GOPATH:$GO_TEST_PATH"

echo "$(date +%H:%M:%S.%N) go fmt"
go fmt ./...
echo "$(date +%H:%M:%S.%N) go vet"
go vet ./...
echo "$(date +%H:%M:%S.%N) go get"
go get -d -v ./...
echo "$(date +%H:%M:%S.%N) go test"
go test -v ./...
echo "$(date +%H:%M:%S.%N) Finished"
