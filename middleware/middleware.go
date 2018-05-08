package middleware

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"workerbook/controller"
)

func Register(g *gin.Engine) {
	// g.Use(log)
}

// check up json web token
func Jwt(c *gin.Context) {
	auth, prefix, token := c.Request.Header.Get("authorization"), "Bearer ", ""

	if len(auth) > len(prefix) {
		token = auth[len(prefix):]

		fmt.Println(token)
		// check up your token here...

		c.Next()
	} else {
		ctx := controller.CreateCtx(c)
		ctx.Forbidden()
	}
}

// print user agent
func log(c *gin.Context) {
	fmt.Println(" >>> UserAgent is: ", c.Request.UserAgent())
	c.Next()
}