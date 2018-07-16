package router

import (
  `github.com/gin-gonic/gin`
  `workerbook/controller`
  `workerbook/middleware`
)

func registerConsoleUserRouter(g *gin.RouterGroup) {
  // 获取用户列表
  g.GET("", middleware.ConsoleJwt, controller.C_GetUsersList)

  // 获取单个用户
  g.GET("/:id", middleware.ConsoleJwt, controller.C_GetUserOne)

  // 添加用户
  g.POST("", middleware.ConsoleJwt, controller.C_CreateUser)

  // 修改用户
  g.PUT("/:id", middleware.ConsoleJwt, controller.C_UpdateUser)

  // 修改用户
  g.DELETE("/:id", middleware.ConsoleJwt, controller.C_DelUserOne)
}

func registerUserRouter(g *gin.RouterGroup) {
  // 登录
  g.POST("/login", middleware.Jwt, controller.UserLogin)

  // 获取个人信息
  g.GET("/profile", middleware.Jwt, controller.GetProfile)

  // g.POST("/:id/dailies/today/items", middleware.Jwt, controller.CreateMyTodayDailyItem)

  // g.DELETE("/:id/dailies/today/items/:itemId", middleware.Jwt, controller.DeleteUserTodayDailyItem)

  // g.GET("/:id/dailies/today", middleware.Jwt, controller.GetMyTodayDaily)

}
