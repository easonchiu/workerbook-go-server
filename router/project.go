package router

import (
  "github.com/gin-gonic/gin"
  "workerbook/controller"
)

func registerProjectRouter(g *gin.RouterGroup) {

  g.GET("", /*middleware.Jwt,*/ controller.GetProjectsList)

  g.POST("", /*middleware.Jwt,*/ controller.CreateProject)

}