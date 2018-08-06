package controller

import (
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "time"
  "workerbook/conf"
  "workerbook/context"
  "workerbook/model"
  "workerbook/service"
)

// 创建日报
func CreateDaily(ctx *context.New) {

  // get
  content, _ := ctx.GetRaw("content")
  progress, _ := ctx.GetRawInt("progress")
  missionId, _ := ctx.GetRaw("missionId")

  // create
  data := model.DailyItem{
    Content:  content,
    Progress: progress,
  }

  // insert
  err := service.CreateDaily(ctx, data, missionId)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(nil)
}

// 修改日报内容
func UpdateDaily(ctx *context.New) {

  // get
  content, _ := ctx.GetRaw("content")
  id, _ := ctx.GetRaw("id")

  // update
  err := service.UpdateDailyContent(ctx, id, content)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(nil)
}

// 删除当天写的某一条日报内容
func DelDaily(ctx *context.New) {
  id, _ := ctx.GetRaw("id")

  // delete
  err := service.DelDailyContent(ctx, id)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(nil)
}

// 获取我今天的日报数据
func GetTodayDaily(ctx *context.New) {

  // get
  ownUserId, _ := ctx.Get(conf.OWN_USER_ID)

  daily, err := service.GetDailyByDay(ctx, ownUserId, time.Now().Format("2006-01-02"))

  // 该接口不管有没有数据都不能返回报错
  if err != nil {
    ctx.Success(gin.H{
      "data": []gin.H{},
    })
    return
  }

  ctx.Success(gin.H{
    "data": daily.GetMap()["dailies"],
  })
}

// 根据天获取日报列表
func GetDailiesListByDay(ctx *context.New) {

  // get
  skip, _ := ctx.GetQueryInt("skip")
  limit, _ := ctx.GetQueryInt("limit")

  data, err := service.GetDailiesList(ctx, skip, limit, bson.M{
    "day": time.Now().Format("2006-01-02"),
  })

  // 该接口不管有没有数据都不能返回报错
  if err != nil || data.Count == 0 {
    ctx.Success(gin.H{
      "data": gin.H{
        "count": 0,
        "limit": 0,
        "skip":  0,
        "list":  []gin.H{},
      },
    })
    return
  }

  ctx.Success(gin.H{
    "data": data.Each(func(item *model.Daily) gin.H {
      each := item.GetMap()
      user, err := service.FindUserRef(ctx, &item.User)
      if err == nil {
        each["user"] = user.GetMap("createTime", "department", "email", "exist", "mobile")
      }
      return each
    }),
  })
}

// 独立更新今天日报的任务进度
func UpdateDailyMissionProgress(ctx *context.New) {
  // get
  progress, _ := ctx.GetRawInt("progress")
  missionId, _ := ctx.GetParam("id")

  err := service.UpdateDailyMissionProgress(ctx, missionId, progress)

  // check
  if err != nil {
    ctx.Error(err)
    return
  }

  // return
  ctx.Success(nil)
}