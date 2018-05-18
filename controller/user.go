package controller

import (
  "errors"
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "workerbook/model"
  "workerbook/service"
)

// login and return jwt
func UserLogin(c *gin.Context) {
  ctx := CreateCtx(c)

  username := ctx.getRaw("username").(string)
  password := ctx.getRaw("password").(string)

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

  gid := ctx.getQuery("gid").(string)
  skip := ctx.getQuery("skip", true).(int)
  limit := ctx.getQuery("limit", true).(int)

  usersList, err := service.GetUsersList(gid, skip, limit)
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

  id := ctx.getParam("id").(string)

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
    NickName: ctx.getRaw("nickname").(string),
    Email:    ctx.getRaw("email").(string),
    UserName: ctx.getRaw("username").(string),
    Gid:      ctx.getRaw("gid").(string),
    Role:     ctx.getRaw("role", true).(int),
    Mobile:   ctx.getRaw("mobile").(string),
    Password: ctx.getRaw("password").(string),
  }

  err := service.CreateUser(data)
  if err != nil {
    ctx.Error(err, 1)
    return
  }

  ctx.Success(nil)
}
