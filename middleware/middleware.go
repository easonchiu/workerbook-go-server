package middleware

import (
  `github.com/gin-gonic/gin`
  `gopkg.in/mgo.v2/bson`
  `regexp`
  `workerbook/controller`
  "workerbook/errno"
)

func Register(g *gin.Engine) {

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
      ctx.Error(errno.ErrUserReLogin)
      c.Abort()
    }
  } else {
    ctx := controller.CreateCtx(c)
    ctx.Error(errno.ErrUserReLogin)
    c.Abort()
  }
}

// check up json web token
func ConsoleJwt(c *gin.Context) {
  auth, token := c.Request.Header.Get("authorization"), ""

  jwtReg := regexp.MustCompile(`^Bearer\s\S+$`)

  if jwtReg.MatchString(auth) {
    token = auth[len("Bearer "):]

    // check up your token here...
    if bson.IsObjectIdHex(token) {
      c.Set("uid", token)
      c.Set("isConsole", true)
      c.Next()
    } else {
      ctx := controller.CreateCtx(c)
      ctx.Error(errno.ErrUserReLogin)
      c.Abort()
    }
  } else {
    ctx := controller.CreateCtx(c)
    ctx.Error(errno.ErrUserReLogin)
    c.Abort()
  }
}
