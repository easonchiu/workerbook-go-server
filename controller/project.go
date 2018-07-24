package controller

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "workerbook/context"
  "workerbook/service"
)

// 获取用户相关的项目列表
func GetProjectsList(c *gin.Context) {
  ctx, err := context.CreateCtx(c)
  defer ctx.Close()

  if err != nil {
    ctx.Error(err)
    return
  }

  // get
  departmentId, _ := c.Get("DEPARTMENT_ID")

  // query
  data, err := service.GetProjectsList(ctx, 0, 0, bson.M{
    "departments.$id": bson.ObjectIdHex(departmentId.(string)),
  }, "missions", "departments", "user")

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
