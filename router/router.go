package router

import (
  `github.com/gin-gonic/gin`
)

func Register(g *gin.Engine) {

  // 注册用户相关的路由
  registerUserRouter(g.Group("/users"))
  registerConsoleUserRouter(g.Group("/console/users"))

  // 注册部门相关的路由
  registerDepartmentRouter(g.Group("/departments"))
  registerConsoleDepartmentRouter(g.Group("/console/departments"))

  // 注册日报相关的路由
  // registerDailyRouter(g.Group("/dailies"))

  // 注册项目相关的路由
  registerProjectRouter(g.Group("/projects"))
  registerConsoleProjectRouter(g.Group("/console/projects"))

  // 注册任务相关的路由
  registerMissionRouter(g.Group("/missions"))
  registerConsoleMissionRouter(g.Group("/console/missions"))

}
