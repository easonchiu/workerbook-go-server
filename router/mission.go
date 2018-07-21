package router

import (
  "github.com/gin-gonic/gin"
  "workerbook/controller"
  "workerbook/middleware"
)

func registerConsoleMissionRouter(g *gin.RouterGroup) {
  // jwt
  g.Use(middleware.Jwt)

  // 任务列表
  g.GET("", controller.C_GetMissionsList)

  // 获取单个任务
  g.GET("/id/:id", controller.C_GetMissionOne)
}

func registerMissionRouter(g *gin.RouterGroup) {
  // jwt
  g.Use(middleware.Jwt)

  // 获取单个任务
  g.GET("/id/:id", controller.GetMissionOne)

  // 添加任务
  g.POST("",
    middleware.AllowRole(middleware.RoleLeader, middleware.RolePM, middleware.RoleAdmin),
    controller.CreateMission)

  // 修改任务
  g.PUT("/id/:id",
    middleware.AllowRole(middleware.RoleLeader, middleware.RolePM, middleware.RoleAdmin),
    controller.UpdateMission)
}
