
# keeper

## About
正在用项目使用vue3 + go 开发数据库可视化，目前只准备做桌面端，目前支持mysql、mongdb、后续会持续更新其他数据库的支持，会兼容window、macox、linux，由于我有正常的工作，都是在晚上跟周末开发。
采用wails提供的通信方式，绝非一般浏览器套壳gui，目前版本0.0.1，可以下载预览版。

[0.0.1预览版](https://github.com/tiyongliu/keeper/releases/tag/v0.0.1-demo)

[使用 Go + HTML + CSS + JS 构建漂亮的跨平台桌面应用](https://wails.io/zh-Hans/#)

## Progress

# frontend install 请按顺序执行[Please execute in sequence]!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
```sh
# 安装了pnpm请跳过这一步
pnpm install -g  

cd ./frontend
pnpm install
pnpm build
```

# backend install
[wails install](https://wails.io/zh-Hans/docs/gettingstarted/installation)