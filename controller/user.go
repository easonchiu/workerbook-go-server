package controller

import (
  "fmt"
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "workerbook/context"
  "workerbook/errgo"
  "workerbook/model"
  "workerbook/service"
)

// 用户登录
func UserLogin(c *gin.Context) {
  ctx, err := context.CreateCtx(c)
  defer ctx.Close()

  if err != nil {
    ctx.Error(err)
    return
  }

  // get
  username, _ := ctx.GetRaw("username")
  password, _ := ctx.GetRaw("password")

  // query
  token, err := service.UserLogin(ctx, username, password)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(gin.H{
    "data": token,
  })
}

// 获取我的信息
func GetProfile(c *gin.Context) {
  ctx, err := context.CreateCtx(c)
  defer ctx.Close()

  if err != nil {
    ctx.Error(err)
    return
  }

  // get
  id, _ := ctx.Get("UID")

  // query
  user, err := service.GetUserInfoById(ctx, id)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(gin.H{
    "data": user.GetMap("department", "username"),
  })
}

// 获取自己及部下人员信息
func GetSubUsersList(c *gin.Context) {
  ctx, err := context.CreateCtx(c)
  defer ctx.Close()

  if err != nil {
    ctx.Error(err)
    return
  }

  // get
  // role, _ := ctx.GetInt("ROLE")
  departmentId, _ := ctx.Get("DEPARTMENT_ID")
  projectId, _ := ctx.GetQuery("projectId")

  // check
  if !bson.IsObjectIdHex(departmentId) {
    ctx.Error(errgo.ErrForbidden)
    return
  }

  if false {
  // if role == conf.RoleAdmin || role == conf.RolePM {
    project, err := service.GetProjectInfoById(ctx, projectId)

    if err != nil {
      ctx.Error(err)
      return
    }

    fmt.Println(project)

    // var projectIdList []bson.ObjectId
    //
    // for _, item := range project.Departments {
    //   if id, ok := item["id"]; ok {
    //     projectIdList = append(projectIdList, id.(bson.ObjectId))
    //   }
    // }
    //
    // users, err := service.GetUsersList(ctx, 0, 0, bson.M{
    //   "department.$id": bson.M{
    //     "$in": projectIdList,
    //   },
    // })

    if err != nil {
      ctx.Error(err)
    } else {
      ctx.Success(gin.H{
        "data": nil,
      })
    }
    return
  // } else if role == conf.RoleLeader {
  } else if true {
    users, err := service.GetUsersList(ctx, 0, 0, bson.M{
      "department.$id": bson.ObjectIdHex(departmentId),
    })

    if err != nil {
      ctx.Error(err)
    } else {
      ctx.Success(gin.H{
        "data": users.Each(func(item model.User) gin.H {
          return item.GetMap("username")
        }),
      })
    }
    return
  }

  // 未匹配中角色
  ctx.Error(errgo.ErrForbidden)
}
