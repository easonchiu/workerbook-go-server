package service

import (
  "errors"
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "time"
  "workerbook/conf"
  "workerbook/errgo"
  "workerbook/model"
  "workerbook/mongo"
)

// 创建任务
func CreateMission(data model.Mission, projectId string, userId string) error {
  db, close, err := mongo.CloneDB()

  if err != nil {
    return err
  } else {
    defer close()
  }

  // 是否存在的标志
  data.Exist = true
  data.CreateTime = time.Now()

  // check
  errgo.ErrorIfStringIsEmpty(data.Name, errgo.ErrMissionNameEmpty)
  errgo.ErrorIfLenMoreThen(data.Name, 15, errgo.ErrMissionNameTooLong)
  errgo.ErrorIfTimeEarlierThen(data.Deadline, time.Now(), errgo.ErrMissionDeadlineTooSoon)
  errgo.ErrorIfStringNotObjectId(projectId, errgo.ErrProjectIdError)
  errgo.ErrorIfStringNotObjectId(userId, errgo.ErrUserIdError)

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
    Database:   conf.DBName,
  }

  // insert it.
  err = db.C(model.MissionCollection).Insert(data)

  if err != nil {
    return errors.New(errgo.ErrCreateMissionFailed)
  }

  // 成功后要在项目的数据中关联这条数据
  err = db.C(model.ProjectCollection).UpdateId(data.ProjectId, bson.M{
    "$push": bson.M{
      "missions": mgo.DBRef{
        Id:         data.Id,
        Collection: model.MissionCollection,
        Database:   conf.DBName,
      },
    },
  })

  // 插入失败就删除任务
  if err != nil {
    db.C(model.MissionCollection).RemoveId(data.Id)
    return errors.New(errgo.ErrCreateMissionFailed)
  }

  return nil
}

// 更新任务
func UpdateMission(id string, data bson.M) error {
  db, close, err := mongo.CloneDB()

  if err != nil {
    return err
  } else {
    defer close()
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
    var project = new(model.Project)

    err = db.C(model.ProjectCollection).Find(bson.M{
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
      Database:   conf.DBName,
    }
  }

  // check
  if err = errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return err
  }

  // update
  if projectId, ok := data["projectId"]; ok {
    data["projectId"] = bson.ObjectIdHex(projectId.(string))
  }

  err = db.C(model.MissionCollection).UpdateId(bson.ObjectIdHex(id), bson.M{
    "$set": data,
  })

  if err != nil {
    return errors.New(errgo.ErrUpdateMissionFailed)
  }

  return nil
}

// 根据id查找任务
func GetMissionInfoById(id string) (gin.H, error) {
  db, close, err := mongo.CloneDB()

  if err != nil {
    return nil, err
  } else {
    defer close()
  }

  // check
  errgo.ErrorIfStringNotObjectId(id, errgo.ErrMissionIdError)

  if err = errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return nil, err
  }

  data := new(model.Mission)

  err = db.C(model.MissionCollection).FindId(bson.ObjectIdHex(id)).One(data)

  if err != nil {
    if err == mgo.ErrNotFound {
      return nil, errors.New(errgo.ErrMissionNotFound)
    }
    return nil, err
  }

  return data.GetMap(db, "user"), nil
}
