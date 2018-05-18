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

  skip := ctx.getQuery("skip", true).(int)
  limit := ctx.getQuery("limit", true).(int)

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
func GetDailyInfo(c *gin.Context) {
  ctx := CreateCtx(c)

  id := ctx.getParam("id").(string)

  if !bson.IsObjectIdHex(id) {
    ctx.Error(errors.New("无效的id号"), 1)
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

// create daily item at today.
func CreateTodayDailyItem(c *gin.Context) {
  ctx := CreateCtx(c)

  pid := ctx.getRaw("project").(string)
  progress := ctx.getRaw("progress").(int)
  record := ctx.getRaw("record").(string)
  uid, _ := c.Get("uid")

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

  // 找到用户今天的日报内容(找不到会创建一个空内容的日报数据)
  dailyInfo, err := service.GetUserTodayDailyByUid(bson.ObjectIdHex(uid.(string)))

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
    Pname:    project.Name,
    Pid:      project.Id.Hex(),
  }

  // insert it.
  err = service.AppendDailyItemIntoUsersDailyList(data, dailyInfo.Id)

  if err != nil {
    ctx.Error(err, 1)
    return
  }

  ctx.Success(gin.H{
    "data": dailyInfo,
  })
}

// 删除今天的日报
func DeleteTodayDailyItem(c *gin.Context) {
  ctx := CreateCtx(c)

  uid := ctx.getRaw("uid").(string)
  itemId := ctx.getParam("itemId").(string)

  if !bson.IsObjectIdHex(uid) {
    ctx.Error(errors.New("无效的id号"), 1)
    return
  }

  // 删除今天中相应的日报
  err := service.DeleteTodayDailyItemFromUsersDailyList(bson.ObjectIdHex(uid), bson.ObjectIdHex(itemId))

  if err != nil {
    ctx.Error(err, 1)
    return
  }
  ctx.Success(gin.H{})
}

// 获取今天的日报
func GetTodayDaily(c *gin.Context) {
  ctx := CreateCtx(c)

  uid := ctx.getParam("id").(string)

  dailyInfo, err := service.GetUserTodayDailyByUid(bson.ObjectIdHex(uid))

  if err != nil {
    ctx.Error(err, 1)
    return
  }

  ctx.Success(gin.H{
    "data": dailyInfo,
  })
}
