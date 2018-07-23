package main

import (
  "flag"
  "github.com/gin-gonic/gin"
  "workerbook/db"
  "workerbook/router"
)

func init() {
  db.ConnectMgoDB()
  db.InitRedisPool()
}

func main() {

  // close db before un-mount
  defer db.CloseMgoDB()

  defer db.RedisPool.Close()

  // initialization
  // Default With the Logger and Recovery middleware already attached
  g := gin.Default() // gin.New()

  // gin.SetMode(gin.ReleaseMode)

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
