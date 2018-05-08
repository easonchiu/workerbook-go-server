package router

import (
	"github.com/gin-gonic/gin"
	"workerbook/controller"
)

func registerUserRouter(g *gin.RouterGroup) {

	g.GET("", /*middleware.Jwt,*/ controller.GetUsersList)

	g.GET("/:id", /*middleware.Jwt,*/ controller.GetUserInfo)

	g.POST("", /*middleware.Jwt,*/ controller.CreateUser)

}