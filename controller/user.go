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

  // query
  var query model.User
  if departmentId != "" {
    // check id
    if !bson.IsObjectIdHex(departmentId) {
      ctx.Error(errgo.ErrDepartmentIdError)
      return
    }
    query.Department = mgo.DBRef{
      Id:         bson.ObjectIdHex(departmentId),
      Collection: model.DepartmentCollection,
      Database:   conf.DBName,
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

  // query
  userInfo, err := service.GetUserInfoById(id)

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

  // query
  userInfo, err := service.GetUserInfoById(id)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  delete(userInfo, "username")

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

  // create
  data := model.User{
    NickName:   nickname,
    Email:      "",
    UserName:   username,
    Title:      title,
    Role:       role,
    Mobile:     "",
    Password:   password,
    CreateTime: time.Now(),
  }

  // insert
  err := service.CreateUser(data, departmentId)

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

  // update
  data := model.User{
    NickName: nickname,
    Title:    title,
    Role:     role,
    Status:   status,
  }

  err := service.UpdateUser(id, data, departmentId)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  ctx.Success(nil)
}
