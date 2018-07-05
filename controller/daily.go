package controller

import (
  `github.com/gin-gonic/gin`
  `gopkg.in/mgo.v2`
  `gopkg.in/mgo.v2/bson`
  "workerbook/errno"
  `workerbook/model`
  `workerbook/service`
)

// 获取日报列表
func GetDailiesList(c *gin.Context) {
  ctx := CreateCtx(c)
  defer ctx.HandleError()

  // get
  skip := ctx.getQueryInt("skip")
  limit := ctx.getQueryInt("limit")

  // check
  ctx.ErrorIfIntLessThen(skip, 0, errno.ErrorSkipRange)
  ctx.ErrorIfIntLessThen(limit, 0, errno.ErrorSkipRange)
  ctx.ErrorIfIntMoreThen(limit, 100, errno.ErrorSkipRange)

  // query
  dailiesList, err := service.GetDailiesList(skip, limit)

  // check
  if err != nil {
    ctx.Error(err)
  }

  // return
  ctx.Success(gin.H{
    "list": dailiesList,
  })
}

// 获取单个日报的信息
func GetDailyOne(c *gin.Context) {
  ctx := CreateCtx(c)
  defer ctx.HandleError()

  // get
  id := ctx.getParam("id")

  // check
  ctx.ErrorIfStringNotObjectId(id, errno.ErrorDailyIdError)

  // query
  dailyInfo, err := service.GetDailyInfoById(bson.ObjectIdHex(id))

  // check
  if err != nil {
    ctx.Error(err)
  }

  // return
  ctx.Success(gin.H{
    "data": dailyInfo,
  })
}

// create my daily item at today.
func CreateMyTodayDailyItem(c *gin.Context) {
  ctx := CreateCtx(c)
  defer ctx.HandleError()

  uid := ctx.get("uid")
  pid := ctx.getRaw("project")
  progress := ctx.getRawInt("progress")
  record := ctx.getRaw("record")

  // check
  ctx.ErrorIfStringNotObjectId(uid, errno.ErrorUserIdError)
  ctx.ErrorIfStringIsEmpty(record, "日报内容不能为空")
  ctx.ErrorIfIntLessThen(progress, 0, "请设置正确的项目进度")
  ctx.ErrorIfIntMoreThen(progress, 100, "请设置正确的项目进度")

  // check user is exist or not.
  _, err := service.GetUserInfoById(bson.ObjectIdHex(uid))

  if err != nil {
    ctx.Error(err)
  }

  // find the project info if has pid
  if pid != "" {
    ctx.ErrorIfStringNotObjectId(pid, "无效的项目")

    // find the project info.
    _, err := service.GetProjectInfoById(bson.ObjectIdHex(pid))

    if err != nil && err != mgo.ErrNotFound {
      panic(err)
    }

    // set project name
    // projectName = project.Name
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
    Progress: progress,
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
  defer ctx.HandleError()

  uid := ctx.get("uid")
  itemId := ctx.getParam("itemId")

  // check
  ctx.ErrorIfStringNotObjectId(uid, errno.ErrorUserIdError)
  ctx.ErrorIfStringNotObjectId(itemId, errno.ErrorDailyIdError)

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
  defer ctx.HandleError()

  uid := ctx.get("uid")

  ctx.ErrorIfStringNotObjectId(uid, errno.ErrorUserIdError)

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
