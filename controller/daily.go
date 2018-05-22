package controller

import (
  `github.com/gin-gonic/gin`
  `gopkg.in/mgo.v2`
  `gopkg.in/mgo.v2/bson`
  `workerbook/model`
  `workerbook/service`
)

// 获取日报列表
func GetDailiesList(c *gin.Context) {
  ctx := CreateCtx(c)

  skip := ctx.getQueryInt("skip")
  limit := ctx.getQueryInt("limit")

  if !ctx.CheckIntIsLessThen(skip, 0, "Skip不能小于0") {
    return
  }

  if !ctx.CheckIntIsLessThen(limit, 0, "Limit不能小于0") {
    return
  }

  if !ctx.CheckIntIsMoreThen(limit, 0, "Limit不能小于100") {
    return
  }

  dailiesList, err := service.GetDailiesList(skip, limit)

  if err != nil {
    ctx.Error("找不到相关数据", 1)
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

  ctx.CheckIsObjectIdHex(id, "无效的用户ID")

  dailyInfo, err := service.GetDailyInfoById(bson.ObjectIdHex(id))

  if err != nil {
    ctx.Error("找不到相关数据", 1)
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

  // check user is exist.
  _, err := service.GetUserInfoById(bson.ObjectIdHex(id))

  if err != nil {
    ctx.Error("没有相关的用户", 1)
    return
  }

  // check args.
  pid := ctx.getRaw("project")
  progress := ctx.getRawInt("progress")
  record := ctx.getRaw("record")
  projectName := ""

  if record == "" {
    ctx.Error("日报内容不能为空", 1)
    return
  }

  if pid != "" {
    if !bson.IsObjectIdHex(pid) {
      ctx.Error("无效的项目", 1)
      return
    }

    // find the project info.
    project, err := service.GetProjectInfoById(bson.ObjectIdHex(pid))

    if err != nil && err != mgo.ErrNotFound {
      ctx.Error(err.Error(), 1)
      return
    }

    // set project name
    projectName = project.Name
  }

  // 找到用户今天的日报内容
  dailyInfo, err := service.GetUserTodayDaily(bson.ObjectIdHex(id))

  // if error and typeof error is not not-found.
  if err != nil && err != mgo.ErrNotFound {
    ctx.Error(err.Error(), 1)
    return
  }

  // if not-found, create empty today daily
  if err == mgo.ErrNotFound {
    newDaily, err := service.CreateMyTodayDaily(bson.ObjectIdHex(id))
    if err != nil {
      ctx.Error("创建日报失败", 1)
      return
    }
    dailyInfo = newDaily
  }

  // progress must be between 0 to 100
  if progress < 0 || progress > 100 {
    ctx.Error("请设置正确的项目进度", 1)
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
    ctx.Error("创建日报失败", 1)
    return
  }

  ctx.Success(nil)
}

// 删除今天的日报
func DeleteUserTodayDailyItem(c *gin.Context) {
  ctx := CreateCtx(c)

  id := ctx.get("uid")

  if !bson.IsObjectIdHex(id) {
    ctx.Error("无效的用户ID", 1)
    return
  }

  itemId := ctx.getParam("itemId")

  // 删除今天中相应的日报
  err := service.DeleteTodayDailyItemFromUsersDailyList(bson.ObjectIdHex(id), bson.ObjectIdHex(itemId))

  if err != nil {
    ctx.Error("删除日报失败", 1)
    return
  }
  ctx.Success(gin.H{})
}

// 获取我今天的日报
func GetMyTodayDaily(c *gin.Context) {
  ctx := CreateCtx(c)

  id := ctx.get("uid")

  dailyInfo, err := service.GetUserTodayDaily(bson.ObjectIdHex(id))

  if err != nil {
    // if not found, return empty data.
    if err == mgo.ErrNotFound {
      ctx.Success(gin.H{
        "data": []model.DailyItem{},
      })
    } else {
      ctx.Error(err.Error(), 1)
    }
    return
  }

  ctx.Success(gin.H{
    "data": dailyInfo.DailyList,
  })
}
