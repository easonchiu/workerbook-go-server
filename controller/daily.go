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
  defer ctx.handleErrorIfPanic()

  skip := ctx.getQueryInt("skip")
  limit := ctx.getQueryInt("limit")

  // check
  ctx.PanicIfIntLessThen(skip, 0, "Skip不能小于0")
  ctx.PanicIfIntLessThen(limit, 0, "Limit不能小于0")
  ctx.PanicIfIntMoreThen(limit, 100, "Limit不能大于100")

  // check pass

  dailiesList, err := service.GetDailiesList(skip, limit)

  if err != nil {
    panic("找不到相关数据")
  }

  ctx.Success(gin.H{
    "list": dailiesList,
  })
}

// 获取单个日报的信息
func GetDailyOne(c *gin.Context) {
  ctx := CreateCtx(c)
  defer ctx.handleErrorIfPanic()

  id := ctx.getParam("id")

  // check
  ctx.PanicIfStringNotObjectId(id, "无效的用户ID")

  // check pass

  dailyInfo, err := service.GetDailyInfoById(bson.ObjectIdHex(id))

  if err != nil {
    panic("找不到相关数据")
  }

  ctx.Success(gin.H{
    "data": dailyInfo,
  })
}

// create my daily item at today.
func CreateMyTodayDailyItem(c *gin.Context) {
  ctx := CreateCtx(c)
  defer ctx.handleErrorIfPanic()

  uid := ctx.get("uid")
  pid := ctx.getRaw("project")
  progress := ctx.getRawInt("progress")
  record := ctx.getRaw("record")
  projectName := ""

  // check
  ctx.PanicIfStringNotObjectId(uid, "无效的用户ID")
  ctx.PanicIfStringIsEmpty(record, "日报内容不能为空")
  ctx.PanicIfIntLessThen(progress, 0, "请设置正确的项目进度")
  ctx.PanicIfIntMoreThen(progress, 100, "请设置正确的项目进度")

  // check user is exist.
  _, err := service.GetUserInfoById(bson.ObjectIdHex(uid))

  if err != nil {
    panic("没有相关的用户")
  }

  if pid != "" {
    ctx.PanicIfStringNotObjectId(pid, "无效的项目")

    // find the project info.
    project, err := service.GetProjectInfoById(bson.ObjectIdHex(pid))

    if err != nil && err != mgo.ErrNotFound {
      panic(err)
    }

    // set project name
    projectName = project.Name
  }

  // 找到用户今天的日报内容
  dailyInfo, err := service.GetUserTodayDaily(bson.ObjectIdHex(uid))

  // if error and typeof error is not not-found.
  if err != nil && err != mgo.ErrNotFound {
    panic("创建日报失败")
  }

  // if not-found, create empty today daily
  if err == mgo.ErrNotFound {
    newDaily, err := service.CreateMyTodayDaily(bson.ObjectIdHex(uid))
    if err != nil {
      panic("创建日报失败")
    }
    dailyInfo = newDaily
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
    panic("创建日报失败")
  }

  ctx.Success(nil)
}

// 删除今天的日报
func DeleteUserTodayDailyItem(c *gin.Context) {
  ctx := CreateCtx(c)
  defer ctx.handleErrorIfPanic()

  uid := ctx.get("uid")
  itemId := ctx.getParam("itemId")

  // check
  ctx.PanicIfStringNotObjectId(uid, "无效的用户ID")
  ctx.PanicIfStringNotObjectId(itemId, "无效的日报ID")

  // 删除今天中相应的日报
  err := service.DeleteTodayDailyItemFromUsersDailyList(bson.ObjectIdHex(uid), bson.ObjectIdHex(itemId))

  if err != nil {
    panic("删除日报失败")
  }

  ctx.Success(nil)
}

// 获取我今天的日报
func GetMyTodayDaily(c *gin.Context) {
  ctx := CreateCtx(c)
  defer ctx.handleErrorIfPanic()

  uid := ctx.get("uid")

  ctx.PanicIfStringNotObjectId(uid, "无效的用户ID")

  dailyInfo, err := service.GetUserTodayDaily(bson.ObjectIdHex(uid))

  if err != nil {
    // if not found, return empty data.
    if err == mgo.ErrNotFound {
      ctx.Success(gin.H{
        "data": []model.DailyItem{},
      })
    } else {
      panic("获取日报失败")
    }
  } else {
    ctx.Success(gin.H{
      "data": dailyInfo.DailyList,
    })
  }
}
