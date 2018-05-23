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
  defer ctx.handleErrorIfPanic()

  skip := ctx.getQueryInt("skip")
  limit := ctx.getQueryInt("limit")
  status := ctx.getQueryInt("status")

  // check
  ctx.PanicIfIntLessThen(skip, 0, "Skip不能小于0")
  ctx.PanicIfIntLessThen(limit, 0, "Limit不能小于0")
  ctx.PanicIfIntMoreThen(limit, 100, "Limit不能大于100")

  // create search sql
  search := bson.M{}

  if status != 0 {
    search["status"] = status
  }

  projectsList, err := service.GetProjectsList(skip, limit, search)

  if err != nil {
    panic("获取项目列表失败")
  }

  ctx.Success(gin.H{
    "list": projectsList,
  })
}

// create project.
func CreateProject(c *gin.Context) {
  ctx := CreateCtx(c)
  defer ctx.handleErrorIfPanic()

  name := ctx.getRaw("name")

  // check
  ctx.PanicIfStringIsEmpty(name, "项目名不能为空")

  data := model.Project{
    Name: name,
  }

  err := service.CreateProject(data)

  if err != nil {
    panic("创建项目失败")
  }

  ctx.Success(nil)
}
