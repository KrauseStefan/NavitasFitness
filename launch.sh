#!/usr/bin/env bash
SF="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
GO_APP_DIR="$SF/App"
cd ${GO_APP_DIR}

GO_TEST_PATH="$HOME/TEST-GOPATH"

rm $GO_TEST_PATH -rf

mkdir -p $GO_TEST_PATH
mkdir -p "$GO_TEST_PATH/src"

for dir in $(ls -d */); do
    if [[ $dir =~ ^[A-Z].* ]]; then
        ln -sf "$GO_APP_DIR/${dir::-1}" "$GO_TEST_PATH/src/"
    fi;
done

cd ${SF}

export GOPATH="$GOPATH:$GO_TEST_PATH"

cd websrc

npm install
npm run clean

gnome-terminal -e "npm start" &

cd -

cd ipn-simulator

gnome-terminal -e "npm start" &

cd -

cp ../NavitasFitnessConfig.Json "$GO_APP_DIR/NavitasFitness/config.json"
go get -v ./... > /dev/null
#gnome-terminal -e "dev_appserver.py --dev_appserver_log_level=warning ." &
gnome-terminal -e "dev_appserver.py $GO_APP_DIR/NavitasFitness/app.yaml" &
