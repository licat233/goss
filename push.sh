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
    comment="修复oss保存路径错误的问题"
fi

git tag -a "v1.1.0-beta.8" -m "$comment"
git add .
git commit -m "$comment"
git push -u origin main
