package middleware

import (
  "github.com/gin-gonic/gin"
  "github.com/jwt-go"
  "regexp"
  "workerbook/conf"
  "workerbook/context"
  "workerbook/errgo"
)

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

      departmentId := tokenMap[conf.OWN_DEPARTMENT_ID].(string)
      userId := tokenMap[conf.OWN_USER_ID].(string)
      role := int(tokenMap[conf.OWN_ROLE].(float64))

      // 验证token内的关键字段
      errgo.ErrorIfStringNotObjectId(departmentId, errgo.ErrServerError)
      errgo.ErrorIfStringNotObjectId(userId, errgo.ErrServerError)

      if err := errgo.PopError(); err != nil {
        get := errgo.Get(err)
        c.JSON(200, gin.H{
          "msg":  get.Message,
          "code": get.Code,
          "data": nil,
        })
        errgo.ClearErrorStack()
        return
      }

      c.Set(conf.OWN_DEPARTMENT_ID, departmentId)
      c.Set(conf.OWN_USER_ID, userId)
      c.Set(conf.OWN_ROLE, role)

      c.Next()
    } else {
      ctx := context.NewBaseCtx(c)
      ctx.Error(errgo.ErrUserReLogin)
      c.Abort()
    }
  } else {
    ctx := context.NewBaseCtx(c)
    ctx.Error(errgo.ErrUserReLogin)
    c.Abort()
  }
}

// 权限控制(匹配中的用户允许访问)
func AllowRole(roles ... int) func(c *gin.Context) {
  return func(c *gin.Context) {
    role := c.GetInt(conf.OWN_ROLE)

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

    ctx := context.NewBaseCtx(c)
    ctx.Error(errgo.ErrForbidden)
    c.Abort()
  }
}
