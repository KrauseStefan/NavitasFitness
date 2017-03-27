#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd ${DIR}

cd websrc

npm install
npm run clean

gnome-terminal -e "npm start" &

cd -

cd ipn-simulator

gnome-terminal -e "npm start" &

cd -

goapp get -v ./...
gnome-terminal -e "/home/stefan/Apps/go_appengine/dev_appserver.py ." &
