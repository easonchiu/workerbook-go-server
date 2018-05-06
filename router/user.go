package router

import (
	"github.com/gin-gonic/gin"
	"web/controller"
	"web/middleware"
)

func registerUserRouter(g *gin.RouterGroup) {

	g.GET("/:id", middleware.Jwt, controller.GetUserInfo)

}