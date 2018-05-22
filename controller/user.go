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

  username := ctx.getRaw("username")
  password := ctx.getRaw("password")

  // check the raw data.
  if username == "" {
    ctx.Error("用户名不能为空", 1)
    return
  } else if password == "" {
    ctx.Error("密码不能为空", 1)
    return
  }

  // query user info from database.
  id, err := service.UserLogin(username, password)

  if err != nil {
    ctx.Error("登录失败", 1)
    return
  } else {
    ctx.Success(gin.H{
      "data": id,
    })
  }
}

// query users list
func GetUsersList(c *gin.Context) {
  ctx := CreateCtx(c)

  gid := ctx.getQuery("gid")
  skip := ctx.getQueryInt("skip")
  limit := ctx.getQueryInt("limit")

  usersList, err := service.GetUsersList(gid, skip, limit)
  if err != nil {
    ctx.Error(err.Error(), 1)
    return
  }

  ctx.Success(gin.H{
    "list": usersList,
  })
}

// query user info
func GetUserOne(c *gin.Context) {
  ctx := CreateCtx(c)

  c.SetCookie("test", "value", 10000, "/", "0.0.0.0", true, true)

  id := ctx.getParam("id")

  if id == "my" {
    id = ctx.get("uid")
  }

  if !bson.IsObjectIdHex(id) {
    ctx.Error("无效的用户ID", 1)
    return
  }

  userInfo, err := service.GetUserInfoById(bson.ObjectIdHex(id))
  if err != nil {
    ctx.Error("获取用户信息失败", 1)
    return
  }

  ctx.Success(gin.H{
    "data": userInfo,
  })
}

// create user
func CreateUser(c *gin.Context) {
  ctx := CreateCtx(c)

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
    ctx.Error("创建用户失败", 1)
    return
  }

  ctx.Success(nil)
}
