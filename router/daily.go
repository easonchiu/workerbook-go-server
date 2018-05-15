package router

import (
	"github.com/gin-gonic/gin"
	"workerbook/controller"
)

func registerDailyRouter(g *gin.RouterGroup) {

	g.GET("", /*middleware.Jwt,*/ controller.GetDailiesList)

	g.GET("/:id", /*middleware.Jwt,*/ controller.GetDailyInfo)

	g.PUT("/item", /*middleware.Jwt,*/ controller.CreateTodayDailyItem)

	g.DELETE("/item/:itemId", controller.DeleteTodayDailyItem)

}
