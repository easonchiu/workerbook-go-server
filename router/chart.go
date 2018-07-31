package router

import (
  "github.com/gin-gonic/gin"
  "workerbook/context"
  "workerbook/controller"
)

func registerChartRouter(g *gin.RouterGroup) {

  // 整体部门概要（各个部门的人数，任务数，任务完成情况）
  g.GET("/departments/summary", context.CreateCtx(controller.GetDepartmentsListChart))

  // 部门成员概要（每个成员的任务数，任务完成情况）
  g.GET("/departments/summary/:id", context.CreateCtx(controller.GetUsersSummaryChart))

}