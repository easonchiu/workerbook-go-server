package router

import (
  "github.com/gin-gonic/gin"
  "workerbook/context"
  "workerbook/controller"
  "workerbook/middleware"
)

func registerDailyRouter(g *gin.RouterGroup) {
  // jwt
  g.Use(middleware.Jwt)

  // 创建日报
  g.POST("/today", context.CreateCtx(controller.CreateDaily))

  // 获取我今天的日报
  g.GET("/today", context.CreateCtx(controller.GetTodayDaily))

  // 获取日报列表
  g.GET("", context.CreateCtx(controller.GetDailiesListByDay))

  // 更新日报
  g.PUT("/today", context.CreateCtx(controller.UpdateDaily))

  // 删除一条日报数据
  g.DELETE("/today", context.CreateCtx(controller.DelDaily))
}
