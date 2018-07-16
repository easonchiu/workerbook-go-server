package controller

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "workerbook/model"
  "workerbook/service"
)

// 创建部门
func C_CreateDepartment(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  name := ctx.getRaw("name")

  // create
  data := model.Department{
    Name:       name,
  }

  // insert
  err := service.CreateDepartment(data)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(nil)
}

// 删除单个部门
func C_DelDepartmentOne(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  id := ctx.getParam("id")

  // query
  err := service.DelDepartmentById(id)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(nil)
}

// 获取部门列表
func C_GetDepartmentsList(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  skip := ctx.getQueryInt("skip")
  limit := ctx.getQueryInt("limit")

  // query
  data, err := service.GetDepartmentsList(skip, limit, bson.M{})

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

// 获取单个部门
func C_GetDepartmentOne(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  id := ctx.getParam("id")

  // query
  departmentInfo, err := service.GetDepartmentInfoById(id)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(gin.H{
    "data": departmentInfo,
  })
}

// 修改部门
func C_UpdateDepartment(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  id := ctx.getParam("id")
  name := ctx.getRaw("name")

  // update
  err := service.UpdateDepartment(id, bson.M{
    "name": name,
  })

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  ctx.Success(nil)
}
