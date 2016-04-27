#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd ${DIR}

./node_modules/.bin/gulp clean


gnome-terminal -e "./node_modules/.bin/gulp buildAndWatch" &

gnome-terminal -e "goapp serve ./app-engine/" &

cd ipn-simulator

gnome-terminal -e "node index.js" &


cd -
