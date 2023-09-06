#!/usr/bin/env bash

SELF_DIR=$(dirname "$0")
cd "${SELF_DIR}"/.. || exit


if [ "$(uname)" = "Darwin" ]; then
  if ! which gsed ; then
    # mac 下 如果没有 gsed，可以通过 brew install gnu-sed 来安装
    echo "no install gsed, please run: brew install gnu-sed"
    exit 1
  fi
fi

if [ "$(uname)" = "Darwin" ]; then
  # 因为 goimports 命令不能查看版本，所以 mac 电脑上直接安装一下项目中可以使用 goimports 的版本，避免有的机器安装的版本与项目指明的不统一
  go install golang.org/x/tools/cmd/goimports@v0.10.0
else
  # 默认 goimports 为 go 1.10 版本的太低了
  GOIMPORT_PATH=$(which goimports)
  if [ "${GOIMPORT_PATH}" == "/usr/bin/goimports" ] ; then
    go install golang.org/x/tools/cmd/goimports@v0.10.0
  fi
fi


GO_VERSION=$(go version)
expect_version="go version go1.19"
if ! [[ ${GO_VERSION} =~ ${expect_version} ]] ; then
  echo "please run 'brew install go@1.19' to update go version to 1.19.x"
  exit 1
fi

# 格式化修改的文件，提交到仓库和还未提交的 go 文件都会被格式化
# 修改的 go 文件较多的时候可能会比较慢
GIT_VERSION=$(git version)
GIT_REQ_VERSION="git version 2.22.0"
if [[ ${GIT_VERSION} < ${GIT_REQ_VERSION} ]];then
  echo "${GIT_VERSION} less than 2.22.0, please upgrade git version first"
  exit 1
fi
branch=$(git branch --show-current)
all=$(git diff "origin/master...${branch}" | grep "diff --git"; git status)
for file in $(echo "$all" | grep ".go$" | awk -F' ' '{print $NF}' | sed 's|^b/|./|'); do

  if [ ! -e "${file}" ]; then
    echo "ignore removed file : ${file}"
    continue
  fi

  echo "fmt: $file"
  bash script/fmt_imports.sh "$file"
  go fmt "$file"
done