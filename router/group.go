package router

import (
	`github.com/gin-gonic/gin`
	`workerbook/controller`
	`workerbook/middleware`
)

func registerGroupRouter(g *gin.RouterGroup) {

	g.GET("", middleware.Jwt, controller.GetGroupsList)

	g.GET("/:id", middleware.Jwt, controller.GetGroupOne)

	g.POST("", middleware.Jwt, controller.CreateGroup)

}