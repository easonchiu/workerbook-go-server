package controller

import (
  "errors"
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "workerbook/model"
  "workerbook/service"
)

// 获取分组列表
func GetGroupsList(c *gin.Context) {
  ctx := CreateCtx(c)

  skip := ctx.getQueryInt("skip")
  limit := ctx.getQueryInt("limit")

  groupsList, err := service.GetGroupsList(skip, limit)
  if err != nil {
    ctx.Error(err, 1)
    return
  }

  ctx.Success(gin.H{
    "list": groupsList,
  })
}

// 获取单个分组的信息
func GetGroupOne(c *gin.Context) {
  ctx := CreateCtx(c)

  id := ctx.getParam("id")

  if !bson.IsObjectIdHex(id) {
    ctx.Error(errors.New("无效的分组ID"), 1)
    return
  }

  groupInfo, err := service.GetGroupInfoById(bson.ObjectIdHex(id))
  if err != nil {
    ctx.Error(err, 1)
    return
  }

  ctx.Success(gin.H{
    "data": groupInfo,
  })
}

// 创建分组
func CreateGroup(c *gin.Context) {
  ctx := CreateCtx(c)

  name := ctx.getRaw("name")

  data := model.Group{
    Name: name,
  }

  err := service.CreateGroup(data)
  if err != nil {
    ctx.Error(err, 1)
    return
  }

  ctx.Success(nil)
}
