package main

import (
  "flag"
  "github.com/gin-gonic/gin"
  "workerbook/mongo"
  "workerbook/router"
)

func init() {
  mongo.ConnectDB()
}

func main() {

  // close db before unmount
  defer mongo.CloseDB()

  // initialization
  // Default With the Logger and Recovery middleware already attached
  g := gin.Default() // gin.New()

  // register middleware
  // middleware.Register(g)

  // register router
  router.Register(g)

  // get port args
  // e.g.  go run main.go --port=:8080
  port := ""
  flag.StringVar(&port, "port", "8080", "port addr")
  flag.Parse()

  // start
  g.Run(":" + port)

}
