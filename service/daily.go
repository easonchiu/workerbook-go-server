package service

import (
  "errors"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "time"
  "workerbook/cache"
  "workerbook/conf"
  "workerbook/context"
  "workerbook/errgo"
  "workerbook/model"
)

// 创建日报
func CreateDaily(ctx *context.New, data model.DailyItem, missionId string) error {

  // get
  ownUserId, _ := ctx.Get(conf.OWN_USER_ID)
  ownDepartmentId, _ := ctx.Get(conf.OWN_DEPARTMENT_ID)

  ctx.Errgo.ErrorIfStringIsEmpty(data.Content, errgo.ErrDailyContentEmpty)
  ctx.Errgo.ErrorIfIntLessThen(data.Progress, 0, errgo.ErrDailyProgressRange)
  ctx.Errgo.ErrorIfIntMoreThen(data.Progress, 100, errgo.ErrDailyProgressRange)
  ctx.Errgo.ErrorIfStringNotObjectId(missionId, errgo.ErrDailyIdError)

  if err := ctx.Errgo.PopError(); err != nil {
    return err
  }

  // 找到任务的信息
  mission, err := GetMissionInfoById(ctx, missionId)

  if err != nil {
    return errors.New(errgo.ErrCreateDailyFailed)
  }

  // 查看这任务是不是自己的
  if mission.User.Id.(bson.ObjectId) != bson.ObjectIdHex(ownUserId) {
    return errors.New(errgo.ErrCreateDailyFailed)
  }

  // 找到项目信息
  project, err := FindProjectRef(ctx, &mission.Project)

  if err != nil {
    return errors.New(errgo.ErrCreateDailyFailed)
  }

  // 补全dailyItem
  data.Id = bson.NewObjectId()
  data.MissionId = mission.Id
  data.MissionName = mission.Name
  data.ProjectId = project.Id
  data.ProjectName = project.Name

  // 先查找该用户今天是否有写日报
  count, err := ctx.MgoDB.C(model.DailyCollection).Find(bson.M{
    "user.$id": bson.ObjectIdHex(ownUserId),
    "day":      time.Now().Format("2006-01-02"),
  }).Count()

  // 没有写的话，插入一条完整日报数据
  if count == 0 {

    // 找到部门信息
    department, err := GetDepartmentInfoById(ctx, ownDepartmentId)

    if err != nil {
      return errors.New(errgo.ErrDepartmentNotFound)
    }

    daily := model.Daily{
      User: mgo.DBRef{
        Database:   conf.MgoDBName,
        Collection: model.UserCollection,
        Id:         bson.ObjectIdHex(ownUserId),
      },
      DepartmentName: department.Name,
      Day:            time.Now().Format("2006-01-02"),
      Dailies:        []*model.DailyItem{&data},
      CreateTime:     time.Now(),
      UpdateTime:     time.Now(),
    }

    err = ctx.MgoDB.C(model.DailyCollection).Insert(daily)

    if err != nil {
      return errors.New(errgo.ErrCreateDailyFailed)
    }

    // 更新任务进度（不管是否成功）
    UpdateMission(ctx, missionId, bson.M{
      "progress": data.Progress,
    })

    return nil
  }

  // 如果有写过，只要在原来基础上追加一条daily
  daily := new(model.Daily)
  err = ctx.MgoDB.C(model.DailyCollection).Find(bson.M{
    "user.$id": bson.ObjectIdHex(ownUserId),
    "day":      time.Now().Format("2006-01-02"),
  }).One(daily)

  if err != nil {
    return errors.New(errgo.ErrCreateDailyFailed)
  }

  // 清缓存
  err = cache.DailyDel(ctx.Redis, ownUserId, time.Now().Format("2006-01-02"))

  if err != nil {
    return errors.New(errgo.ErrCreateDailyFailed)
  }

  // 追加并同步所有相同任务的进度
  var dailyList []model.DailyItem
  for _, item := range daily.Dailies {
    if item.MissionId == data.MissionId {
      item.Progress = data.Progress
    }
    dailyList = append(dailyList, *item)
  }

  // 这次的这条放在最后
  dailyList = append(dailyList, data)

  // 把修改好的数据塞回去
  err = ctx.MgoDB.C(model.DailyCollection).UpdateId(daily.Id, bson.M{
    "$set": bson.M{
      "dailies":    dailyList,
      "updateTime": time.Now(),
    },
  })

  if err != nil {
    return errors.New(errgo.ErrCreateDailyFailed)
  }

  // 更新任务进度（不管是否成功）
  UpdateMission(ctx, missionId, bson.M{
    "progress": data.Progress,
  })

  return nil
}

// 获取用户某一天的日报数据
func GetDailyByDay(ctx *context.New, userId string, day string) (*model.Daily, error) {

  ctx.Errgo.ErrorIfStringNotObjectId(userId, errgo.ErrUserIdError)
  ctx.Errgo.ErrorIfStringIsEmpty(day, errgo.ErrDayError)

  if err := ctx.Errgo.PopError(); err != nil {
    return nil, err
  }

  data := new(model.Daily)

  // 先从缓存查数据，找不到的话再从数据库找
  rok := cache.DailyGet(ctx.Redis, userId, day, data)

  if !rok {
    err := ctx.MgoDB.C(model.DailyCollection).Find(bson.M{
      "user.$id": bson.ObjectIdHex(userId),
      "day":      day,
    }).One(data)

    if err != nil {
      if err == mgo.ErrNotFound {
        return nil, errors.New(errgo.ErrDailyNotFound)
      }
      return nil, err
    }

    // 存到缓存
    cache.DailySet(ctx.Redis, userId, day, data)
  }

  return data, nil
}

// 更新日报内容(注意，这里的id是每一条日报的id，并非外层的这个)
func UpdateDailyContent(ctx *context.New, id string, content string) error {

  if content == "" {
    return nil
  }

  // get
  ownUserId, _ := ctx.Get(conf.OWN_USER_ID)

  // 只允许更新今天写的
  day := time.Now().Format("2006-01-02")

  // check
  ctx.Errgo.ErrorIfStringNotObjectId(id, errgo.ErrDailyIdError)

  if err := ctx.Errgo.PopError(); err != nil {
    return err
  }

  // 更新前必须成功清空缓存
  err := cache.DailyDel(ctx.Redis, ownUserId, day)

  if err != nil {
    return errors.New(errgo.ErrUpdateDailyFailed)
  }

  daily := new(model.Daily)

  err = ctx.MgoDB.C(model.DailyCollection).Find(bson.M{
    "day":         day,
    "dailies._id": bson.ObjectIdHex(id),
  }).One(daily)

  if err != nil {
    return errors.New(errgo.ErrUpdateDailyFailed)
  }

  err = cache.DailyDel(ctx.Redis, ownUserId, day)

  if err != nil {
    return errors.New(errgo.ErrUpdateDailyFailed)
  }

  // 更新数据
  err = ctx.MgoDB.C(model.DailyCollection).Update(bson.M{
    "day":         day,
    "dailies._id": bson.ObjectIdHex(id),
  }, bson.M{
    "$set": bson.M{
      "dailies.$.content": content,
      "updateTime":        time.Now(),
    },
  })

  if err != nil {
    if err == mgo.ErrNotFound {
      return errors.New(errgo.ErrDailyNotFound)
    }
    return errors.New(errgo.ErrUpdateDailyFailed)
  }

  return nil
}

// 获取日报列表
func GetDailiesList(ctx *context.New, skip int, limit int, query bson.M) (*model.DailyList, error) {

  if query != nil {
    query["dailies"] = bson.M{
      "$elemMatch": bson.M{
        "$ne": nil,
      },
    }
  }

  // check
  if limit != 0 {
    ctx.Errgo.ErrorIfIntLessThen(skip, 0, errgo.ErrSkipRange)
    ctx.Errgo.ErrorIfIntLessThen(limit, 1, errgo.ErrLimitRange)
    ctx.Errgo.ErrorIfIntMoreThen(limit, 100, errgo.ErrLimitRange)
  }

  if err := ctx.Errgo.PopError(); err != nil {
    return nil, err
  }

  data := new([]*model.Daily)

  // find it
  var err error
  if skip == 0 && limit == 0 {
    err = ctx.MgoDB.C(model.DailyCollection).Find(query).Sort("-createTime").All(data)
  } else {
    err = ctx.MgoDB.C(model.DailyCollection).Find(query).Skip(skip).Limit(limit).Sort("-createTime").All(data)
  }

  if err != nil {
    if err == mgo.ErrNotFound {
      return nil, errors.New(errgo.ErrDailyNotFound)
    }
    return nil, err
  }

  // result
  if skip == 0 && limit == 0 {
    return &model.DailyList{
      Count: len(*data),
      List:  *data,
    }, nil
  }

  // get count
  count, err := ctx.MgoDB.C(model.DailyCollection).Find(query).Count()

  if err != nil {
    return nil, errors.New(errgo.ErrDailyNotFound)
  }

  return &model.DailyList{
    List:  *data,
    Count: count,
    Skip:  skip,
    Limit: limit,
  }, nil
}

// 删除日报内容(注意，这里的id是每一条日报的id，并非外层的这个)
func DelDailyContent(ctx *context.New, id string) error {

  // get
  ownUserId, _ := ctx.Get(conf.OWN_USER_ID)

  // 只允许删除今天写的
  day := time.Now().Format("2006-01-02")

  // check
  ctx.Errgo.ErrorIfStringNotObjectId(id, errgo.ErrDailyIdError)

  if err := ctx.Errgo.PopError(); err != nil {
    return err
  }

  // 先找到这条数据
  daily, err := GetDailyByDay(ctx, ownUserId, day)

  if err != nil {
    return err
  }

  // 遍历一下找到这条日报的相关信息
  var missionId bson.ObjectId

  for _, item := range daily.Dailies {
    if item.Id.Hex() == id {
      missionId = item.MissionId
      break
    }
  }

  if missionId == "" {
    return errors.New(errgo.ErrDailyNotFound)
  }

  // 查找一下有几条这任务的的日报
  dailyCountOfMission := 0

  for _, item := range daily.Dailies {
    if item.MissionId == missionId {
      dailyCountOfMission += 1
    }
  }

  if dailyCountOfMission == 0 {
    return errors.New(errgo.ErrDailyNotFound)
  }

  // 删除前必须成功清空缓存
  err = cache.DailyDel(ctx.Redis, ownUserId, day)

  if err != nil {
    return errors.New(errgo.ErrDeleteDailyFailed)
  }

  // 删除
  err = ctx.MgoDB.C(model.DailyCollection).Update(bson.M{
    "user.$id":    bson.ObjectIdHex(ownUserId),
    "day":         day,
    "dailies._id": bson.ObjectIdHex(id),
  }, bson.M{
    "$pull": bson.M{
      "dailies": bson.M{
        "_id": bson.ObjectIdHex(id),
      },
    },
  })

  if err != nil {
    return errors.New(errgo.ErrDeleteDailyFailed)
  }

  // 如果和这任务相关的日报只有1条，那这条被删除后，任务的进度要恢复
  if dailyCountOfMission == 1 {
    // 恢复任务进度（不管是否成功）
    mission, err := GetMissionInfoById(ctx, missionId.Hex())

    if err == nil {
      UpdateMission(ctx, missionId.Hex(), bson.M{
        "progress": mission.PreProgress,
      })
    }
  }

  return nil
}

// 更新日报的任务进度
func UpdateDailyMissionProgress(ctx *context.New, missionId string, progress int) error {
  ownUserId, _ := ctx.Get(conf.OWN_USER_ID)

  // check
  ctx.Errgo.ErrorIfIntLessThen(progress, 0, errgo.ErrDailyProgressRange)
  ctx.Errgo.ErrorIfIntMoreThen(progress, 100, errgo.ErrDailyProgressRange)
  ctx.Errgo.ErrorIfStringNotObjectId(missionId, errgo.ErrDailyIdError)

  if err := ctx.Errgo.PopError(); err != nil {
    return err
  }

  // find
  var data = new(model.Daily)
  err := ctx.MgoDB.C(model.DailyCollection).Find(bson.M{
    "user.$id": bson.ObjectIdHex(ownUserId),
    "day":      time.Now().Format("2006-01-02"),
  }).One(data)

  if err != nil {
    if err == mgo.ErrNotFound {
      return errors.New(errgo.ErrDailyNotFound)
    }
    return errors.New(errgo.ErrUpdateDailyFailed)
  }

  // 更新前必须成功清空缓存
  err = cache.DailyDel(ctx.Redis, ownUserId, data.Day)

  if err != nil {
    return errors.New(errgo.ErrUpdateDailyFailed)
  }

  // 同步所有相同任务的进度
  var dailyList []model.DailyItem
  for _, item := range data.Dailies {
    if item.MissionId.Hex() == missionId {
      item.Progress = progress
    }
    dailyList = append(dailyList, *item)
  }

  // 把修改好的数据塞回去
  err = ctx.MgoDB.C(model.DailyCollection).UpdateId(data.Id, bson.M{
    "$set": bson.M{
      "dailies":    dailyList,
      "updateTime": time.Now(),
    },
  })

  if err != nil {
    return errors.New(errgo.ErrUpdateDailyFailed)
  }

  // 更新任务进度（不管是否成功）
  UpdateMission(ctx, missionId, bson.M{
    "progress": progress,
  })

  return nil
}
