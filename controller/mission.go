package controller

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "workerbook/conf"
  "workerbook/context"
  "workerbook/model"
  "workerbook/service"
)

// 获取分配到自己的任务列表
func GetOwnsMissionsList(ctx *context.New) {

  // get
  skip, _ := ctx.GetQueryInt("skip")
  limit, _ := ctx.GetQueryInt("limit")
  ownUserId, _ := ctx.Get(conf.OWN_USER_ID)

  // query
  data, err := service.GetMissionsList(ctx, skip, limit, bson.M{
    "user.$id": bson.ObjectIdHex(ownUserId),
  })

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(gin.H{
    "data": data.Each(func(item model.Mission) gin.H {
      each := item.GetMap("createTime", "user", "preProgress", "chartTime", "editor", "editTime", "exist")
      project, err := service.FindProjectRef(ctx, &item.Project)
      if err == nil {
        each["project"] = project.GetMap("departments", "missions", "editor", "editTime", "exist")
      }
      return each
    }),
  })
}

// 获取单个任务
func GetMissionOne(ctx *context.New) {

  // get
  id, _ := ctx.GetParam("id")

  // query
  mission, err := service.GetMissionInfoById(ctx, id)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // query project
  project, err := service.GetProjectInfoById(ctx, mission.Project.Id.(bson.ObjectId).Hex())

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // query user
  user, err := service.FindUserRef(ctx, &mission.User)

  // return
  data := mission.GetMap("preProgress", "chartTime", "editor", "editTime", "exist")
  data["project"] = project.GetMap("departments", "missions", "editor", "editTime", "exist")
  data["user"] = user.GetMap("username", "department", "editor", "editTime")

  ctx.Success(gin.H{
    "data": data,
  })
}

// 创建任务
func CreateMission(ctx *context.New) {

  // get
  name, _ := ctx.GetRaw("name")
  deadline, _ := ctx.GetRawTime("deadline")
  userId, _ := ctx.GetRaw("userId")
  projectId, _ := ctx.GetRaw("projectId")

  // create
  data := model.Mission{
    Name:     name,
    Deadline: deadline,
  }

  // insert
  err := service.CreateMission(ctx, data, projectId, userId)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  ctx.Success(nil)
}

// 更新任务
func UpdateMission(ctx *context.New) {

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

  err := service.UpdateMission(ctx, id, data)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  ctx.Success(nil)
}
