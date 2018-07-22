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
  nickname, _ := ctx.getRaw("nickname")
  username, _ := ctx.getRaw("username")
  departmentId, _ := ctx.getRaw("departmentId")
  title, _ := ctx.getRaw("title")
  role, _ := ctx.getRawInt("role")
  password, _ := ctx.getRaw("password")

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
  id, _ := ctx.getParam("id")

  // update
  data := bson.M{}

  if nickname, ok := ctx.getRaw("nickname"); ok {
    data["nickname"] = nickname
  }

  if departmentId, ok := ctx.getRaw("departmentId"); ok {
    data["department.$id"] = departmentId
  }

  if title, ok := ctx.getRaw("title"); ok {
    data["title"] = title
  }

  if role, ok := ctx.getRawInt("role"); ok {
    data["role"] = role
  }

  if status, ok := ctx.getRawInt("status"); ok {
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
  departmentId, didExist := ctx.getQuery("departmentId")
  skip := ctx.getQueryIntDefault("skip", 0)
  limit := ctx.getQueryIntDefault("limit", 10)

  // query
  var query = bson.M{}
  if didExist {
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
  id, _ := ctx.getParam("id")

  // query
  userInfo, err := service.GetUserInfoById(id, "department")

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
  id, _ := ctx.getParam("id")

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
