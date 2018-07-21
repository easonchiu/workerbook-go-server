package router

import (
  `github.com/gin-gonic/gin`
  `workerbook/controller`
  `workerbook/middleware`
)

func registerConsoleProjectRouter(g *gin.RouterGroup) {
  // jwt
  g.Use(middleware.Jwt)

  // 获取项目列表
  g.GET("", controller.C_GetProjectsList)

  // 获取单个项目
  g.GET("/id/:id", controller.C_GetProjectOne)

  // 删除单个项目
  g.DELETE("/id/:id", controller.C_DelProjectOne)

  // 添加项目
  g.POST("", controller.C_CreateProject)

  // 修改项目
  g.PUT("/id/:id", controller.C_UpdateProject)
}

func registerProjectRouter(g *gin.RouterGroup) {
  // jwt
  g.Use(middleware.Jwt)

  // 获取项目列表
  g.GET("", middleware.Jwt, controller.GetProjectsList)
}
