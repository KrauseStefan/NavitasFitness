#!/usr/bin/env bash

./goTest.sh

cd websrc
npm test

cd -
cd e2e
npm install
npm test
