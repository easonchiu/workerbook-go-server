package router

import (
  `github.com/gin-gonic/gin`
  `workerbook/controller`
)

func Register(g *gin.Engine) {

  // 注册用户相关的路由
  registerUserRouter(g.Group("/users"))

  // 注册分组相关的路由
  registerGroupRouter(g.Group("/groups"))

  // 注册日报相关的路由
  registerDailyRouter(g.Group("/dailies"))

  // 注册项目相关的路由
  registerProjectRouter(g.Group("/projects"))

  // 其他路由
  g.POST("/login", controller.UserLogin)

}
