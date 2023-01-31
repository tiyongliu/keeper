#!/usr/bin/env bash

node_home=0
node_version=0

go_home=0
go_version=0
function version_ge() { test "$(echo "$@" | tr " " "\n" | sort -rV | head -n 1)" == "$1"; }

echo
echo 'check node env'
if ! type node >/dev/null 2>&1; then
  node_home=0
  echo "请安装>=14.0以后的node.js，可跳转到https://nodejs.org/zh-cn/download安装"
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
  go_version=$(bash -c "go version" | grep -oP '([1-9]\d*\.+?\d*)|(0\.\d*[1-9])') #非贪婪匹配
fi


if [ $go_home -eq 1 ] && version_ge $go_version 1.18; then
   echo "check go env success"
  else
    echo "check go env failed"
    echo "请安装>=1.18以后的go，可跳转到https://studygolang.com/dl安装"
fi

source wails.sh

echo
echo "check wails"
if ! type wails >/dev/null 2>&1; then
  bash -c "wails update"
  echo -e "\e[31m If your system is reporting that the wails command is missing, make sure you have followed the Go installation guide correctly. Normally, it means that the go/bin directory in your User's home directory is not in the PATH environment variable. You will also normally need to close and reopen any open command prompts so that changes to the environment made by the installer are reflected at the command prompt. \e[0m"
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
echo -e "\e[34m cd $(pwd) \e[0m"
echo -e "\e[34m wails dev \e[0m"