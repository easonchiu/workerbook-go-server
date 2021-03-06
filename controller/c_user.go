package controller

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "workerbook/context"
  "workerbook/errgo"
  "workerbook/model"
  "workerbook/service"
)

// 创建用户
func C_CreateUser(ctx *context.New) {

  // get
  nickname, _ := ctx.GetRaw("nickname")
  username, _ := ctx.GetRaw("username")
  departmentId, _ := ctx.GetRaw("departmentId")
  title, _ := ctx.GetRaw("title")
  role, _ := ctx.GetRawInt("role")
  password, _ := ctx.GetRaw("password")

  // create
  data := model.User{
    NickName: nickname,
    Email:    "",
    UserName: username,
    Title:    title,
    Role:     role,
    Mobile:   "",
    Password: password,
  }

  // insert
  err := service.CreateUser(ctx, data, departmentId)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(nil)
}

// 修改用户
func C_UpdateUser(ctx *context.New) {

  // get
  id, _ := ctx.GetParam("id")

  // update
  data := bson.M{}

  if nickname, ok := ctx.GetRaw("nickname"); ok {
    data["nickname"] = nickname
  }

  if departmentId, ok := ctx.GetRaw("departmentId"); ok {
    data["department.$id"] = departmentId
  }

  if title, ok := ctx.GetRaw("title"); ok {
    data["title"] = title
  }

  if role, ok := ctx.GetRawInt("role"); ok {
    data["role"] = role
  }

  if status, ok := ctx.GetRawInt("status"); ok {
    data["status"] = status
  }

  err := service.UpdateUser(ctx, id, data)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  ctx.Success(nil)
}

// 获取用户列表
func C_GetUsersList(ctx *context.New) {

  // get
  departmentId, didExist := ctx.GetQuery("departmentId")
  skip := ctx.GetQueryIntDefault("skip", 0)
  limit := ctx.GetQueryIntDefault("limit", 10)

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

  data, err := service.GetUsersList(ctx, skip, limit, query)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(gin.H{
    "data": data.Each(func(item *model.User) gin.H {
      each := item.GetMap()
      department, err := service.FindDepartmentRef(ctx, &item.Department)
      if err == nil {
        each["department"] = department.GetMap("userCount", "createTime")
      }
      return each
    }),
  })
}

// 获取用户信息
func C_GetUserOne(ctx *context.New) {

  // get
  id, _ := ctx.GetParam("id")

  // query
  user, err := service.GetUserInfoById(ctx, id)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(gin.H{
    "data": user.GetMap("username"),
  })
}

// 删除用户
func C_DelUserOne(ctx *context.New) {

  // get
  id, _ := ctx.GetParam("id")

  // query
  err := service.UpdateUser(ctx, id, bson.M{
    "exist": false,
  })

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(nil)
}
