package controller

import "github.com/gin-gonic/gin"

// 获取任务列表
func GetMissionsList(c *gin.Context) {
  ctx := CreateCtx(c)

  ctx.Success(nil)
}