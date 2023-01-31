![Go reference](https://img.shields.io/badge/go-v1.18-blue?logo=go&logoColor=white)
[![wails](https://img.shields.io/badge/wails-v2.3.1-brightgreen.svg)](https://wails.io)
[![vue3](https://img.shields.io/badge/vue-v3.2.0-7289da.svg?logo=v&logoColor=42b883)](https://vuejs.org/)
[![MIT license](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![dagate](https://img.shields.io/badge/dbgate-reference-brightgreen?&logoColor=white)](https://github.com/dbgate/dbgate)
[![styled with prettier](https://img.shields.io/badge/vben_admin-reference-ff69b4.svg)](https://vvbin.cn/doc-next/)

# keeper [构建漂亮的跨平台桌面应用](https://wails.io)

使用`go`、`wails`、`vue3`，目前只开发开发桌面端，目前仅支持`mysql`、`mongodb`、后续会持续更新其他数据库的支持，会兼容`window`、`macOs`、`linux`，由于我有正常的工作，都是利用空闲时间晚上跟周末开发。

## 预览版
![](https://cdn.jsdelivr.net/gh/422720735/easy_go@master/keeper.png)
![](https://cdn.jsdelivr.net/gh/422720735/easy_go@master/keeper-2.png)
![](https://cdn.jsdelivr.net/gh/422720735/easy_go@master/dark.png)

目前版本0.0.2.2，很多功能是不完善，并不能作为生成环境使用，如果你正准备用`wails`架构自己的桌面程序，本项目可作为一个参考，可以下载预览版尝鲜。[0.0.2.2](https://github.com/tiyongliu/keeper/releases/tag/0.0.2.2)

## 功能
- 使用`go`、`wails`、`vue3`构建桌面数据库可视化工具
- 使用wails提供事件系统go和Javascript之间的通信
- 目前仅支持`mongodb`、`mysql`的使用
- `windows` `macOs` `lunux` 发行版

# 开发环境安装

## 2. install
```shell
cd script
./install.sh
```

## 2. run app
```shell
wails dev #keeper同级目录下运行
```

## 如果安装过程中遇到问题，可以参考wails官方文档


## License
© tiyongliu, 2022-03.14~time.Now

Released under the [MIT License](./LICENSE)



# Supported by

[Jetbrains](https://www.jetbrains.com/)
![https://www.jetbrains.com/](release/img/jetbrains.svg)