package main

import (
	"github.com/gin-gonic/gin"
	"web/middleware"
	"web/router"
	"web/db"
	"os"
)

func init() {
	db.ConnectDB()
}

func main() {

	// close db before unmount
	defer db.CloseDB()

	// 初始化
	g := gin.Default() // gin.New()

	// 注册中间件
	middleware.Register(g)

	// 添加路由
	router.Register(g)

	// 获取port参数
	port := ""
	if len(os.Args) >= 2 {
		port = os.Args[1]
	}
	if port == "" {
		port = ":8080"
	} else {
		port = ":" + port
	}

	// 启动服务
	g.Run(port)

}

