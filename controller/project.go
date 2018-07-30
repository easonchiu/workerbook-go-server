package controller

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "workerbook/context"
  "workerbook/model"
  "workerbook/service"
)

// 获取用户相关的项目列表
func GetProjectsList(ctx *context.New) {

  // get
  departmentId, _ := ctx.Get("DEPARTMENT_ID")

  // query
  data, err := service.GetProjectsList(ctx, 0, 0, bson.M{
    "departments.$id": bson.ObjectIdHex(departmentId),
  })

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(gin.H{
    "data": data.Each(func(item model.Project) gin.H {
      each := item.GetMap()

      // departments
      var departments []gin.H

      for _, ref := range item.Departments {
        department, err := service.FindDepartmentRef(ctx, &ref)
        if err == nil {
          departments = append(departments, department.GetMap())
        }
      }

      each["departments"] = departments

      // missions
      var missions []gin.H

      for _, ref := range item.Missions {
        mission, err := service.FindMissionRef(ctx, &ref)
        if err == nil {
          m := mission.GetMap("editor", "editTime", "exist")
          user, err := service.FindUserRef(ctx, &mission.User)
          if err == nil {
            m["user"] = user.GetMap("username", "department", "editor", "editTime")
            missions = append(missions, m)
          }
        }
      }

      each["missions"] = missions

      return each
    }),
  })
}
