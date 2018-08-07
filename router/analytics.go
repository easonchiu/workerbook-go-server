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

  // 部门列表
  g.GET("/departments", context.CreateCtx(controller.GetDepartmentListAnalytics))

  // 部门概要
  g.GET("/departments/summary/:id", context.CreateCtx(controller.GetDepartmentSummaryAnalytics))

  // 部门详情
  g.GET("/departments/detail/:id", context.CreateCtx(controller.GetDepartmentDetailAnalytics))

  // 项目列表
  g.GET("/projects", context.CreateCtx(controller.GetProjectListAnalytics))

  // 项目概要
  g.GET("/projects/summary/:id", context.CreateCtx(controller.GetProjectSummaryAnalytics))

  // 项目详情
  g.GET("/projects/detail/:id", context.CreateCtx(controller.GetProjectDetailAnalytics))

}
