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
    comment="优化新增对文件格式的选择"
fi

git tag -a "v1.1.0-beta.5" -m "$comment"
git add .
git commit -m "$comment"
git push -u origin main
