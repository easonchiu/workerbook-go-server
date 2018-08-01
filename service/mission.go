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
  "workerbook/util"
)

// 创建任务
func CreateMission(ctx *context.New, data model.Mission, projectId string, userId string) error {

  // 是否存在的标志
  data.Exist = true

  data.Progress = 0
  data.PreProgress = 0
  data.Status = 1
  data.CreateTime = time.Now()
  ownUserId, _ := ctx.Get(conf.OWN_USER_ID)
  data.Editor = mgo.DBRef{
    Database:   conf.MgoDBName,
    Collection: model.UserCollection,
    Id:         bson.ObjectIdHex(ownUserId),
  }
  data.ChartTime = "19900101"
  data.EditTime = time.Now()

  // 修改任务结束时间到今天的最晚时间，即23:59:59
  {
    y1, m1, d1 := data.Deadline.Local().Date()
    data.Deadline = time.Date(y1, m1, d1, 23, 59, 59, 0, time.Local)
  }

  // check
  errgo.ErrorIfStringIsEmpty(data.Name, errgo.ErrMissionNameEmpty)
  errgo.ErrorIfLenMoreThen(data.Name, 30, errgo.ErrMissionNameTooLong)
  errgo.ErrorIfTimeEarlierThen(data.Deadline, time.Now(), errgo.ErrMissionDeadlineTooSoon)
  errgo.ErrorIfStringNotObjectId(projectId, errgo.ErrProjectIdError)
  errgo.ErrorIfStringNotObjectId(userId, errgo.ErrUserIdError)

  if err := errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return err
  }

  // 查找项目
  project, err := GetProjectInfoById(ctx, projectId)

  if err != nil {
    return err
  }

  // 任务结束时间不能晚于项目结束时间
  errgo.ErrorIfTimeLaterThen(data.Deadline, project.Deadline, errgo.ErrMissionDeadlineTooLate)

  if err = errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return err
  }

  data.Id = bson.NewObjectId()
  data.Project = mgo.DBRef{
    Database:   conf.MgoDBName,
    Collection: model.ProjectCollection,
    Id:         bson.ObjectIdHex(projectId),
  }

  // find user
  _, err = GetUserInfoById(ctx, userId)

  if err != nil {
    return err
  }

  data.User = mgo.DBRef{
    Database:   conf.MgoDBName,
    Collection: model.UserCollection,
    Id:         bson.ObjectIdHex(userId),
  }

  // insert it.
  err = ctx.MgoDB.C(model.MissionCollection).Insert(data)

  if err != nil {
    return errors.New(errgo.ErrCreateMissionFailed)
  }

  // 成功后要在项目的数据中关联这条数据
  err = ctx.MgoDB.C(model.ProjectCollection).UpdateId(bson.ObjectIdHex(projectId), bson.M{
    "$push": bson.M{
      "missions": mgo.DBRef{
        Database:   conf.MgoDBName,
        Collection: model.MissionCollection,
        Id:         data.Id,
      },
    },
  })

  // 插入失败就删除任务
  if err != nil {
    ctx.MgoDB.C(model.MissionCollection).RemoveId(data.Id)
    return errors.New(errgo.ErrCreateMissionFailed)
  }

  return nil
}

// 更新任务
func UpdateMission(ctx *context.New, id string, data bson.M) error {

  if data == nil {
    return errors.New(errgo.ErrServerError)
  }

  // 限制更新字段
  util.Only(
    data,
    util.Keys{
      "name":        util.TypeString,
      "preProgress": util.TypeInt,
      "progress":    util.TypeInt,
      "chartTime":   util.TypeTime,
      "userId":      util.TypeString,
      "deadline":    util.TypeTime,
      "status":      util.TypeInt,
      "project":     util.TypeBsonM,
      "exist":       util.TypeBool,
      "editor":      util.TypeBsonM,
      "editTime":    util.TypeTime,
    },
  )

  // check
  errgo.ErrorIfStringNotObjectId(id, errgo.ErrMissionIdError)

  if name, ok := util.GetString(data, "name"); ok {
    errgo.ErrorIfLenMoreThen(name, 15, errgo.ErrMissionNameTooLong)
  }

  if userId, ok := util.GetString(data, "userId"); ok {
    errgo.ErrorIfStringNotObjectId(userId, errgo.ErrUserIdError)
  }

  if progress, ok := util.GetInt(data, "progress"); ok {
    errgo.ErrorIfIntLessThen(progress, 0, errgo.ErrMissionProgressRange)
    errgo.ErrorIfIntMoreThen(progress, 100, errgo.ErrMissionProgressRange)
  }

  // check
  if err := errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return err
  }

  // 找到这个任务和项目
  if deadline, ok := util.GetTime(data, "deadline"); ok {
    mission, err := GetMissionInfoById(ctx, id)

    if err != nil {
      return err
    }

    project, err := FindProjectRef(ctx, &mission.Project)

    if err != nil {
      return err
    }

    // 如果有设置结束时间，将时间改为这天的最晚时间，即23:59:59
    // 任务截至时间不能晚于项目截至时间，不能早于当前时间
    y1, m1, d1 := deadline.Local().Date()
    t := time.Date(y1, m1, d1, 23, 59, 59, 0, time.Local)
    data["deadline"] = t

    errgo.ErrorIfTimeEarlierThen(t, time.Now(), errgo.ErrMissionDeadlineTooSoon)
    errgo.ErrorIfTimeLaterThen(t, project.Deadline, errgo.ErrMissionDeadlineTooLate)
  } else {
    // 没有项目id时不允许修改任务结束时间
    delete(data, "deadline")
  }

  // 查找执行人并存入执行人ref
  if userId, ok := util.GetString(data, "userId"); ok {
    user, err := GetUserInfoById(ctx, userId)

    if err != nil {
      return err
    }

    data["user"] = mgo.DBRef{
      Database:   conf.MgoDBName,
      Collection: model.UserCollection,
      Id:         user.Id,
    }
  }

  // check
  if err := errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return err
  }

  // 先要清缓存，清成功后才可以更新数据
  err := cache.MissionDel(ctx.Redis, id)

  if err != nil {
    if exist, ok := util.GetBool(data, "exist"); ok && exist == false {
      return errors.New(errgo.ErrDeleteMissionFailed)
    }
    return errors.New(errgo.ErrUpdateMissionFailed)
  }

  // update
  ownUserId, _ := ctx.Get(conf.OWN_USER_ID)
  data["editor.$id"] = bson.ObjectIdHex(ownUserId)
  data["editTime"] = time.Now()

  err = ctx.MgoDB.C(model.MissionCollection).UpdateId(bson.ObjectIdHex(id), bson.M{
    "$set": data,
  })

  if err != nil {
    if err == mgo.ErrNotFound {
      return errors.New(errgo.ErrMissionNotFound)
    }
    if exist, ok := util.GetBool(data, "exist"); ok && exist == false {
      return errors.New(errgo.ErrDeleteMissionFailed)
    }
    return errors.New(errgo.ErrUpdateMissionFailed)
  }

  return nil
}

// 查找任务列表(当limit都为0时，查找全部)
func GetMissionsList(ctx *context.New, skip int, limit int, query bson.M) (*model.MissionList, error) {

  // check
  if limit != 0 {
    errgo.ErrorIfIntLessThen(skip, 0, errgo.ErrSkipRange)
    errgo.ErrorIfIntLessThen(limit, 1, errgo.ErrLimitRange)
    errgo.ErrorIfIntMoreThen(limit, 100, errgo.ErrLimitRange)
  }

  if err := errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return nil, err
  }

  data := new([]model.Mission)
  query["exist"] = true

  // find it
  var err error
  if skip == 0 && limit == 0 {
    err = ctx.MgoDB.C(model.MissionCollection).Find(query).Sort("-createTime").All(data)
  } else {
    err = ctx.MgoDB.C(model.MissionCollection).Find(query).Skip(skip).Limit(limit).Sort("-createTime").All(data)
  }

  if err != nil {
    if err == mgo.ErrNotFound {
      return nil, errors.New(errgo.ErrMissionNotFound)
    }
    return nil, err
  }

  // result

  if skip == 0 && limit == 0 {
    return &model.MissionList{
      List: data,
    }, nil
  }

  // get count
  count, err := ctx.MgoDB.C(model.MissionCollection).Find(query).Count()

  if err != nil {
    return nil, errors.New(errgo.ErrMissionNotFound)
  }

  return &model.MissionList{
    List:  data,
    Count: count,
    Skip:  skip,
    Limit: limit,
  }, nil
}

// 根据id查找任务
func GetMissionInfoById(ctx *context.New, id string) (*model.Mission, error) {

  // check
  errgo.ErrorIfStringNotObjectId(id, errgo.ErrMissionIdError)

  if err := errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return nil, err
  }

  data := new(model.Mission)

  // 先从缓存取数据，如果缓存没取到，从数据库取
  rok := cache.MissionGet(ctx.Redis, id, data)
  if !rok {
    err := ctx.MgoDB.C(model.MissionCollection).FindId(bson.ObjectIdHex(id)).One(data)

    if err != nil {
      if err == mgo.ErrNotFound {
        return nil, errors.New(errgo.ErrMissionNotFound)
      }
      return nil, err
    }

    // 存到缓存
    cache.MissionSet(ctx.Redis, data)
  }

  return data, nil
}
