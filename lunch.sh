#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd ${DIR}

gulp clean

gnome-terminal -e "gulp buildAndWatch" &

gnome-terminal -e "goapp serve ./app-engine/" &

cd -