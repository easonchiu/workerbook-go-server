package controller

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "workerbook/errgo"
  "workerbook/model"
  "workerbook/service"
)

// 获取部门列表
func GetDepartmentsList(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  skip := ctx.getQueryInt("skip")
  limit := ctx.getQueryInt("limit")

  // check
  if limit != 0 {
    errgo.ErrorIfIntLessThen(skip, 0, errgo.ErrSkipRange)
    errgo.ErrorIfIntLessThen(limit, 1, errgo.ErrLimitRange)
    errgo.ErrorIfIntMoreThen(limit, 100, errgo.ErrLimitRange)
  }

  // query
  data, err := service.GetDepartmentsList(skip, limit, nil)

  // check
  if err != nil {
    ctx.Error(errgo.ErrDepartmentNotFound)
    return
  }

  // return
  ctx.Success(gin.H{
    "data": data,
  })
}

// 创建部门
func CreateDepartment(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  name := ctx.getRaw("name")

  // check
  errgo.ErrorIfStringIsEmpty(name, errgo.ErrDepartmentNameEmpty)

  if errgo.HandleError(ctx.Error) {
    return
  }

  // create
  data := model.Department{
    Name:      name,
    UserCount: 0,
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

// 获取单个部门
func GetDepartmentOne(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  id := ctx.getParam("id")

  // check
  errgo.ErrorIfStringNotObjectId(id, errgo.ErrDepartmentIdError)

  if errgo.HandleError(ctx.Error) {
    return
  }

  // query
  departmentInfo, err := service.GetDepartmentInfoById(bson.ObjectIdHex(id))

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
func UpdateDepartment(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  id := ctx.getParam("id")
  name := ctx.getRaw("name")

  // check
  errgo.ErrorIfStringNotObjectId(id, errgo.ErrDepartmentIdError)
  errgo.ErrorIfStringIsEmpty(name, errgo.ErrDepartmentNameEmpty)

  if errgo.HandleError(ctx.Error) {
    return
  }

  // update
  data := model.Department{
    Name: name,
  }
  err := service.UpdateDepartment(bson.ObjectIdHex(id), data)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  ctx.Success(nil)
}
