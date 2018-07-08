package router

import (
  "github.com/gin-gonic/gin"
  "workerbook/controller"
  "workerbook/middleware"
)

func registerConsoleRouter(g  *gin.RouterGroup) {

  // 用户列表
  g.GET("/users", middleware.ConsoleJwt, controller.GetUsersList)

}
