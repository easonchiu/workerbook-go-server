package controller

import (
  `github.com/gin-gonic/gin`
  `gopkg.in/mgo.v2/bson`
  `workerbook/model`
  `workerbook/service`
)

// query projects list
func GetProjectsList(c *gin.Context) {
  ctx := CreateCtx(c)

  skip := ctx.getQueryInt("skip")
  limit := ctx.getQueryInt("limit")
  status := ctx.getQueryInt("status")

  // create search sql
  search := bson.M{}

  if status != 0 {
    search["status"] = status
  }

  projectsList, err := service.GetProjectsList(skip, limit, search)

  if err != nil {
    ctx.Error(err.Error(), 1)
    return
  }

  ctx.Success(gin.H{
    "list": projectsList,
  })
}

// create project.
func CreateProject(c *gin.Context) {
  ctx := CreateCtx(c)

  data := model.Project{
    Name: ctx.getRaw("name"),
  }

  err := service.CreateProject(data)

  if err != nil {
    ctx.Error("创建项目失败", 1)
    return
  }

  ctx.Success(nil)
}
