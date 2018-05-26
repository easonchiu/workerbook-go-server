package controller

import (
  `github.com/gin-gonic/gin`
  `gopkg.in/mgo.v2/bson`
  `workerbook/model`
  `workerbook/service`
)

// login and return jwt
func UserLogin(c *gin.Context) {
  ctx := CreateCtx(c)
  defer ctx.handleErrorIfPanic()

  username := ctx.getRaw("username")
  password := ctx.getRaw("password")
  
  // check
  ctx.PanicIfStringIsEmpty(username, "用户名不能为空")
  ctx.PanicIfStringIsEmpty(password, "密码不能为空")

  // query user info from database.
  id, err := service.UserLogin(username, password)

  if err != nil {
    panic("登录失败")
  } else {
    ctx.Success(gin.H{
      "data": id,
    })
  }
}

// query users list
func GetUsersList(c *gin.Context) {
  ctx := CreateCtx(c)
  defer ctx.handleErrorIfPanic()

  gid := ctx.getQuery("gid")
  skip := ctx.getQueryInt("skip")
  limit := ctx.getQueryInt("limit")

  // check
  if gid != "" {
    ctx.PanicIfStringNotObjectId(gid, "分组ID错误")
  }
  ctx.PanicIfIntLessThen(skip, 0, "Skip不能小于0")
  ctx.PanicIfIntLessThen(limit, 0, "Limit不能小于0")
  ctx.PanicIfIntMoreThen(limit, 100, "Limit不能大于100")

  usersList, err := service.GetUsersList(gid, skip, limit)

  if err != nil {
    panic("获取用户列表失败")
  }

  ctx.Success(gin.H{
    "list": usersList,
  })
}

// query user info
func GetUserOne(c *gin.Context) {
  ctx := CreateCtx(c)
  defer ctx.handleErrorIfPanic()

  id := ctx.getParam("id")

  if id == "my" {
    id = ctx.get("uid")
  }

  ctx.PanicIfStringNotObjectId(id, "无效的用户ID")

  userInfo, err := service.GetUserInfoById(bson.ObjectIdHex(id))

  if err != nil {
    panic("获取用户信息失败")
  }

  ctx.Success(gin.H{
    "data": userInfo,
  })
}

// create user
func CreateUser(c *gin.Context) {
  ctx := CreateCtx(c)
  defer ctx.handleErrorIfPanic()

  data := model.User{
    NickName: ctx.getRaw("nickname"),
    Email:    ctx.getRaw("email"),
    UserName: ctx.getRaw("username"),
    Gid:      ctx.getRaw("gid"),
    Role:     ctx.getRawInt("role"),
    Mobile:   ctx.getRaw("mobile"),
    Password: ctx.getRaw("password"),
  }

  err := service.CreateUser(data)

  if err != nil {
    panic("创建用户失败")
  }

  ctx.Success(nil)
}
