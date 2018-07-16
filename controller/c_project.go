package controller

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "workerbook/model"
  "workerbook/service"
)

// 获取单个项目
func C_GetProjectOne(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  id := ctx.getParam("id")

  // query
  projectInfo, err := service.GetProjectInfoById(id)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(gin.H{
    "data": projectInfo,
  })
}

// 删除单个项目
func C_DelProjectOne(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  id := ctx.getParam("id")

  // query
  err := service.DelProjectById(id)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(nil)
}

// 获取项目列表
func C_GetProjectsList(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  skip := ctx.getQueryInt("skip")
  limit := ctx.getQueryInt("limit")

  // query
  data, err := service.GetProjectsList(skip, limit, bson.M{})

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(gin.H{
    "data": data,
  })
}

// 创建项目
func C_CreateProject(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  name := ctx.getRaw("name")
  deadline := ctx.getRawTime("deadline")
  departments := ctx.getRawArray("departments")
  description := ctx.getRaw("description")
  weight := ctx.getRawInt("weight")

  // create
  data := model.Project{
    Name:        name,
    Deadline:    deadline,
    Description: description,
    Weight:      weight,
    Status:      1,
  }

  // insert
  err := service.CreateProject(data, departments)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(nil)
}

// 修改项目
func C_UpdateProject(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  id := ctx.getParam("id")

  // update
  data := bson.M{}

  if name := ctx.getRaw("name"); name != "" {
    data["name"] = name
  }

  if deadline := ctx.getRawTime("deadline"); !deadline.IsZero() {
    data["deadline"] = deadline
  }

  if departments := ctx.getRawArray("departments"); len(departments) != 0 {
    data["departments"] = departments
  }

  if description := ctx.getRaw("description"); description != "" {
    data["description"] = description
  }

  if weight := ctx.getRawInt("weight"); weight != 0 {
    data["weight"] = weight
  }

  err := service.UpdateProject(id, data)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  ctx.Success(nil)
}
