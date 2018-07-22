package controller

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "workerbook/model"
  "workerbook/service"
  "workerbook/util"
)

// 获取单个任务
func GetMissionOne(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  id, _ := ctx.getParam("id")

  // query
  missionInfo, err := service.GetMissionInfoById(id)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // query project
  if missionInfo["projectId"] != "" {
    projectId := missionInfo["projectId"].(bson.ObjectId)
    projectInfo, err := service.GetProjectInfoById(projectId.Hex())

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
  ctx := CreateCtx(c)

  // get
  name, _ := ctx.getRaw("name")
  deadline, _ := ctx.getRawTime("deadline")
  userId, _ := ctx.getRaw("userId")
  projectId, _ := ctx.getRaw("projectId")

  // create
  data := model.Mission{
    Name:        name,
    Deadline:    deadline,
    Progress:    0,
    Status:      1,
  }

  // insert
  err := service.CreateMission(data, projectId, userId)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  ctx.Success(nil)
}

// 更新任务
func UpdateMission(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  id, _ := ctx.getParam("id")

  // update
  data := bson.M{}

  if name, ok := ctx.getRaw("name"); ok {
    data["name"] = name
  }

  if description, ok := ctx.getRaw("description"); ok {
    data["description"] = description
  }

  if deadline, ok := ctx.getRawTime("deadline"); ok {
    data["deadline"] = deadline
  }

  if userId, ok := ctx.getRaw("userId"); ok {
    data["userId"] = userId
  }

  if projectId, ok := ctx.getRaw("projectId"); ok {
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