#!/bin/bash

#进入monitor mode
set -m

current_path=$(
    cd $(dirname $0)
    pwd
)

cd "$current_path" || exit

cd ../

go build .

cd example

../goss html --dir=.
