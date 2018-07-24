package controller

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "workerbook/context"
  "workerbook/service"
  "workerbook/util"
)

func GetDepartmentsList(c *gin.Context) {
  ctx, err := context.CreateCtx(c)
  defer ctx.Close()

  if err != nil {
    ctx.Error(err)
    return
  }

  // get
  skip, _ := ctx.GetQueryInt("skip")
  limit, _ := ctx.GetQueryInt("limit")

  // query
  data, err := service.GetDepartmentsList(ctx, skip, limit, bson.M{})

  // 过滤
  util.ForgetArr(data["list"], "createTime")

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