package controller

import (
  "github.com/gin-gonic/gin"
  "workerbook/errno"
  "workerbook/model"
  "workerbook/service"
)

// 获取部门列表
// func GetDepartmentsList(c *gin.Context) {
//   ctx := CreateCtx(c)
//
//   // get
//   skip := ctx.getQueryInt("skip")
//   limit := ctx.getQueryInt("limit")
//
//   // check
//   ctx.ErrorIfIntLessThen(skip, 0, errno.ErrSkipRange)
//   ctx.ErrorIfIntLessThen(limit, 0, errno.ErrLimitRange)
//   ctx.ErrorIfIntMoreThen(limit, 100, errno.ErrLimitRange)
//
//   // query
//   groupsList, err := service.GetGroupsList(skip, limit)
//
//   // check
//   if err != nil {
//     ctx.Error(errno.ErrorDepartmentNotFound)
//   }
//
//   // query
//   count, err := service.GetCountOfGroup()
//
//   // check
//   if err != nil {
//     ctx.Error(errno.ErrorDepartmentNotFound)
//   }
//
//   // return
//   ctx.Success(gin.H{
//     "list": groupsList,
//     "skip": skip,
//     "limit": limit,
//     "count": count,
//   })
// }

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

  // create
  data := model.Department{
    Name: name,
    UserCount: 0,
  }

  // insert
  err := service.CreateDepartment(data)

  // check
  if err != nil {
    ctx.Error(err)
  }

  // return
  ctx.Success(nil)
}
