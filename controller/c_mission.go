package controller

import (
  "github.com/gin-gonic/gin"
  "workerbook/service"
)

// 获取单个任务
func C_GetMissionOne(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  id, _ := ctx.getParam("id")

  // query
  missionInfo, err := service.GetMissionInfoById(id)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(gin.H{
    "data": missionInfo,
  })
}

// 获取任务列表
func C_GetMissionsList(c *gin.Context) {
  ctx := CreateCtx(c)

  // return
  ctx.Success(nil)
}

