package controller

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "workerbook/service"
  "workerbook/util"
)

func GetDepartmentsList(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  skip, _ := ctx.getQueryInt("skip")
  limit, _ := ctx.getQueryInt("limit")

  // query
  data, err := service.GetDepartmentsList(skip, limit, bson.M{})

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