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

  // get
  username := ctx.getRaw("username")
  password := ctx.getRaw("password")

  // check
  ctx.ErrorIfStringIsEmpty(username, errno.ErrUsernameEmpty)
  ctx.ErrorIfStringIsEmpty(password, errno.ErrPasswordEmpty)

  if ctx.HandleErrorIf() {
    return
  }

  // query
  id, err := service.UserLogin(username, password)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(gin.H{
    "data": id,
  })
}

// query users list
func GetUsersList(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  departmentId := ctx.getQuery("departmentId")
  skip := ctx.getQueryInt("skip")
  limit := ctx.getQueryInt("limit")

  // check
  if departmentId != "" {
    ctx.ErrorIfStringNotObjectId(departmentId, errno.ErrDepartmentIdError)
  }
  ctx.ErrorIfIntLessThen(skip, 0, errno.ErrSkipRange)
  ctx.ErrorIfIntLessThen(limit, 0, errno.ErrLimitRange)
  ctx.ErrorIfIntMoreThen(limit, 100, errno.ErrLimitRange)

  if ctx.HandleErrorIf() {
    return
  }

  // query
  usersList, err := service.GetUsersList(departmentId, skip, limit)

  // check
  if err != nil {
    ctx.Error(err)
    return
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
  ctx.ErrorIfStringNotObjectId(id, errno.ErrUserIdError)

  if ctx.HandleErrorIf() {
    return
  }

  // query
  userInfo, err := service.GetUserInfoById(bson.ObjectIdHex(id))

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(gin.H{
    "data": userInfo,
  })
}

// create user
func CreateUser(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  nickName := ctx.getRaw("nickname")
  userName := ctx.getRaw("username")
  departmentId := ctx.getRaw("departmentId")
  role := ctx.getRawInt("role")
  password := ctx.getRaw("password")

  // check
  ctx.ErrorIfStringIsEmpty(userName, errno.ErrUsernameEmpty)
  ctx.ErrorIfStringIsEmpty(password, errno.ErrPasswordEmpty)
  ctx.ErrorIfStringIsEmpty(nickName, errno.ErrNicknameEmpty)
  ctx.ErrorIfStringNotObjectId(departmentId, errno.ErrDepartmentIdError)

  if ctx.HandleErrorIf() {
    return
  }

  // create
  data := model.User{
    NickName:     nickName,
    Email:        "",
    UserName:     userName,
    DepartmentId: departmentId,
    Role:         role,
    Mobile:       "",
    Password:     password,
  }

  // insert
  err := service.CreateUser(data)

  // check
  if err != nil {
    ctx.Error(err)
  }

  // update count
  service.UpdateDepartmentCount(departmentId)

  // return
  ctx.Success(nil)
}
