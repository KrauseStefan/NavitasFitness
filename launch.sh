#!/usr/bin/env bash
root="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $root
cp "$root/../NavitasFitnessConfig.Json" "$root/src/NavitasFitness/config.json"

export GOPATH="$GOPATH:$root"

cd $root/websrc
npm install
npm run clean
cd -

cd $root/src/
go get -v ./... > /dev/null
cd -

#gnome-terminal -e "dev_appserver.py $root/src/NavitasFitness/app.yaml" &
gnome-terminal --tab -e "npm start" --working-directory="$root/websrc" --title="client" --tab -e "npm start" --working-directory="$root/ipn-simulator" --title="IPN" --tab -e "dev_appserver.py --dev_appserver_log_level=warning $root/src/NavitasFitness/app.yaml" --title="server"
