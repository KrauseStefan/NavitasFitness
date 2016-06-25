#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd ${DIR}


cd websrc

#npm install
npm run clean

gnome-terminal -e "npm run watch" &

cd -

cd ipn-simulator

gnome-terminal -e "npm start" &

cd -

gnome-terminal -e "goapp serve" &

