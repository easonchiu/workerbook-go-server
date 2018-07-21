package controller

import (
  "github.com/gin-gonic/gin"
  "workerbook/service"
  "workerbook/util"
)

// 用户登录
func UserLogin(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  username, _ := ctx.getRaw("username")
  password, _ := ctx.getRaw("password")

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
  id, _ := ctx.get("UID")

  // query
  userInfo, err := service.GetUserInfoById(id)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  util.Forget(userInfo, "username")

  // return
  ctx.Success(gin.H{
    "data": userInfo,
  })
}