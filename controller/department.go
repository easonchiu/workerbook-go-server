package controller

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "workerbook/context"
  "workerbook/model"
  "workerbook/service"
)

// 获取部门列表
func GetDepartmentsList(ctx *context.New) {

  // get
  skip, _ := ctx.GetQueryInt("skip")
  limit, _ := ctx.GetQueryInt("limit")

  // query
  data, err := service.GetDepartmentsList(ctx, skip, limit, bson.M{})

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(gin.H{
    "data": data.Each(func(item model.Department) gin.H {
      return item.GetMap("createTime", "editor", "editTime", "exist")
    }),
  })
}
