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
    comment="新增对文件格式的选择，可以通过文件格式进行文件筛选"
fi

git tag -a "v1.1.0-beta.4" -m "$comment"
git add .
git commit -m "$comment"
git push -u origin main
