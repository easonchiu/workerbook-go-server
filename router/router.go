package router

import (
  `github.com/gin-gonic/gin`
  "workerbook/controller"
  "workerbook/middleware"
)

func Register(g *gin.Engine) {

  // 注册用户相关的路由
  registerUserRouter(g.Group("/users"))

  // 注册部门相关的路由
  registerDepartmentRouter(g.Group("/departments"))

  // 注册日报相关的路由
  // registerDailyRouter(g.Group("/dailies"))

  // 注册项目相关的路由
  registerProjectRouter(g.Group("/projects"))

  // 登录
  // g.POST("/login", controller.UserLogin)

  // 获取个人信息
  g.GET("/profile", middleware.Jwt, controller.GetProfile)

}
