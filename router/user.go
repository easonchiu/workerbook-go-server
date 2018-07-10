package router

import (
  `github.com/gin-gonic/gin`
  `workerbook/controller`
  `workerbook/middleware`
)

func registerUserRouter(g *gin.RouterGroup) {

  // 获取用户列表
  g.GET("", middleware.Jwt, controller.GetUsersList)

  // 获取单个用户
  g.GET("/:id", middleware.Jwt, controller.GetUserOne)

  // 添加用户
  g.POST("", middleware.ConsoleJwt, controller.CreateUser)

  // 修改用户
  g.PUT("/:id", middleware.ConsoleJwt, controller.UpdateUser)

  // g.POST("/:id/dailies/today/items", middleware.Jwt, controller.CreateMyTodayDailyItem)

  // g.DELETE("/:id/dailies/today/items/:itemId", middleware.Jwt, controller.DeleteUserTodayDailyItem)

  // g.GET("/:id/dailies/today", middleware.Jwt, controller.GetMyTodayDaily)

}
