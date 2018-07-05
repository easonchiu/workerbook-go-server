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
  defer ctx.HandleError()

  // get
  skip := ctx.getQueryInt("skip")
  limit := ctx.getQueryInt("limit")

  // check
  ctx.ErrorIfIntLessThen(skip, 0, errno.ErrorSkipRange)
  ctx.ErrorIfIntLessThen(limit, 0, errno.ErrorLimitRange)
  ctx.ErrorIfIntMoreThen(limit, 100, errno.ErrorLimitRange)

  // query
  groupsList, err := service.GetGroupsList(skip, limit)

  // check
  if err != nil {
    ctx.Error(errno.ErrorDepartmentNotFound)
  }

  // query
  count, err := service.GetCountOfGroup()

  // check
  if err != nil {
    ctx.Error(errno.ErrorDepartmentNotFound)
  }

  // return
  ctx.Success(gin.H{
    "list": groupsList,
    "skip": skip,
    "limit": limit,
    "count": count,
  })
}

// 获取单个部门的信息
func GetDepartmentOne(c *gin.Context) {
  ctx := CreateCtx(c)
  defer ctx.HandleError()

  // get
  id := ctx.getParam("id")

  // check
  ctx.ErrorIfStringNotObjectId(id, errno.ErrorDepartmentIdError)

  // query
  groupInfo, err := service.GetGroupInfoById(bson.ObjectIdHex(id))

  // check
  if err != nil {
    ctx.Error(errno.ErrorDepartmentNotFound)
  }

  // return
  ctx.Success(gin.H{
    "data": groupInfo,
  })
}

// 创建部门
func CreateDepartment(c *gin.Context) {
  ctx := CreateCtx(c)
  defer ctx.HandleError()

  // get
  name := ctx.getRaw("name")

  // check
  ctx.ErrorIfStringIsEmpty(name, errno.ErrorDepartmentNameEmpty)

  // create
  data := model.Group{
    Name: name,
  }

  // insert
  err := service.CreateGroup(data)

  // check
  if err != nil {
    ctx.Error(errno.ErrorCreateDepartmentFailed)
  }

  // return
  ctx.Success(nil)
}
