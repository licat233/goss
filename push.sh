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
    comment="完善一些列功能，结构调整"
fi

git tag -a "v1.1.0-beta.2" -m "$comment"
git add .
git commit -m "$comment"
git push -u origin main
