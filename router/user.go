package router

import (
  `github.com/gin-gonic/gin`
  `workerbook/controller`
  `workerbook/middleware`
)

func registerUserRouter(g *gin.RouterGroup) {

  g.GET("", middleware.Jwt, controller.GetUsersList)

  g.GET("/:id", middleware.Jwt, controller.GetUserOne)

  // g.POST("/:id/dailies/today/items", middleware.Jwt, controller.CreateMyTodayDailyItem)

  // g.DELETE("/:id/dailies/today/items/:itemId", middleware.Jwt, controller.DeleteUserTodayDailyItem)

  // g.GET("/:id/dailies/today", middleware.Jwt, controller.GetMyTodayDaily)

  g.POST("", middleware.Jwt, controller.CreateUser)

}
