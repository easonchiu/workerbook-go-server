package router

import (
  "github.com/gin-gonic/gin"
  "workerbook/controller"
  "workerbook/middleware"
)

func registerMissionRouter(g *gin.RouterGroup) {
  // 任务列表
  g.GET("", middleware.Jwt, controller.GetMissionsList)

  // 获取单个任务
  g.GET("/:id", middleware.Jwt, controller.GetMissionOne)

  // 添加任务
  g.POST("", middleware.Jwt, controller.CreateMission)

  // 修改任务
  g.PUT("/:id", middleware.Jwt, controller.UpdateMission)
}
