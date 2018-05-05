package router

import (
	"github.com/gin-gonic/gin"
	"web/controller"
)

func registerUserRouter(g *gin.RouterGroup) {

	g.GET("/:id", controller.GetUserInfo)

}