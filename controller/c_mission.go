package controller

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "workerbook/model"
  "workerbook/service"
)

// 获取单个任务
func C_GetMissionOne(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  id := ctx.getParam("id")

  // query
  missionInfo, err := service.GetMissionInfoById(id)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(gin.H{
    "data": missionInfo,
  })
}

// 获取任务列表
func C_GetMissionsList(c *gin.Context) {
  ctx := CreateCtx(c)

  // return
  ctx.Success(nil)
}

// 创建任务
func C_CreateMission(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  name := ctx.getRaw("name")
  description := ctx.getRaw("description")
  deadline := ctx.getRawTime("deadline")
  projectId := ctx.getRaw("projectId")

  // create
  data := model.Mission{
    Name:        name,
    Description: description,
    Deadline:    deadline,
    Progress:    0,
    Status:      1,
  }

  // insert
  err := service.CreateMission(data, projectId)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  ctx.Success(nil)
}

// 更新任务
func C_UpdateMission(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  id := ctx.getParam("id")

  // update
  data := bson.M{}

  if name := ctx.getRaw("name"); name != "" {
    data["name"] = name
  }

  if description := ctx.getRaw("description"); description != "" {
    data["description"] = description
  }

  if deadline := ctx.getRawTime("deadline"); !deadline.IsZero() {
    data["deadline"] = deadline
  }

  if projectId := ctx.getRaw("projectId"); projectId != "" {
    data["projectId"] = projectId
  }

  err := service.UpdateMission(id, data)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  ctx.Success(nil)
}
