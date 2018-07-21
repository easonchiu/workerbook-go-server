package controller

import (
  "fmt"
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "workerbook/model"
  "workerbook/service"
)

// 获取单个项目
func C_GetProjectOne(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  id, _ := ctx.getParam("id")

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
  id, _ := ctx.getParam("id")

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
  skip, _ := ctx.getQueryInt("skip")
  limit, _ := ctx.getQueryInt("limit")

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
  name, _ := ctx.getRaw("name")
  deadline, _ := ctx.getRawTime("deadline")
  departments, _ := ctx.getRawArray("departments")
  description, _ := ctx.getRaw("description")
  weight, _ := ctx.getRawInt("weight")

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
  id, _ := ctx.getParam("id")

  // update
  data := bson.M{}

  fmt.Println(string(ctx.RawData), "< 123123333")

  if name, ok := ctx.getRaw("name"); ok {
    data["name"] = name
  }

  if deadline, ok := ctx.getRawTime("deadline"); ok {
    data["deadline"] = deadline
  }

  if departments, ok := ctx.getRawArray("departments"); ok {
    data["departments"] = departments
  }

  if description, ok := ctx.getRaw("description"); ok {
    data["description"] = description
  }

  if weight, ok := ctx.getRawInt("weight"); ok {
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
