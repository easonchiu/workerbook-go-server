package router

import (
  "github.com/gin-gonic/gin"
  "workerbook/conf"
  "workerbook/middleware"
)

func Register(g *gin.Engine) {

  // 管理后台的路由（管理后台必须jwt通过并只开放给pm与admin）
  console := g.Group("/console")
  console.Use(middleware.Jwt)
  console.Use(middleware.AllowRole(conf.RolePM, conf.RoleAdmin))

  // 注册用户相关的路由
  registerUserRouter(g.Group("/users"))
  registerConsoleUserRouter(console.Group("/users"))

  // 注册部门相关的路由
  registerDepartmentRouter(g.Group("/departments"))
  registerConsoleDepartmentRouter(console.Group("/departments"))

  // 注册日报相关的路由
  registerDailyRouter(g.Group("/dailies"))

  // 注册项目相关的路由
  registerProjectRouter(g.Group("/projects"))
  registerConsoleProjectRouter(console.Group("/projects"))

  // 注册任务相关的路由
  registerMissionRouter(g.Group("/missions"))
  registerConsoleMissionRouter(console.Group("/missions"))

  // 注册数据统计相关的路由
  registerAnalyticsRouter(g.Group("/analytics"))

}
