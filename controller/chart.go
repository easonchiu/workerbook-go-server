package controller

import (
  "github.com/gin-gonic/gin"
  "workerbook/context"
  "workerbook/service"
)

// 获取部门成员的概要信息
func GetDepartmentUserSummary(ctx *context.New) {
  // get
  departmentId, _ := ctx.GetParam("id")

  data, err := service.GetDepartmentUserSummary(ctx, departmentId)

  if err != nil {
    ctx.Error(err)
    return
  }

  ctx.Success(gin.H{
    "data": data,
  })

}
