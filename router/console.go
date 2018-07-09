package router

import (
  "github.com/gin-gonic/gin"
  "workerbook/controller"
  "workerbook/middleware"
)

func registerConsoleRouter(g  *gin.RouterGroup) {

  // 用户列表
  g.GET("/users", middleware.ConsoleJwt, controller.GetUsersList)

  // 添加用户
  g.POST("/users", middleware.ConsoleJwt, controller.CreateUser)

  // 修改用户
  g.PUT("/users", middleware.ConsoleJwt, controller.UpdateUser)

}
