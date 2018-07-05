package router

import (
  `github.com/gin-gonic/gin`
  `workerbook/controller`
  `workerbook/middleware`
)

func registerGroupRouter(g *gin.RouterGroup) {

  g.GET("", middleware.Jwt, controller.GetDepartmentsList)

  g.GET("/:id", middleware.Jwt, controller.GetDepartmentOne)

  g.POST("", middleware.Jwt, controller.CreateDepartment)

}
