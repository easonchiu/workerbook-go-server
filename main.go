package main

import (
	"github.com/gin-gonic/gin"
	"web/middleware"
	"web/router"
	"web/db"
	"os"
)

func init() {
	db.ConnectDB()
}

func main() {

	// close db before unmount
	defer db.CloseDB()

	// initialization
	g := gin.Default() // gin.New()

	// register middleware
	middleware.Register(g)

	// register router
	router.Register(g)

	// get port args
	port := ""
	if len(os.Args) >= 2 {
		port = os.Args[1]
	}
	if port == "" {
		port = ":8080"
	} else {
		port = ":" + port
	}

	// start
	g.Run(port)

}

