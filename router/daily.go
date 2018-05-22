package router

import (
  `github.com/gin-gonic/gin`
  `workerbook/controller`
  `workerbook/middleware`
)

func registerDailyRouter(g *gin.RouterGroup) {

  g.GET("", controller.GetDailiesList)

  g.GET("/:id", middleware.Jwt, controller.GetDailyOne)

}
