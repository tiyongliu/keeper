#!/bin/bash
cd $(pwd)/../frontend/
echo
echo "check pnpm package manager"

pnpm=0
if [ ! -d "$(pwd)/node_modules/" ];then #不存在
    pnpm=0
else
    pnpm=1
fi

echo
if [ $pnpm -eq 0 ] && ! type pnpm >/dev/null 2>&1; then
  read -p "需要安装pnpm软件包管理器, 是否使用淘宝镜像以实现高效下载?  [Y/n]" input
  case $input in
      [yY][eE][sS]|[yY])
  		echo "Yes"
  		npm install -g pnpm --registry=http://registry.npmmirror.com
  		pnpm install --registry https://registry.npm.taobao.org install any-touch
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
  elif [ $pnpm -eq 0 ] && type pnpm >/dev/null 2>&1; then
      pnpm install --registry https://registry.npm.taobao.org install any-touch
fi

if [ ! -d "$(pwd)/dist/" ];then
   npm run build
fi
