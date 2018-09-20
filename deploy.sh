#!/usr/bin/env bash

root="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
export GOPATH="$GOPATH:$root"

# gcloud app deploy $root/src/NavitasFitness/app.yaml -v beta --no-promote --no-stop-previous-version
gcloud app deploy $root/src/NavitasFitness/app.yaml -v 1-3
