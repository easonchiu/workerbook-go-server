package controller

import (
  "github.com/gin-gonic/gin"
  "workerbook/context"
  "workerbook/service"
)

// 获取单个部门成员的概要信息
func GetUsersSummaryChart(ctx *context.New) {
  // get
  departmentId, _ := ctx.GetParam("id")

  data, err := service.GetUsersSummaryChartByDepartmentId(ctx, departmentId)

  if err != nil {
    ctx.Error(err)
    return
  }

  ctx.Success(gin.H{
    "data": data,
  })
}

// 获取部门列表的概要信息
func GetDepartmentsListChart(ctx *context.New) {
  // get
  skip, _ := ctx.GetQueryInt("skip")
  limit, _ := ctx.GetQueryInt("limit")

  data, err := service.GetDepartmentsListChart(ctx, skip, limit)

  if err != nil {
    ctx.Error(err)
    return
  }

  ctx.Success(gin.H{
    "data": data,
  })
}