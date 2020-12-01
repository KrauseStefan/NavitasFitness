# NavitasFitnes

# Environment setup

install google cloud SDK
run `gcloud components install cloud-datastore-emulator`

## Getting started

install google app engine, note the version in `src/NavitasFitness/app.yaml`

configure the go compiler with a go path (Tested with `go1.13.8`)

run `./launch.sh`

Configure dropbox integration
Use url: http://localhost:9000/rest/dropbox/authenticate

run all tests with `./runTests.sh`


## script descriptions
TODO

## relevant links

https://cloud.google.com/datastore/docs/tools/datastore-emulator
