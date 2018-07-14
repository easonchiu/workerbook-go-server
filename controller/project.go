package controller

import (
  "github.com/gin-gonic/gin"
  "time"
  "workerbook/model"
  "workerbook/service"
)

// 获取单个项目
func GetProjectOne(c *gin.Context) {
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
func DelProjectOne(c *gin.Context) {
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
func GetProjectsList(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  skip := ctx.getQueryInt("skip")
  limit := ctx.getQueryInt("limit")

  // query
  data, err := service.GetProjectsList(skip, limit, model.Project{})

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
func CreateProject(c *gin.Context) {
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
    CreateTime:  time.Now(),
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
func UpdateProject(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  id := ctx.getParam("id")
  name := ctx.getRaw("name")
  deadline := ctx.getRawTime("deadline")
  departments := ctx.getRawArray("departments")
  description := ctx.getRaw("description")
  weight := ctx.getRawInt("weight")

  // update
  data := model.Project{
    Name:        name,
    Deadline:    deadline,
    Description: description,
    Weight:      weight,
  }

  err := service.UpdateProject(id, data, departments)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  ctx.Success(nil)
}
