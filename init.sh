#!/usr/bin/env bash

#Install golang
sudo apt install golang-1.12

$GOHOME="$HOME/gopath"
echo '' >> ~/.profile
echo '# Golang Settings' >> ~/.profile
echo 'PATH="/usr/lib/go-1.12/bin:$PATH"' >> ~/.profile
echo 'GOPATH="$GOHOME"' >> ~/.profile
mkdir $GOHOME

#download dependencies
go get -v ./...
go get gopkg.in/validator.v2


# Install google cloud tools
echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | sudo tee -a /etc/apt/sources.list.d/google-cloud-sdk.list
sudo apt install apt-transport-https ca-certificates
curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key --keyring /usr/share/keyrings/cloud.google.gpg add -
sudo apt update
sudo apt install google-cloud-sdk google-cloud-sdk-app-engine-python google-cloud-sdk-app-engine-go google-cloud-sdk-datastore-emulator
gcloud init


