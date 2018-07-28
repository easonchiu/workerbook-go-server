package router

import (
  "github.com/gin-gonic/gin"
  "workerbook/conf"
  "workerbook/context"
  "workerbook/controller"
  "workerbook/middleware"
)

func registerConsoleMissionRouter(g *gin.RouterGroup) {
  // jwt
  g.Use(middleware.Jwt)

}

func registerMissionRouter(g *gin.RouterGroup) {
  // jwt
  g.Use(middleware.Jwt)

  // 获取单个任务
  g.GET("/id/:id", context.CreateCtx(controller.GetMissionOne))

  // 获取分配到自己的任务列表
  g.GET("/owns",
    // middleware.AllowRole(conf.RoleDev, conf.RoleLeader),
    context.CreateCtx(controller.GetOwnsMissionsList))

  // 添加任务
  g.POST("",
    middleware.AllowRole(conf.RoleLeader, conf.RolePM, conf.RoleAdmin),
    context.CreateCtx(controller.CreateMission))

  // 修改任务
  g.PUT("/id/:id",
    middleware.AllowRole(conf.RoleLeader, conf.RolePM, conf.RoleAdmin),
    context.CreateCtx(controller.UpdateMission))
}
