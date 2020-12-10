#!/usr/bin/env bash
set -o errexit
set -o pipefail
set -o errtrace
set -o nounset
# set -o xtrace

. ./configure-gopath.sh
cd "$root" || exit 1
cp "$root/../NavitasFitnessConfig.Json" "$root/src/config.json"

cd "$root/websrc" || exit 1
yarn --frozen
yarn clean
cd - || exit 1

cd "$root/src/" || exit 1
go get -v ./... > /dev/null
cd - || exit 1

echo "setting env variables:"
gcloud beta emulators datastore env-init
eval "$(gcloud beta emulators datastore env-init)"

# https://cloud.google.com/sdk/gcloud/reference/beta/emulators/datastore/start
# gcloud beta emulators datastore start --host-port=0.0.0.0:8081
# In January the datastore will be upgraded and will no longer have consistency issues...
# gcloud beta emulators datastore start --host-port=0.0.0.0:8081 --consistency=1.0
# gcloud beta emulators datastore start --host-port=0.0.0.0:8081 --consistency=0.1

# Ubuntu code
# gnome-terminal \
#   --tab -e "npm start" --working-directory="$root/websrc" --title="client" \
#   --tab -e "npm start" --working-directory="$root/ipn-simulator" --title="IPN" \
#   --tab -e "gcloud beta emulators datastore start --consistency=1.0" --title="server" \
#   --tab -e "go run NavitasFitness" --title="server"

# WSL2 with conemu code
bash.exe -new_console:t:"Client" -c "source \$HOME/.profile; cd $root/websrc; npm start" || true
bash.exe -new_console:t:"IPN" -c "source \$HOME/.profile; cd $root/ipn-simulator; npm start" || true
bash.exe -new_console:t:"Datastore simulator" -c "gcloud beta emulators datastore start --consistency=1.0" || true
bash.exe -new_console:t:"server" -c "source \$HOME/.profile; source ./configure-gopath.sh; eval \"$(gcloud beta emulators datastore env-init)\"; cd src/ ; go run main.go" || true
