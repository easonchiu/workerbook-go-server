package controller

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "workerbook/context"
  "workerbook/model"
  "workerbook/service"
  "workerbook/util"
)

// 获取单个任务
func GetMissionOne(c *gin.Context) {
  ctx, err := context.CreateCtx(c)
  defer ctx.Close()

  if err != nil {
    ctx.Error(err)
    return
  }

  // get
  id, _ := ctx.GetParam("id")

  // query
  missionInfo, err := service.GetMissionInfoById(ctx, id, "user")

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // query project
  if missionInfo["projectId"] != "" {
    projectId := missionInfo["projectId"].(bson.ObjectId)
    projectInfo, err := service.GetProjectInfoById(ctx, projectId.Hex())

    if err != nil {
      ctx.Error(err)
      return
    }

    util.Forget(projectInfo, "departments missions description")

    missionInfo["project"] = projectInfo
  }

  util.Forget(missionInfo, "projectId")

  // return
  ctx.Success(gin.H{
    "data": missionInfo,
  })
}

// 创建任务
func CreateMission(c *gin.Context) {
  ctx, err := context.CreateCtx(c)
  defer ctx.Close()

  if err != nil {
    ctx.Error(err)
    return
  }

  // get
  name, _ := ctx.GetRaw("name")
  deadline, _ := ctx.GetRawTime("deadline")
  userId, _ := ctx.GetRaw("userId")
  projectId, _ := ctx.GetRaw("projectId")

  // create
  data := model.Mission{
    Name:     name,
    Deadline: deadline,
    Progress: 0,
    Status:   1,
  }

  // insert
  err = service.CreateMission(ctx, data, projectId, userId)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  ctx.Success(nil)
}

// 更新任务
func UpdateMission(c *gin.Context) {
  ctx, err := context.CreateCtx(c)
  defer ctx.Close()

  if err != nil {
    ctx.Error(err)
    return
  }

  // get
  id, _ := ctx.GetParam("id")

  // update
  data := bson.M{}

  if name, ok := ctx.GetRaw("name"); ok {
    data["name"] = name
  }

  if description, ok := ctx.GetRaw("description"); ok {
    data["description"] = description
  }

  if deadline, ok := ctx.GetRawTime("deadline"); ok {
    data["deadline"] = deadline
  }

  if userId, ok := ctx.GetRaw("userId"); ok {
    data["userId"] = userId
  }

  if projectId, ok := ctx.GetRaw("projectId"); ok {
    data["projectId"] = projectId
  }

  err = service.UpdateMission(ctx, id, data)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  ctx.Success(nil)
}
