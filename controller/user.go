package controller

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "workerbook/conf"
  "workerbook/errgo"
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
  data, err := service.UserLogin(username, password)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(data)
}

// 获取我的信息
func GetProfile(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  id, _ := ctx.get("UID")

  // query
  userInfo, err := service.GetUserInfoById(id, "department")

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

// 获取自己及部下人员信息
func GetSubUsersList(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  role, _ := ctx.getInt("ROLE")
  departmentId, _ := ctx.get("DEPARTMENT_ID")
  projectId, _ := ctx.getQuery("projectId")

  // check
  if !bson.IsObjectIdHex(departmentId) {
    ctx.Error(errgo.ErrForbidden)
    return
  }


  if role == conf.RoleAdmin || role == conf.RolePM {
    project, err := service.GetProjectInfoById(projectId)

    if err != nil {
      ctx.Error(err)
      return
    }

    var projectIdList []bson.ObjectId

    for _, item := range project["departments"].([]gin.H) {
      if id, ok := item["id"]; ok {
        projectIdList = append(projectIdList, id.(bson.ObjectId))
      }
    }

    users, err := service.GetUsersList(0,0, bson.M{
      "department.$id": bson.M{
        "$in": projectIdList,
      },
    })

    if err != nil {
      ctx.Error(err)
    } else {
      ctx.Success(gin.H{
        "data": users,
      })
    }
    return
  } else if role == conf.RoleLeader {
    users, err := service.GetUsersList(0,0, bson.M{
      "department.$id": bson.ObjectIdHex(departmentId),
    })

    if err != nil {
      ctx.Error(err)
    } else {
      ctx.Success(gin.H{
        "data": users,
      })
    }
    return
  }

  // 未匹配中角色
  ctx.Error(errgo.ErrForbidden)
}