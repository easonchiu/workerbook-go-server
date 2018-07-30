package router

import (
  "github.com/gin-gonic/gin"
  "workerbook/context"
  "workerbook/controller"
)

func registerChartRouter(g *gin.RouterGroup) {

  // 整体部门概要（各个部门的任务数，完成情况以及任务饱和指标）
  g.GET("/department/summary")

  // 部门成员概要（每个成员的任务数，完成情况）
  g.GET("/department/summary/:id", context.CreateCtx(controller.GetDepartmentUserSummary))

}