package controller

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "workerbook/errgo"
  "workerbook/model"
  "workerbook/service"
)

// 创建用户
func C_CreateUser(c *gin.Context) {
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
func C_UpdateUser(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  id := ctx.getParam("id")

  // update
  data := bson.M{}

  if nickname := ctx.getRaw("nickname"); nickname != "" {
    data["nickname"] = nickname
  }

  if departmentId := ctx.getRaw("departmentId"); departmentId != "" {
    data["department.$id"] = departmentId
  }

  if title := ctx.getRaw("title"); title != "" {
    data["title"] = title
  }

  if role := ctx.getRawInt("role"); role != 0 {
    data["role"] = role
  }

  if status := ctx.getRawInt("status"); status != 0 {
    data["status"] = status
  }

  err := service.UpdateUser(id, data)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  ctx.Success(nil)
}

// 获取用户列表
func C_GetUsersList(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  departmentId := ctx.getQuery("departmentId")
  skip := ctx.getQueryIntDefault("skip", 0)
  limit := ctx.getQueryIntDefault("limit", 10)

  // query
  var query = bson.M{}
  if departmentId != "" {
    // check id
    if !bson.IsObjectIdHex(departmentId) {
      ctx.Error(errgo.ErrDepartmentIdError)
      return
    }
    query["department.$id"] = bson.ObjectIdHex(departmentId)
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
func C_GetUserOne(c *gin.Context) {
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

// 删除用户
func C_DelUserOne(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  id := ctx.getParam("id")

  // query
  err := service.DelUserById(id)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(nil)
}
