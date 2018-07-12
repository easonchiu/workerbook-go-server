package controller

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "time"
  "workerbook/conf"
  "workerbook/errgo"
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
  errgo.ErrorIfStringIsEmpty(username, errgo.ErrUsernameEmpty)
  errgo.ErrorIfStringIsEmpty(password, errgo.ErrPasswordEmpty)

  if errgo.HandleError(ctx.Error) {
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
  departmentId := ctx.getQuery("departmentId")
  skip := ctx.getQueryIntDefault("skip", 0)
  limit := ctx.getQueryIntDefault("limit", 10)

  // check
  if departmentId != "" {
    errgo.ErrorIfStringNotObjectId(departmentId, errgo.ErrDepartmentIdError)
  }
  errgo.ErrorIfIntLessThen(skip, 0, errgo.ErrSkipRange)
  errgo.ErrorIfIntLessThen(limit, 1, errgo.ErrLimitRange)
  errgo.ErrorIfIntMoreThen(limit, 100, errgo.ErrLimitRange)

  if errgo.HandleError(ctx.Error) {
    return
  }

  // query
  var query bson.M
  if departmentId != "" {
    query = bson.M{
      "departmentId": departmentId,
    }
  }

  data, err := service.GetUsersList(skip, limit, query)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(gin.H{
    "data": data,
  })
}

// 获取用户信息
func GetUserOne(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  id := ctx.getParam("id")

  // check
  errgo.ErrorIfStringNotObjectId(id, errgo.ErrUserIdError)

  if errgo.HandleError(ctx.Error) {
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
  errgo.ErrorIfStringNotObjectId(id, errgo.ErrUserIdError)

  if errgo.HandleError(ctx.Error) {
    return
  }

  // query
  userInfo, err := service.GetUserInfoById(bson.ObjectIdHex(id))
  delete(userInfo, "username")

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
  nickname := ctx.getRaw("nickname")
  username := ctx.getRaw("username")
  departmentId := ctx.getRaw("departmentId")
  title := ctx.getRaw("title")
  role := ctx.getRawInt("role")
  password := ctx.getRaw("password")

  // check
  errgo.ErrorIfStringIsEmpty(nickname, errgo.ErrNicknameEmpty)
  errgo.ErrorIfLenLessThen(nickname, 2, errgo.ErrNicknameTooShort)
  errgo.ErrorIfLenMoreThen(nickname, 14, errgo.ErrNicknameTooLong)
  errgo.ErrorIfStringNotObjectId(departmentId, errgo.ErrDepartmentIdError)
  errgo.ErrorIfStringIsEmpty(title, errgo.ErrUserTitleIsEmpty)
  errgo.ErrorIfLenMoreThen(title, 14, errgo.ErrUserTitleTooLong)
  if role != 1 && role != 2 && role != 3 {
    ctx.Error(errgo.ErrUserRoleError)
    return
  }
  errgo.ErrorIfStringIsEmpty(username, errgo.ErrUsernameEmpty)
  errgo.ErrorIfLenLessThen(username, 6, errgo.ErrUsernameTooShort)
  errgo.ErrorIfLenMoreThen(username, 14, errgo.ErrUsernameTooLong)
  errgo.ErrorIfStringIsEmpty(password, errgo.ErrPasswordEmpty)

  if errgo.HandleError(ctx.Error) {
    return
  }

  // create
  data := model.User{
    NickName: nickname,
    Email:    "",
    UserName: username,
    Department: mgo.DBRef{
      Id:         bson.ObjectIdHex(departmentId),
      Collection: model.DepartmentCollection,
      Database:   conf.DBName,
    },
    Title:      title,
    Role:       role,
    Mobile:     "",
    Password:   password,
    CreateTime: time.Now(),
  }

  // insert
  err := service.CreateUser(data)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(nil)
}

// 修改用户
func UpdateUser(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  id := ctx.getParam("id")
  nickname := ctx.getRaw("nickname")
  departmentId := ctx.getRaw("departmentId")
  title := ctx.getRaw("title")
  role := ctx.getRawInt("role")
  status := ctx.getRawInt("status")

  // check
  errgo.ErrorIfStringIsEmpty(nickname, errgo.ErrNicknameEmpty)
  errgo.ErrorIfLenLessThen(nickname, 2, errgo.ErrNicknameTooShort)
  errgo.ErrorIfLenMoreThen(nickname, 14, errgo.ErrNicknameTooLong)
  errgo.ErrorIfStringNotObjectId(departmentId, errgo.ErrDepartmentIdError)
  errgo.ErrorIfStringIsEmpty(title, errgo.ErrUserTitleIsEmpty)
  errgo.ErrorIfLenMoreThen(title, 14, errgo.ErrUserTitleTooLong)
  if role != 1 && role != 2 && role != 3 {
    ctx.Error(errgo.ErrUserRoleError)
    return
  }

  if errgo.HandleError(ctx.Error) {
    return
  }

  // update
  data := model.User{
    NickName: nickname,
    Department: mgo.DBRef{
      Id:         bson.ObjectIdHex(departmentId),
      Collection: model.DepartmentCollection,
      Database:   conf.DBName,
    },
    Title:  title,
    Role:   role,
    Status: status,
  }

  err := service.UpdateUser(bson.ObjectIdHex(id), data)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  ctx.Success(nil)
}
