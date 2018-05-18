package router

import "github.com/gin-gonic/gin"

func Register(g *gin.Engine) {

	// 注册用户相关的路由
	registerUserRouter(g.Group("/user"))

	// 注册分组相关的路由
	registerGroupRouter(g.Group("/group"))

	// 注册日报相关的路由
	registerDailyRouter(g.Group("/daily"))

	// 注册项目相关的路由
	registerProjectRouter(g.Group("/project"))

}
