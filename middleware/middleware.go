package middleware

import (
  `errors`
  "fmt"
  "github.com/gin-gonic/gin"
  `gopkg.in/mgo.v2/bson`
  `regexp`
  "workerbook/controller"
)

func Register(g *gin.Engine) {
  // g.Use(log)
}

// check up json web token
func Jwt(c *gin.Context) {
  auth, token := c.Request.Header.Get("authorization"), ""

  jwtReg := regexp.MustCompile(`^Bearer\s\S+$`)

  if jwtReg.MatchString(auth) {
    token = auth[len("Bearer "):]

    // check up your token here...
    if bson.IsObjectIdHex(token) {
      c.Set("uid", token)

      c.Next()
    } else {
      ctx := controller.CreateCtx(c)
      ctx.Error(errors.New("bad token."), 401)
      return
    }

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
