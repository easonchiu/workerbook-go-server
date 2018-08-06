package controller

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "workerbook/conf"
  "workerbook/context"
  "workerbook/model"
  "workerbook/service"
)

// 获取用户相关的项目列表
// 如果是开发者或是leader，则拉取和自己部门有关的任务
func GetProjectsList(ctx *context.New) {

  // get
  departmentId, _ := ctx.Get(conf.OWN_DEPARTMENT_ID)
  role, _ := ctx.GetInt(conf.OWN_ROLE)

  skip := ctx.GetQueryIntDefault("skip", 0)
  limit := ctx.GetQueryIntDefault("limit", 9)

  // query
  query := bson.M{
    "status": 1,
  }

  if role == conf.RoleDev || role == conf.RoleLeader {
    query["departments.$id"] = bson.ObjectIdHex(departmentId)
  }

  data, err := service.GetProjectsList(ctx, skip, limit, query)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(gin.H{
    "data": data.Each(func(item *model.Project) gin.H {
      each := item.GetMap("editor", "exist", "editTime")

      // progress
      if len(item.Missions) > 0 {
        progress := service.GetProjectProgress(ctx, item.Id)
        each["progress"] = progress
      }

      // departments
      var departments []gin.H

      for _, ref := range item.Departments {
        department, err := service.FindDepartmentRef(ctx, &ref)
        if err == nil {
          departments = append(departments, department.GetMap(model.REMEMBER, "id", "name"))
        }
      }

      each["departments"] = departments

      // missions
      var missions []gin.H

      for _, ref := range item.Missions {
        mission, err := service.FindMissionRef(ctx, &ref)
        if err == nil {
          m := mission.GetMap("editor", "editTime", "exist", "project")
          user, err := service.FindUserRef(ctx, &mission.User)
          if err == nil {
            m["user"] = user.GetMap(model.REMEMBER, "id", "nickname", "exist", "status", "title", "role")
            missions = append(missions, m)
          }
        }
      }

      each["missions"] = missions

      return each
    }),
  })
}
