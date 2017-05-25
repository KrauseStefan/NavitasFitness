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

cp ../NavitasFitnessConfig.Json App/config.json
go get -v ./...
#gnome-terminal -e "dev_appserver.py --dev_appserver_log_level=warning ." &
gnome-terminal -e "dev_appserver.py App/app.yaml" &
