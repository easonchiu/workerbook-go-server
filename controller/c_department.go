package controller

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "workerbook/context"
  "workerbook/model"
  "workerbook/service"
)

// 创建部门
func C_CreateDepartment(ctx *context.New) {

  // get
  name, _ := ctx.GetRaw("name")

  // create
  data := model.Department{
    Name: name,
  }

  // insert
  err := service.CreateDepartment(ctx, data)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(nil)
}

// 删除单个部门
func C_DelDepartmentOne(ctx *context.New) {

  // get
  id, _ := ctx.GetParam("id")

  // query
  err := service.DelDepartmentById(ctx, id)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(nil)
}

// 获取部门列表
func C_GetDepartmentsList(ctx *context.New) {

  // get
  skip, _ := ctx.GetQueryInt("skip")
  limit, _ := ctx.GetQueryInt("limit")

  // query
  data, err := service.GetDepartmentsList(ctx, skip, limit, bson.M{})

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(gin.H{
    "data": data.Each(func(item model.Department) gin.H {
      return item.GetMap()
    }),
  })
}

// 获取单个部门
func C_GetDepartmentOne(ctx *context.New) {

  // get
  id, _ := ctx.GetParam("id")

  // query
  department, err := service.GetDepartmentInfoById(ctx, id)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(gin.H{
    "data": department.GetMap(),
  })
}

// 修改部门
func C_UpdateDepartment(ctx *context.New) {

  // get
  id, _ := ctx.GetParam("id")

  data := bson.M{}

  if name, ok := ctx.GetRaw("name"); ok {
    data["name"] = name
  }

  // update
  err := service.UpdateDepartment(ctx, id, data)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  ctx.Success(nil)
}
