package controller

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "workerbook/conf"
  "workerbook/context"
  "workerbook/errgo"
  "workerbook/model"
  "workerbook/service"
)

// 用户登录
func UserLogin(ctx *context.New) {

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
func GetProfile(ctx *context.New) {

  // get
  id, _ := ctx.Get("USER_ID")

  // query
  user, err := service.GetUserInfoById(ctx, id)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(gin.H{
    "data": func() gin.H {
      data := user.GetMap(model.REMEMBER, "nickname", "id", "role", "status", "title")
      department, err := service.FindDepartmentRef(ctx, &user.Department)
      if err == nil {
        data["department"] = department.GetMap(model.REMEMBER, "name", "id")
      }
      return data
    }(),
  })
}

// 获取自己及部下人员信息
func GetSubUsersList(ctx *context.New) {

  // get
  role, _ := ctx.GetInt(conf.OWN_ROLE)
  departmentId, _ := ctx.Get(conf.OWN_DEPARTMENT_ID)
  projectId, _ := ctx.GetQuery("projectId")

  // check
  if !bson.IsObjectIdHex(departmentId) {
    ctx.Error(errgo.ErrForbidden)
    return
  }

  if role == conf.RoleAdmin || role == conf.RolePM {
    project, err := service.GetProjectInfoById(ctx, projectId)

    if err != nil {
      ctx.Error(err)
      return
    }

    var projectIdList []bson.ObjectId

    for _, item := range project.Departments {
      projectIdList = append(projectIdList, item.Id.(bson.ObjectId))
    }

    users, err := service.GetUsersList(ctx, 0, 0, bson.M{
      "department.$id": bson.M{
        "$in": projectIdList,
      },
    })

    if err != nil {
      ctx.Error(err)
    } else {
      ctx.Success(gin.H{
        "data": users.Each(func(item *model.User) gin.H {
          return item.GetMap(model.REMEMBER, "id", "nickname", "title")
        }),
      })
    }
    return
  } else if role == conf.RoleLeader {
    users, err := service.GetUsersList(ctx, 0, 0, bson.M{
      "department.$id": bson.ObjectIdHex(departmentId),
    })

    if err != nil {
      ctx.Error(err)
    } else {
      ctx.Success(gin.H{
        "data": users.Each(func(item *model.User) gin.H {
          return item.GetMap(model.REMEMBER, "id", "nickname", "title")
        }),
      })
    }
    return
  }

  // 未匹配中角色
  ctx.Error(errgo.ErrForbidden)
}

// 获取用户列表
func GetUsersList(ctx *context.New) {

  // get
  departmentId, didExist := ctx.GetQuery("departmentId")
  skip := ctx.GetQueryIntDefault("skip", 0)
  limit := ctx.GetQueryIntDefault("limit", 10)

  // query
  var query = bson.M{}
  if didExist {
    // check id
    if !bson.IsObjectIdHex(departmentId) {
      ctx.Error(errgo.ErrDepartmentIdError)
      return
    }
    query["department.$id"] = bson.ObjectIdHex(departmentId)
  }

  data, err := service.GetUsersList(ctx, skip, limit, query)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(gin.H{
    "data": data.Each(func(item *model.User) gin.H {
      return item.GetMap("editor", "editTime", "exist", "department", "username")
    }),
  })
}
