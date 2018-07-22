package router

import (
  `github.com/gin-gonic/gin`
  "workerbook/conf"
  `workerbook/controller`
  `workerbook/middleware`
)

func registerConsoleUserRouter(g *gin.RouterGroup) {
  // jwt
  g.Use(middleware.Jwt)

  // 获取用户列表
  g.GET("", controller.C_GetUsersList)

  // 获取单个用户
  g.GET("/id/:id", controller.C_GetUserOne)

  // 添加用户
  g.POST("", controller.C_CreateUser)

  // 修改用户
  g.PUT("/id/:id", controller.C_UpdateUser)

  // 修改用户
  g.DELETE("/id/:id", controller.C_DelUserOne)
}

func registerUserRouter(g *gin.RouterGroup) {

  // 登录
  g.POST("/login", controller.UserLogin)

  // 获取个人信息
  g.GET("/profile",
    middleware.Jwt,
    controller.GetProfile)

  // 获取下级用户（包括自己）
  g.GET("/subordinate",
    middleware.Jwt,
    middleware.AllowRole(conf.RoleLeader, conf.RolePM, conf.RoleAdmin),
    controller.GetSubUsersList)
}
