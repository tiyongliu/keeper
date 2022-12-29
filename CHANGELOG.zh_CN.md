## 0.0.2.2-alpha.2(2022-12-29)

### 🐛 Bug Fixes

## 0.0.2.2-alpha.1(2022-11-28)

### ✨ Features
- 新增数据库展示单行`SqlFormView.tsx`组件的开发

### 🐛 Bug Fixes
- 修复mysql连接池database被写死
- 改造后端连接数据库流程逻辑
- 重构数据库open方式
- 前端连接数据库代码精简
- driver storageSession 重构
- 修复`CellValue.vue` props value 如果传值为''，界面上会显示true。
- 修复`databaseConnections.go`SqlSelect方法 call sendRequest

- HandleSqlSelect方法runtime.EventsOn改成runtime.EventsOnce
- 程序加载完毕，判断是否自动跳转到上次记录页面

## 0.0.2.1-alpha.1(2022-09-23)

### ✨ Features

- 新增标签页`TabsPanel`组件的开发
- 新增列表页`TabRegister`组件的开发
- 新增mysql`TableDataTab`列表功能查询的开发
- 新增mysql`Columns`组件动态字段的显示、隐藏
- 新增前端`Columns`Filters 字段的模糊过滤
- 新增后端util方法

### 🐛 Bug Fixes
- 收藏、取消收藏关闭重启应用后失效
- 拖动窗口大小重新计算宽度
- 修复mysqlColumnInfo `AutoIncrement` 未正确赋值
- 修复数据库连接池引发程序崩溃



## 0.0.1-rc.1

### 🐛 Bug Fixes

- 修复重新加载database statusIcon消失
- 修复频繁刷新databases程序出现崩溃
- 修复刷新后databases无数据