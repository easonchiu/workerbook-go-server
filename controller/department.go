package controller

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "workerbook/errno"
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
    ctx.ErrorIfIntLessThen(skip, 0, errno.ErrSkipRange)
    ctx.ErrorIfIntLessThen(limit, 1, errno.ErrLimitRange)
    ctx.ErrorIfIntMoreThen(limit, 100, errno.ErrLimitRange)
  }

  // query
  data, err := service.GetDepartmentsList(skip, limit, nil)

  // check
  if err != nil {
    ctx.Error(errno.ErrDepartmentNotFound)
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
  ctx.ErrorIfStringIsEmpty(name, errno.ErrDepartmentNameEmpty)

  if ctx.HandleErrorIf() {
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
  ctx.ErrorIfStringNotObjectId(id, errno.ErrDepartmentIdError)

  if ctx.HandleErrorIf() {
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
  ctx.ErrorIfStringNotObjectId(id, errno.ErrDepartmentIdError)
  ctx.ErrorIfStringIsEmpty(name, errno.ErrDepartmentNameEmpty)

  if ctx.HandleErrorIf() {
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
