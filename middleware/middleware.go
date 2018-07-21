package middleware

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "regexp"
  "workerbook/controller"
  "workerbook/errgo"
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
      c.Set("DEPARTMENT_ID", "5b424feeaea6f431c2655006")
      c.Set("UID", token)
      c.Set("ROLE", RoleAdmin)
      c.Next()
    } else {
      ctx := controller.CreateCtx(c)
      ctx.Error(errgo.ErrUserReLogin)
      c.Abort()
    }
  } else {
    ctx := controller.CreateCtx(c)
    ctx.Error(errgo.ErrUserReLogin)
    c.Abort()
  }
}

// 权限控制(匹配中的用户允许访问)
// 1: 开发者 2: 部门管理者 3: 观察员 4: 项目管理者 99: 管理员
const (
  RoleAdmin  = 99 // 管理员
  RolePM     = 4  // 项目管理者
  RoleOB     = 3  // 观察员
  RoleLeader = 2  // 部门管理者
  RoleDev    = 1  // 开发者
)

func AllowRole(roles ... int) func(c *gin.Context) {
  return func(c *gin.Context) {
    role := c.GetInt("ROLE")

    exist := false
    for _, i := range roles {
      if i == role {
        exist = true
        continue
      }
    }

    if exist {
      c.Next()
      return
    }

    ctx := controller.CreateCtx(c)
    ctx.Error(errgo.ErrForbidden)
    c.Abort()
  }
}
