#!/usr/bin/env bash

# 格式化一个目录下所有 go 文件中，import 时的空行


# 示例：
# bash script/fmt_imports.sh content-warmup
# bash script/fmt_imports.sh framework/rpc/member.go

usage() {
    echo "Usage: bash $0 path"
    exit 1
}

if [ $# -ne 1 ];then
    usage
fi

FILE_PATH=$1


if [ "$(uname)" = "Darwin" ]; then
  # mac 下 如果没有 gsed，可以通过 brew install gnu-sed 来安装
  gsed -i '/import (/, /)$/{/^$/d}' $(find $FILE_PATH -type f -name '*.go') && goimports -w "${FILE_PATH}"
else
  sed -i '/import (/, /)$/{/^$/d}' $(find $FILE_PATH -type f -name '*.go') && goimports -w "${FILE_PATH}"
fi