package router

import (
  `github.com/gin-gonic/gin`
  `workerbook/controller`
  `workerbook/middleware`
)

func registerProjectRouter(g *gin.RouterGroup) {

  // 获取项目列表
  g.GET("", middleware.Jwt, controller.GetProjectsList)

  // 获取单个用户
  g.GET("/:id", middleware.Jwt, controller.GetProjectOne)

  // 添加项目
  g.POST("", middleware.Jwt, controller.CreateProject)

  // 修改项目
  g.PUT("/:id", middleware.ConsoleJwt, controller.UpdateProject)

}
