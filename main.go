package main

import (
	"github.com/gin-gonic/gin"
	"workerbook/middleware"
	"workerbook/router"
	"workerbook/db"
	"os"
)

func init() {
	db.ConnectDB()
}

func main() {

	// close db before unmount
	defer db.CloseDB()

	// initialization
	// Default With the Logger and Recovery middleware already attached
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

