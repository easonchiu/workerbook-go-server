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
  ctx.ErrorIfIntLessThen(skip, 0, errno.ErrSkipRange)
  ctx.ErrorIfIntLessThen(limit, 0, errno.ErrLimitRange)
  ctx.ErrorIfIntMoreThen(limit, 100, errno.ErrLimitRange)

  // query
  data, err := service.GetDepartmentsList(skip, limit)

  // check
  if err != nil {
    ctx.Error(errno.ErrDepartmentNotFound)
  }

  // return
  ctx.Success(gin.H{
    "data": data,
  })
}

// 获取单个部门的信息
// func GetDepartmentOne(c *gin.Context) {
//   ctx := CreateCtx(c)
//   defer ctx.HandleError()
//
//   // get
//   id := ctx.getParam("id")
//
//   // check
//   ctx.ErrorIfStringNotObjectId(id, errno.ErrorDepartmentIdError)
//
//   // query
//   groupInfo, err := service.GetGroupInfoById(bson.ObjectIdHex(id))
//
//   // check
//   if err != nil {
//     ctx.Error(errno.ErrorDepartmentNotFound)
//   }
//
//   // return
//   ctx.Success(gin.H{
//     "data": groupInfo,
//   })
// }

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
  err := service.UpdateDepartment(bson.ObjectIdHex(id), bson.M{
    "name": name,
  })

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  ctx.Success(nil)
}
