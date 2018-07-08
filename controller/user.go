package controller

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "workerbook/errno"
  "workerbook/model"
  "workerbook/service"
)

// 用户登录
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

// 获取用户列表
func GetUsersList(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  isConsole := ctx.getBool("isConsole")
  departmentId := ctx.getQuery("departmentId")
  skip := ctx.getQueryIntDefault("skip", 0)
  limit := ctx.getQueryIntDefault("limit", 10)

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
  if isConsole == true {
    usersList, err := service.GetConsoleUsersList(departmentId, skip, limit)

    // check
    if err != nil {
      ctx.Error(err)
      return
    }

    // return
    ctx.Success(gin.H{
      "list": usersList,
    })
  } else {
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
}

// 获取用户信息
func GetUserOne(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  id := ctx.getParam("id")

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

// 获取我的信息
func GetProfile(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  id := ctx.get("uid")

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

// 创建用户
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
