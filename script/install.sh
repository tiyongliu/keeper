#!/bin/bash

node_home=0
node_version=0

go_home=0
go_version=0
function version_ge() { test "$(echo "$@" | tr " " "\n" | sort -rV | head -n 1)" == "$1"; }

echo
echo 'check node env'
if ! type node >/dev/null 2>&1; then
  node_home=0
  echo "请安装>=15.0以后的node.js，可跳转到https://nodejs.org/zh-cn/download安装"
   exit
else
  node_home=1
  node_version=$(bash -c "node --version" | sed 's/^v//')
fi

if [ $node_home -eq 1 ] && version_ge $node_version 15; then
    echo "check node env success"
  else
    echo "check node env failed"
    echo "请安装>=14.0以后的node.js，可跳转到https://nodejs.org/zh-cn/download安装"
    exit
fi

echo
echo "check go env"
if ! type go >/dev/null 2>&1; then
  go_home=0
  echo "请安装>=1.18以后的go，可跳转到https://studygolang.com/dl安装"
  exec
else
  go_home=1
  go_version=`go version | { read _ _ v _; echo ${v#go}; }`
fi

if [ $go_home -eq 1 ] && version_ge $go_version 1.18; then
   echo "check go env success"
  else
    echo "check go env failed"
    echo "请安装>=1.18以后的go，可跳转到https://studygolang.com/dl安装"
    exit
fi

source ./wails.sh

echo "check wails"
if ! type wails >/dev/null 2>&1; then
  read -p "是否使用go七牛云代理以实现高效下载?  [Y/n]" input
    case $input in
        [yY][eE][sS]|[yY])
    		echo "Yes"
    		go env -w GO111MODULE=on
    		go env -w GOPROXY=https://goproxy.cn,direct
    		;;
        [nN][oO]|[nN])
    		echo "No"
    		npm install -g pnpm
           	;;
        *)
    		echo "Invalid input..."
    		exit 1
    		;;
    esac
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  bash -c "go install github.com/wailsapp/wails/v2/cmd/wails@latest"
  exec
else
  go_home=1
  bash -c "wails update"
fi

echo
echo "================================= Successfully ================================="
echo "wails install Successfully"
echo "Get started with the following commands"

echo
echo "cd $(pwd)"
echo "wails dev"
