package router

import (
  "github.com/gin-gonic/gin"
  "workerbook/context"
  "workerbook/controller"
)

func registerAnalyticsRouter(g *gin.RouterGroup) {

  // 整体部门概要（各个部门的人数，任务数，任务完成情况）
  g.GET("/departments", context.CreateCtx(controller.GetDepartmentsAnalytics))

  // 部门成员概要（每个成员的任务数，任务完成情况）
  g.GET("/departments/id/:id", context.CreateCtx(controller.GetDepartmentOneAnalytics))

  // 整体项目该要
  g.GET("/projects", context.CreateCtx(controller.GetProjectsAnalytics))

}