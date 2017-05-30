#!/usr/bin/env bash
root="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $root
cp "$root/../NavitasFitnessConfig.Json" "$root/src/NavitasFitness/config.json"

export GOPATH="$GOPATH:$root"

cd websrc
npm install
npm run clean
gnome-terminal -e "npm start" &
cd -

cd ipn-simulator
gnome-terminal -e "npm start" &
cd -

cd $root/src/
go get -v ./... > /dev/null
cd -

#gnome-terminal -e "dev_appserver.py $root/src/NavitasFitness/app.yaml" &
gnome-terminal -e "dev_appserver.py --dev_appserver_log_level=warning $root/src/NavitasFitness/app.yaml" &
