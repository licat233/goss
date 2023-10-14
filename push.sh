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
    comment="改为使用代码固定后version号，摒弃使用git的tag，以解决因为git带来的各种问题"
fi

git tag -a "v1.0.6" -m "$comment"
git add .
git commit -m "$comment"
git push -u origin main
