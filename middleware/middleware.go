package middleware

import (
	"github.com/gin-gonic/gin"
	"fmt"
)

func Register(g *gin.Engine) {
	g.Use(log)
}

func log(c *gin.Context) {
	fmt.Println(" >>> UserAgent is: ", c.Request.UserAgent())
	c.Next()
}