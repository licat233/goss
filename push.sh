#!/bin/bash

#进入monitor mode
set -m

current_path=$(
    cd $(dirname $0)
    pwd
)

cd "$current_path" || exit

# ./examples/update.sh

comment=$1

if [ -z "$comment" ]; then
    comment="原来的bar库有Bug，已替换成另一个库的bar"
fi

git tag -a "v1.0.8" -m "$comment"
git add .
git commit -m "$comment"
git push -u origin main
