package controller

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "workerbook/errno"
  "workerbook/model"
  "workerbook/service"
)

// login and return jwt
func UserLogin(c *gin.Context) {
  ctx := CreateCtx(c)
  defer ctx.HandleError()

  // get
  username := ctx.getRaw("username")
  password := ctx.getRaw("password")

  // check
  ctx.ErrorIfStringIsEmpty(username, errno.ErrorUsernameEmpty)
  ctx.ErrorIfStringIsEmpty(password, errno.ErrorPasswordEmpty)

  // query
  id, err := service.UserLogin(username, password)

  // check
  if err != nil {
    ctx.Error(0)
  }

  // return
  ctx.Success(gin.H{
    "data": id,
  })
}

// query users list
func GetUsersList(c *gin.Context) {
  ctx := CreateCtx(c)
  defer ctx.HandleError()

  // get
  gid := ctx.getQuery("gid")
  skip := ctx.getQueryInt("skip")
  limit := ctx.getQueryInt("limit")

  // check
  if gid != "" {
    ctx.ErrorIfStringNotObjectId(gid, errno.ErrorDepartmentIdError)
  }
  ctx.ErrorIfIntLessThen(skip, 0, errno.ErrorSkipRange)
  ctx.ErrorIfIntLessThen(limit, 0, errno.ErrorLimitRange)
  ctx.ErrorIfIntMoreThen(limit, 100, errno.ErrorLimitRange)

  // query
  usersList, err := service.GetUsersList(gid, skip, limit)

  // check
  if err != nil {
    ctx.Error(0)
  }

  // return
  ctx.Success(gin.H{
    "list": usersList,
  })
}

// query user info
func GetUserOne(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  id := ctx.getParam("id")

  if id == "my" {
    id = ctx.get("uid")
  }

  // check
  ctx.ErrorIfStringNotObjectId(id, errno.ErrorUserIdError)

  // query
  userInfo, err := service.GetUserInfoById(bson.ObjectIdHex(id))

  // check
  if err != nil {
    ctx.Error(0)
  }

  // return
  ctx.Success(gin.H{
    "data": userInfo,
  })
}

// create user
func CreateUser(c *gin.Context) {
  ctx := CreateCtx(c)
  defer ctx.HandleError()

  // get
  nickName := ctx.getRaw("nickname")
  userName := ctx.getRaw("username")
  groupId := ctx.getRaw("groupId")
  role := ctx.getRawInt("role")
  password := ctx.getRaw("password")

  // create
  data := model.User{
    NickName: nickName,
    Email:    "",
    UserName: userName,
    GroupId:  groupId,
    Role:     role,
    Mobile:   "",
    Password: password,
  }

  // insert
  err := service.CreateUser(data)

  if err != nil {
    ctx.Error(0)
  }

  // return
  ctx.Success(nil)
}
