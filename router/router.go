package router

import "github.com/gin-gonic/gin"

func Register(g *gin.Engine) {

	// 注册用户相关的路由
	registerUserRouter(g.Group("/user"))

}
