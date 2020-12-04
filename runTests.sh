#!/usr/bin/env bash
set -o errexit
set -o pipefail
set -o errtrace
set -o nounset
# set -o xtrace

echo "setting env variables:"
gcloud beta emulators datastore env-init
eval "$(gcloud beta emulators datastore env-init)"

./goTest.sh

cd websrc || exit 1
yarn test

cd - || exit 1
cd e2e || exit 1
yarn --frozen
yarn test

# Validate that we are not leaking file descriptersa
pid="$(ps -C NavitasFitness -o pid=)"
openFds="$(find "/proc/${pid##*( )}/fd" | wc -l)"
maxOpenFds=12

# This will likly first error on 2nd run, triggering and getting a baseline before the E2E test would be beneficial
if (( openFds>maxOpenFds )); then
red="\e[31m"
    white="\e[0m"
    echo -e "${red}Too many open file descripters ($openFds), validate that all resource are closed!$white"
    exit 1
else
    echo "Currently there are $openFds open file descripter for pid: $pid (hopefully) beloing to the NavitasFitness process"
fi
