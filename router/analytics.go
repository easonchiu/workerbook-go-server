package router

import (
  "github.com/gin-gonic/gin"
  "workerbook/context"
  "workerbook/controller"
  "workerbook/middleware"
)

func registerAnalyticsRouter(g *gin.RouterGroup) {
  // jwt
  g.Use(middleware.Jwt)

  // 整体部门概要（各个部门的人数，任务数，任务完成情况）
  g.GET("/departments", context.CreateCtx(controller.GetDepartmentsAnalytics))

  // 部门成员概要（每个成员的任务数，任务完成情况）
  g.GET("/departments/summary/:id", context.CreateCtx(controller.GetDepartmentOneAnalytics))

  // 整体项目概要
  g.GET("/projects", context.CreateCtx(controller.GetProjectsAnalytics))

  // 单个项目的任务概要(即项目的完成指标)
  g.GET("/projects/summary/:id", context.CreateCtx(controller.GetProjectOneAnalytics))

}
