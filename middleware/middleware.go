package middleware

import (
  "github.com/gin-gonic/gin"
  "github.com/jwt-go"
  "regexp"
  "workerbook/conf"
  "workerbook/context"
  "workerbook/errgo"
)

func Register(g *gin.Engine) {

}

// check up json web token
func Jwt(c *gin.Context) {
  headerAuth, headerJwt := c.Request.Header.Get("authorization"), ""

  jwtReg := regexp.MustCompile(`^Bearer\s\S+$`)

  if jwtReg.MatchString(headerAuth) {
    headerJwt = headerAuth[len("Bearer "):]

    token, _ := jwt.Parse(headerJwt, func(t *jwt.Token) (interface{}, error) {
      return conf.JwtSecret, nil
    })

    if token.Valid {
      tokenMap := token.Claims.(jwt.MapClaims)

      c.Set("DEPARTMENT_ID", tokenMap["departmentId"])
      c.Set("UID", tokenMap["id"])
      c.Set("ROLE", int(tokenMap["role"].(float64)))

      c.Next()
    } else {
      ctx := context.CreateBaseCtx(c)
      ctx.Error(errgo.ErrUserReLogin)
      c.Abort()
    }
  } else {
    ctx := context.CreateBaseCtx(c)
    ctx.Error(errgo.ErrUserReLogin)
    c.Abort()
  }
}

// 权限控制(匹配中的用户允许访问)
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

    ctx := context.CreateBaseCtx(c)
    ctx.Error(errgo.ErrForbidden)
    c.Abort()
  }
}
