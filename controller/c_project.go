package controller

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "workerbook/context"
  "workerbook/model"
  "workerbook/service"
)

// 获取单个项目
func C_GetProjectOne(ctx *context.New) {

  // get
  id, _ := ctx.GetParam("id")

  // query
  project, err := service.GetProjectInfoById(ctx, id)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(gin.H{
    "data": project.GetMap(),
  })
}

// 删除单个项目
func C_DelProjectOne(ctx *context.New) {

  // get
  id, _ := ctx.GetParam("id")

  // query
  err := service.UpdateProject(ctx, id, bson.M{
    "exist": false,
  })

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(nil)
}

// 获取项目列表
func C_GetProjectsList(ctx *context.New) {

  // get
  skip, _ := ctx.GetQueryInt("skip")
  limit, _ := ctx.GetQueryInt("limit")

  // query
  data, err := service.GetProjectsList(ctx, skip, limit, bson.M{})

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(gin.H{
    "data": data.Each(func(item *model.Project) gin.H {
      each := item.GetMap()
      if len(item.Missions) > 0 {
        each["progress"] = service.GetProjectProgress(ctx, item.Id)
      }
      return each
    }),
  })
}

// 创建项目
func C_CreateProject(ctx *context.New) {

  // get
  name, _ := ctx.GetRaw("name")
  deadline, _ := ctx.GetRawTime("deadline")
  departments, _ := ctx.GetRawArray("departments")
  description, _ := ctx.GetRaw("description")
  weight, _ := ctx.GetRawInt("weight")

  // create
  data := model.Project{
    Name:        name,
    Deadline:    deadline,
    Description: description,
    Weight:      weight,
    Status:      1,
  }

  // insert
  err := service.CreateProject(ctx, data, departments)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(nil)
}

// 修改项目
func C_UpdateProject(ctx *context.New) {

  // get
  id, _ := ctx.GetParam("id")

  // update
  data := bson.M{}

  if name, ok := ctx.GetRaw("name"); ok {
    data["name"] = name
  }

  if deadline, ok := ctx.GetRawTime("deadline"); ok {
    data["deadline"] = deadline
  }

  if departments, ok := ctx.GetRawArray("departments"); ok {
    data["departments"] = departments
  }

  if description, ok := ctx.GetRaw("description"); ok {
    data["description"] = description
  }

  if weight, ok := ctx.GetRawInt("weight"); ok {
    data["weight"] = weight
  }

  err := service.UpdateProject(ctx, id, data)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  ctx.Success(nil)
}
