#!/usr/bin/env bash
SF="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd ${SF}

export GOPATH="$GOPATH:$SF"

cd websrc
npm install
npm run clean
gnome-terminal -e "npm start" &
cd -

cd ipn-simulator
gnome-terminal -e "npm start" &
cd -

cd $SF/src/
cp ../NavitasFitnessConfig.Json "$SF/src/NavitasFitness/config.json"
go get -v ./... > /dev/null
#gnome-terminal -e "dev_appserver.py --dev_appserver_log_level=warning ." &
gnome-terminal -e "dev_appserver.py $SF/src/NavitasFitness/app.yaml" &
cd -
