package controller

import (
  "errors"
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "strconv"
  "workerbook/model"
  "workerbook/service"
)

// 获取分组列表
func GetGroupsList(c *gin.Context) {
  ctx := CreateCtx(c)

  skip, _ := c.GetQuery("skip")
  limit, _ := c.GetQuery("limit")

  intSkip, err := strconv.Atoi(skip)

  if err != nil {
    intSkip = 0
  }

  intLimit, err := strconv.Atoi(limit)

  if err != nil {
    intLimit = 10
  }

  groupsList, err := service.GetGroupsList(intSkip, intLimit)
  if err != nil {
    ctx.Error(err, 1)
    return
  }

  ctx.Success(gin.H{
    "list": groupsList,
  })
}

// 获取单个分组的信息
func GetGroupInfo(c *gin.Context) {
  ctx := CreateCtx(c)

  id := ctx.getParam("id")

  if !bson.IsObjectIdHex(id) {
    ctx.Error(errors.New("无效的id号"), 1)
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

  data := model.Group{
    Name: ctx.getRaw("name"),
  }

  err := service.CreateGroup(data)
  if err != nil {
    ctx.Error(err, 1)
    return
  }

  ctx.Success(nil)
}
