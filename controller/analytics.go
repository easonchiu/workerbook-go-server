package controller

import (
  "github.com/gin-gonic/gin"
  "workerbook/context"
  "workerbook/service"
)

// 获取单个部门成员的概要信息
func GetDepartmentOneAnalytics(ctx *context.New) {
  // get
  departmentId, _ := ctx.GetParam("id")

  data, err := service.GetDepartmentAnalysisById(ctx, departmentId)

  if err != nil {
    ctx.Error(err)
    return
  }

  ctx.Success(gin.H{
    "data": data,
  })
}

// 获取部门列表的概要信息
func GetDepartmentsAnalytics(ctx *context.New) {
  // get
  skip, _ := ctx.GetQueryInt("skip")
  limit, _ := ctx.GetQueryInt("limit")

  data, err := service.GetDepartmentsAnalysis(ctx, skip, limit)

  if err != nil {
    ctx.Error(err)
    return
  }

  ctx.Success(gin.H{
    "data": data,
  })
}

// 获取项目列表的概要信息
func GetProjectsAnalytics(ctx *context.New) {
  // get
  skip, _ := ctx.GetQueryInt("skip")
  limit, _ := ctx.GetQueryInt("limit")

  data, err := service.GetProjectsAnalysis(ctx, skip, limit)

  if err != nil {
    ctx.Error(err)
    return
  }

  ctx.Success(gin.H{
    "data": data,
  })
}