package router

import (
  "github.com/gin-gonic/gin"
  "workerbook/controller"
  "workerbook/middleware"
)

func registerMissionRouter(g *gin.RouterGroup) {
  // 任务列表
  g.GET("", middleware.Jwt, controller.GetDepartmentsList)

  // 获取单个任务
  g.GET("/:id", middleware.Jwt, controller.GetDepartmentOne)

  // 添加任务
  g.POST("", middleware.Jwt, controller.CreateDepartment)

  // 修改任务
  g.PUT("/:id", middleware.Jwt, controller.UpdateDepartment)
}
