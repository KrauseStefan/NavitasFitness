#!/usr/bin/env bash

root="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
export GOPATH="$GOPATH:$root"

version=v1-9-0

gcloud app deploy $root/src/NavitasFitness/app.yaml -v $version --no-promote --no-stop-previous-version
git tag -f $version
git push -f origin --tags

# gcloud app deploy src/NavitasFitness/cron.yaml