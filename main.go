package main

import (
	"github.com/gin-gonic/gin"
	"web/middleware"
	"web/router"
	"web/conf"
)

func init() {
	conf.ConnectDB()
}

func main() {

	// 初始化
	g := gin.Default() // gin.New()

	// 注册中间件
	middleware.Register(g)

	// 添加路由
	router.Register(g)

	// 启动服务
	g.Run(":8080")

}

