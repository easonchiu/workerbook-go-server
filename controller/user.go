package controller

import (
  "github.com/gin-gonic/gin"
  "workerbook/service"
)

// 用户登录
func UserLogin(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  username := ctx.getRaw("username")
  password := ctx.getRaw("password")

  // query
  id, err := service.UserLogin(username, password)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(gin.H{
    "data": id,
  })
}

// 获取我的信息
func GetProfile(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  id := ctx.get("UID")

  // query
  userInfo, err := service.GetUserInfoById(id)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  delete(userInfo, "username")

  // return
  ctx.Success(gin.H{
    "data": userInfo,
  })
}