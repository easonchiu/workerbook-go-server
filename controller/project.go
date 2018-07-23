package controller

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "workerbook/service"
)

// 获取用户相关的项目列表
func GetProjectsList(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  departmentId, _ := c.Get("DEPARTMENT_ID")

  // query
  data, err := service.GetProjectsList(0, 0, bson.M{
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