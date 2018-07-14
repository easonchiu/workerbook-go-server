package controller

import (
  "github.com/gin-gonic/gin"
  "time"
  "workerbook/model"
  "workerbook/service"
)


// 获取单个任务
func GetMissionOne(c *gin.Context) {
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
func GetMissionsList(c *gin.Context) {
  ctx := CreateCtx(c)

  // return
  ctx.Success(nil)
}

// 创建任务
func CreateMission(c *gin.Context) {
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
    CreateTime:  time.Now(),
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
func UpdateMission(c *gin.Context) {
  ctx := CreateCtx(c)

  // get
  id := ctx.getParam("id")
  name := ctx.getRaw("name")
  description := ctx.getRaw("description")
  deadline := ctx.getRawTime("deadline")
  projectId := ctx.getRaw("projectId")

  // update
  data := model.Mission{
    Name:        name,
    Description: description,
    Deadline:    deadline,
    CreateTime:  time.Now(),
    Progress:    0,
    Status:      1,
  }

  err := service.UpdateMission(id, data, projectId)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  ctx.Success(nil)
}
