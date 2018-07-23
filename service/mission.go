package service

import (
  "errors"
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "time"
  "workerbook/cache"
  "workerbook/conf"
  "workerbook/db"
  "workerbook/errgo"
  "workerbook/model"
)

// 创建任务
func CreateMission(data model.Mission, projectId string, userId string) error {
  mg, closer, err := db.CloneMgoDB()

  if err != nil {
    return err
  } else {
    defer closer()
  }

  // 是否存在的标志
  data.Exist = true
  data.CreateTime = time.Now()

  // check
  errgo.ErrorIfStringIsEmpty(data.Name, errgo.ErrMissionNameEmpty)
  errgo.ErrorIfLenMoreThen(data.Name, 30, errgo.ErrMissionNameTooLong)
  errgo.ErrorIfTimeEarlierThen(data.Deadline, time.Now(), errgo.ErrMissionDeadlineTooSoon)
  errgo.ErrorIfStringNotObjectId(projectId, errgo.ErrProjectIdError)
  errgo.ErrorIfStringNotObjectId(userId, errgo.ErrUserIdError)

  if err = errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return err
  }

  // 查找项目
  project, err := GetProjectInfoById(projectId)

  if err != nil {
    return err
  }

  // 任务结束时间不能晚于项目结束时间
  errgo.ErrorIfTimeLaterThen(data.Deadline, project["deadline"].(time.Time), errgo.ErrMissionDeadlineTooLate)

  if err = errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return err
  }


  data.Id = bson.NewObjectId()
  data.ProjectId = bson.ObjectIdHex(projectId)

  // find user
  _, err = GetUserInfoById(userId)

  if err != nil {
    return err
  }

  data.User = mgo.DBRef{
    Id:         bson.ObjectIdHex(userId),
    Collection: model.UserCollection,
    Database:   conf.MgoDBName,
  }

  // insert it.
  err = mg.C(model.MissionCollection).Insert(data)

  if err != nil {
    return errors.New(errgo.ErrCreateMissionFailed)
  }

  // 成功后要在项目的数据中关联这条数据
  err = mg.C(model.ProjectCollection).UpdateId(data.ProjectId, bson.M{
    "$push": bson.M{
      "missions": mgo.DBRef{
        Id:         data.Id,
        Collection: model.MissionCollection,
        Database:   conf.MgoDBName,
      },
    },
  })

  // 插入失败就删除任务
  if err != nil {
    mg.C(model.MissionCollection).RemoveId(data.Id)
    return errors.New(errgo.ErrCreateMissionFailed)
  }

  return nil
}

// 更新任务
func UpdateMission(id string, data bson.M) error {
  mg, closer, err := db.CloneMgoDB()

  if err != nil {
    return err
  } else {
    defer closer()
  }

  // check
  errgo.ErrorIfStringNotObjectId(id, errgo.ErrMissionIdError)

  if name, ok := data["name"]; ok {
    errgo.ErrorIfLenMoreThen(name.(string), 15, errgo.ErrMissionNameTooLong)
  }

  if deadline, ok := data["deadline"]; ok {
    errgo.ErrorIfTimeEarlierThen(deadline.(time.Time), time.Now(), errgo.ErrMissionDeadlineTooSoon)
  }

  if userId, ok := data["userId"]; ok {
    errgo.ErrorIfStringNotObjectId(userId.(string), errgo.ErrUserIdError)
  }

  if projectId, ok := data["projectId"]; ok {
    errgo.ErrorIfStringNotObjectId(projectId.(string), errgo.ErrProjectIdError)
  }

  // check
  if err = errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return err
  }

  // 查找项目
  if projectId, ok := data["projectId"]; ok {
    if !bson.IsObjectIdHex(projectId.(string)) {
      return errors.New(errgo.ErrProjectIdError)
    }
    var project = new(model.Project)
    err = mg.C(model.ProjectCollection).Find(bson.M{
      "_id":   bson.ObjectIdHex(projectId.(string)),
      "exist": true,
    }).One(project)

    if err != nil {
      return errors.New(errgo.ErrProjectNotFound)
    }

    // 任务截至时间不能晚于项目截至时间
    if deadline, ok := data["deadline"]; ok {
      errgo.ErrorIfTimeLaterThen(deadline.(time.Time), project.Deadline, errgo.ErrMissionDeadlineTooLate)
    }
  }

  // 查找执行人
  if userId, ok := data["userId"]; ok {
    _, err := GetUserInfoById(userId.(string))

    if err != nil {
      return err
    }

    data["user"] = mgo.DBRef{
      Id:         bson.ObjectIdHex(userId.(string)),
      Collection: model.UserCollection,
      Database:   conf.MgoDBName,
    }
  }

  // check
  if err = errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return err
  }

  // set project id
  if projectId, ok := data["projectId"]; ok {
    data["projectId"] = bson.ObjectIdHex(projectId.(string))
  }

  // 先要清缓存，清成功后才可以更新数据
  err = cache.MissionDel(id)

  if err != nil {
    return errors.New(errgo.ErrUpdateMissionFailed)
  }

  // update
  data["exist"] = true
  err = mg.C(model.MissionCollection).UpdateId(bson.ObjectIdHex(id), bson.M{
    "$set": data,
  })

  if err != nil {
    return errors.New(errgo.ErrUpdateMissionFailed)
  }

  return nil
}

// 根据id查找任务
func GetMissionInfoById(id string, refs ... string) (gin.H, error) {
  mg, closer, err := db.CloneMgoDB()

  if err != nil {
    return nil, err
  } else {
    defer closer()
  }

  // check
  errgo.ErrorIfStringNotObjectId(id, errgo.ErrMissionIdError)

  if err = errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return nil, err
  }

  data := new(model.Mission)

  // 先从缓存取数据，如果缓存没取到，从数据库取
  rok := cache.MissionGet(id, data)
  if !rok {
    err = mg.C(model.MissionCollection).FindId(bson.ObjectIdHex(id)).One(data)
  }

  if err != nil {
    if err == mgo.ErrNotFound {
      return nil, errors.New(errgo.ErrMissionNotFound)
    }
    return nil, err
  }

  return data.GetMap(FindRef(mg), refs...), nil
}
