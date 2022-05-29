package app

import (
	"keeper/app/startup"
	"keeper/app/variable"
)

func RunApplication() {
	application := &variable.SystemApplication{
		Application: &variable.Application{
			Port: 8980,
		},
		RegisterHttpRoute: startup.RegisterHttpRoute,
	}
	run(application)
}

//gin 参考
//https://github.com/gin-gonic/examples

//一个支持多存储的文件列表程序，使用 Gin 和 React
//https://github.com/Xhofe/alist.git
