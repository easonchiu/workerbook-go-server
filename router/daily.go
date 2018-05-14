package router

import (
	"github.com/gin-gonic/gin"
	"workerbook/controller"
)

func registerDailyRouter(g *gin.RouterGroup) {

	g.GET("/list", /*middleware.Jwt,*/ controller.GetDailiesList)

	g.GET("/detail/:id", /*middleware.Jwt,*/ controller.GetDailyInfo)

	g.GET("/today/:uid", controller.GetTodayDaily)

	g.PUT("/item", /*middleware.Jwt,*/ controller.CreateDailyItem)

	g.DELETE("/item/:id", controller.DeleteDailyItem)

}