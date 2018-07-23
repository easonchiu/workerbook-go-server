package router

import (
  "github.com/gin-gonic/gin"
  "workerbook/conf"
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
  g.GET("/id/:id", controller.GetMissionOne)

  // 添加任务
  g.POST("",
    middleware.AllowRole(conf.RoleLeader, conf.RolePM, conf.RoleAdmin),
    controller.CreateMission)

  // 修改任务
  g.PUT("/id/:id",
    middleware.AllowRole(conf.RoleLeader, conf.RolePM, conf.RoleAdmin),
    controller.UpdateMission)
}
