package router

import (
	"github.com/gin-gonic/gin"
	"workerbook/controller"
)

func registerGroupRouter(g *gin.RouterGroup) {

	g.GET("", /*middleware.Jwt,*/ controller.GetGroupsList)

	g.GET("/:id", /*middleware.Jwt,*/ controller.GetGroupInfo)

	g.POST("", /*middleware.Jwt,*/ controller.CreateGroup)

}