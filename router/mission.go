package router

import (
  "github.com/gin-gonic/gin"
  "workerbook/controller"
  "workerbook/middleware"
)

func registerConsoleMissionRouter(g *gin.RouterGroup) {
  // 任务列表
  g.GET("", middleware.ConsoleJwt, controller.C_GetMissionsList)

  // 获取单个任务
  g.GET("/:id", middleware.ConsoleJwt, controller.C_GetMissionOne)

  // 添加任务
  g.POST("", middleware.ConsoleJwt, controller.C_CreateMission)

  // 修改任务
  g.PUT("/:id", middleware.ConsoleJwt, controller.C_UpdateMission)
}

func registerMissionRouter(g *gin.RouterGroup) {

}
