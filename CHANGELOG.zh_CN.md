## 0.0.2.3-alpha.1(2023-1-18)

- **其它**
- 前端基础库`vben admin`升级到2023/1/18的版本

### ✨ Features
- 新增`mongodb`动态展示和查询`CollectionDataGridCore`、`CollectionJsonView`组件
- 新增鼠标右键菜单`ContextMenu.vue`组件，完全第一个版本的右键功能的封装
- 新增半自动化安装开发环境脚本



### 🐛 Bug Fixes
- v-splitterDrag自定义拖拽组件优化
- 修复theme`dark`、`light`切换html现实对应theme
- 修复标准盒子模型下高度超过预期，在原有高度-多余的高度，不改变原有盒子模型



## 0.0.2.2-alpha.2(2022-12-29)

- **其它**
- 替换`vben admin`自带解决白屏的loading动画，删除相loading图片资源
- 前端基础库`vben admin`升级到2023/1/7的版本

### 🐛 Bug Fixes
- 修复`WidgetColumnBarItem` props show属性
- 修复`PrimaryKeyFilterEditor`主键名、主键值显示
- 修复`VerticalScrollBar`组件初始化高度不正确，列表过滤scroll值不正确
- 修复wheel方法没被动触发到`VerticalScrollBar`的handleScroll方法
- 修复`DataGridCore.vue`wheel滚轮事件转发到子组件scroll，通过`dispatchEvent(new MouseEvent('scroll'))`
- 修复数据库查询错误处理方法(`setStatus`)，前端页面显示错误消息及图标
- 修复前端加载完毕后触发`doDatabasePing`、`doDatabasePing`、`subscribeRecentDatabaseSwitch`方法
- 修复前端加载完毕后如果当前显示`SqlFormView.tsx`点击标签页的关闭，需要点击2次
- 修复`TabsPanel.vue`wheel事件 scrollBy偏移量不生效，标签页没有去重



  
## 0.0.2.2-alpha.1(2022-11-28)

### ✨ Refactor

- 数据库连接池重构

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