package controller

import (
  "errors"
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "strconv"
  "workerbook/model"
  "workerbook/service"
)

// login and return jwt
func UserLogin(c *gin.Context) {
  ctx := CreateCtx(c)

  username := ctx.getRaw("username")
  password := ctx.getRaw("password")

  // check the raw data.
  if username == "" {
    ctx.Error(errors.New("用户名不能为空"), 1)
    return
  } else if password == "" {
    ctx.Error(errors.New("密码不能为空"), 1)
    return
  }

  // query user info from database.
  id, err := service.UserLogin(username, password)

  if err != nil {
    ctx.Error(err, 1)
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

  skip, _ := c.GetQuery("skip")
  limit, _ := c.GetQuery("limit")

  intSkip, err := strconv.Atoi(skip)

  if err != nil {
    intSkip = 0
  }

  intLimit, err := strconv.Atoi(limit)

  if err != nil {
    intLimit = 10
  }

  usersList, err := service.GetUsersList(intSkip, intLimit)
  if err != nil {
    ctx.Error(err, 1)
    return
  }

  ctx.Success(gin.H{
    "list": usersList,
  })
}

// query user info
func GetUserInfo(c *gin.Context) {
  ctx := CreateCtx(c)

  id := ctx.getParam("id")

  if !bson.IsObjectIdHex(id) {
    ctx.Error(errors.New("无效的id号"), 1)
    return
  }

  userInfo, err := service.GetUserInfoById(bson.ObjectIdHex(id))
  if err != nil {
    ctx.Error(err, 1)
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
    Mobile:   ctx.getRaw("mobile"),
    Password: ctx.getRaw("password"),
  }

  err := service.CreateUser(data)
  if err != nil {
    ctx.Error(err, 1)
    return
  }

  ctx.Success(nil)
}
