package controller

import (
  `github.com/gin-gonic/gin`
  `gopkg.in/mgo.v2/bson`
  `workerbook/model`
  `workerbook/service`
)

// 获取分组列表
func GetGroupsList(c *gin.Context) {
  ctx := CreateCtx(c)
  defer ctx.handleErrorIfPanic()

  skip := ctx.getQueryInt("skip")
  limit := ctx.getQueryInt("limit")

  // check
  ctx.PanicIfIntLessThen(skip, 0, "Skip不能小于0")
  ctx.PanicIfIntLessThen(limit, 0, "Limit不能小于0")
  ctx.PanicIfIntMoreThen(limit, 100, "Limit不能大于100")

  groupsList, err := service.GetGroupsList(skip, limit)

  if err != nil {
    panic("获取分组失败")
  }

  count, err := service.GetCountOfGroup()

  if err != nil {
    panic("获取分组失败")
  }

  ctx.Success(gin.H{
    "list": groupsList,
    "skip": skip,
    "limit": limit,
    "count": count,
  })
}

// 获取单个分组的信息
func GetGroupOne(c *gin.Context) {
  ctx := CreateCtx(c)
  defer ctx.handleErrorIfPanic()

  id := ctx.getParam("id")

  // check
  ctx.PanicIfStringNotObjectId(id, "无效的分组ID")

  groupInfo, err := service.GetGroupInfoById(bson.ObjectIdHex(id))

  if err != nil {
    panic("获取分组失败")
  }

  ctx.Success(gin.H{
    "data": groupInfo,
  })
}

// 创建分组
func CreateGroup(c *gin.Context) {
  ctx := CreateCtx(c)
  defer ctx.handleErrorIfPanic()

  name := ctx.getRaw("name")

  // check
  ctx.PanicIfStringIsEmpty(name, "分组名不能为空")

  data := model.Group{
    Name: name,
  }

  err := service.CreateGroup(data)
  if err != nil {
    panic("创建分组失败")
  }

  ctx.Success(nil)
}
