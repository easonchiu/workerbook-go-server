package main

import (
	"github.com/gin-gonic/gin"
	"workerbook/middleware"
	"workerbook/router"
	"workerbook/db"
	"flag"
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
	// e.g.  go run main.go --port=:8080
	port := ""
	flag.StringVar(&port, "port", ":8080", "port addr")
	flag.Parse()

	// start
	g.Run(port)

}

