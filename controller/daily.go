package controller

import (
  "errors"
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2/bson"
  "workerbook/model"
  "workerbook/service"
)

// 获取日报列表
func GetDailiesList(c *gin.Context) {
  ctx := CreateCtx(c)

  skip := ctx.getQueryInt("skip")
  limit := ctx.getQueryInt("limit")

  dailiesList, err := service.GetDailiesList(skip, limit)

  if err != nil {
    ctx.Error(err, 1)
    return
  }

  ctx.Success(gin.H{
    "list": dailiesList,
  })
}

// 获取单个日报的信息
func GetDailyOne(c *gin.Context) {
  ctx := CreateCtx(c)

  id := ctx.getParam("id")

  if !bson.IsObjectIdHex(id) {
    ctx.Error(errors.New("无效的用户ID"), 1)
    return
  }

  dailyInfo, err := service.GetDailyInfoById(bson.ObjectIdHex(id))

  if err != nil {
    ctx.Error(err, 1)
    return
  }

  ctx.Success(gin.H{
    "data": dailyInfo,
  })
}

// create my daily item at today.
func CreateMyTodayDailyItem(c *gin.Context) {
  ctx := CreateCtx(c)

  id := ctx.get("uid")

  if !bson.IsObjectIdHex(id) {
    ctx.Error(errors.New("无效的用户ID"), 1)
    return
  }

  pid := ctx.getRaw("project")
  progress := ctx.getRawInt("progress")
  record := ctx.getRaw("record")
  projectName := ""

  if pid != "" {
    if !bson.IsObjectIdHex(pid) {
      ctx.Error(errors.New("无效的项目"), 1)
      return
    }

    // find the project info.
    project, err := service.GetProjectInfoById(bson.ObjectIdHex(pid))

    if err != nil {
      ctx.Error(errors.New("找不到相关项目"), 1)
      return
    }

    // set project name
    projectName = project.Name
  }

  // 找到用户今天的日报内容(找不到会创建一个空内容的日报数据)
  dailyInfo, err := service.GetUserTodayDaily(bson.ObjectIdHex(id))

  if err != nil {
    ctx.Error(err, 1)
    return
  }

  if err != nil {
    ctx.Error(errors.New("请设置正确的项目进度"), 1)
    return
  }

  // create record data.
  data := model.DailyItem{
    Id:       bson.NewObjectId(),
    Record:   record,
    Progress: progress,
    Pname:    projectName,
    Pid:      pid,
  }

  // insert it.
  err = service.AppendDailyItemIntoUsersDailyList(data, dailyInfo.Id)

  if err != nil {
    ctx.Error(err, 1)
    return
  }

  ctx.Success(nil)
}

// 删除今天的日报
func DeleteUserTodayDailyItem(c *gin.Context) {
  ctx := CreateCtx(c)

  id := ctx.get("uid")

  if !bson.IsObjectIdHex(id) {
    ctx.Error(errors.New("无效的用户ID"), 1)
    return
  }

  itemId := ctx.getParam("itemId")

  // 删除今天中相应的日报
  err := service.DeleteTodayDailyItemFromUsersDailyList(bson.ObjectIdHex(id), bson.ObjectIdHex(itemId))

  if err != nil {
    ctx.Error(err, 1)
    return
  }
  ctx.Success(gin.H{})
}

// 获取我今天的日报
func GetMyTodayDaily(c *gin.Context) {
  ctx := CreateCtx(c)

  id := ctx.get("uid")

  if !bson.IsObjectIdHex(id) {
    ctx.Error(errors.New("无效的用户ID"), 1)
    return
  }

  dailyInfo, err := service.GetUserTodayDaily(bson.ObjectIdHex(id))

  if err != nil {
    ctx.Error(err, 1)
    return
  }

  ctx.Success(gin.H{
    "data": dailyInfo.DailyList,
  })
}
