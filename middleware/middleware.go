package middleware

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"web/conf"
)

func Register(g *gin.Engine) {
	g.Use(print)

	g.Use(mongo)
}

// makes the `db` object available for each handler
func mongo(c *gin.Context) {
	s := conf.Session.Clone()

	defer s.Close()

	c.Set("db", s.DB(conf.Mongo.Database))

	c.Next()
}

func print(c *gin.Context) {
	fmt.Println("its a middleware")

	c.Next()

	fmt.Println("middleware after Next()")
}