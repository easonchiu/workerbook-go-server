package service

import (
  `errors`
  `gopkg.in/mgo.v2/bson`
  `time`
  `workerbook/db`
  `workerbook/model`
)

// Query daily info by id.
func GetDailyInfoById(id bson.ObjectId) (*model.Daily, error) {
  db, close, err := db.CloneDB()

  data := new(model.Daily)

  if err != nil {
    return data, err
  } else {
    defer close()
  }

  err = db.C(model.DailyCollection).FindId(id).One(data)

  if err != nil {
    return data, err
  }

  return data, nil
}

// Query dailies list with skip and limit.
func GetDailiesList(skip int, limit int) (*[]model.Daily, error) {
  db, close, err := db.CloneDB()

  if err != nil {
    return nil, err
  } else {
    defer close()
  }

  data := new([]model.Daily)

  if limit < 0 {
    limit = 0
  }

  err = db.C(model.DailyCollection).Find(bson.M{
    "dailyList": bson.M{
      "$not": bson.M{
        "$size": 0,
      },
    },
  }).Sort("-updateTime").Skip(skip).Limit(limit).All(data)

  if err != nil {
    return nil, err
  }

  return data, nil
}

// Query today's daily with some user.
func GetUserTodayDaily(uid bson.ObjectId) (*model.Daily, error) {
  db, close, err := db.CloneDB()

  data := new(model.Daily)

  if err != nil {
    return data, err
  } else {
    defer close()
  }

  // check user is existed first.
  _, err = GetUserInfoById(uid)

  if err != nil {
    return data, errors.New("没有相关的用户")
  }

  // time to string
  today := time.Now().Format("2006-01-02")

  // find daily with uid and string time.
  err = db.C(model.DailyCollection).Find(bson.M{"uid": uid.Hex(), "day": today}).One(data)

  return data, err
}

// create today daily
func CreateMyTodayDaily(uid bson.ObjectId) (*model.Daily, error) {
  db, close, err := db.CloneDB()

  if err != nil {
    return nil, err
  } else {
    defer close()
  }

  // check user is existed first.
  userInfo, err := GetUserInfoById(uid)

  if err != nil {
    return nil, errors.New("没有相关的用户")
  }

  // time to string
  today := time.Now().Format("2006-01-02")

  data := &model.Daily{
    Id:         bson.NewObjectId(),
    Uid:        userInfo.Id.Hex(),
    GroupId:    userInfo.GroupId,
    Day:        today,
    DailyList:  []model.DailyItem{},
    CreateTime: time.Now(),
    UpdateTime: time.Now(),
  }

  err = db.C(model.DailyCollection).Insert(data)

  return data, err
}

// append daily item into users daily list.
func AppendDailyItemIntoUsersDailyList(data model.DailyItem, id bson.ObjectId) error {
  db, close, err := db.CloneDB()

  if err != nil {
    return err
  } else {
    defer close()
  }

  err = db.C(model.DailyCollection).UpdateId(id, bson.M{
    "$push": bson.M{
      "dailyList": data,
    },
  })

  return err
}

// delete daily item in today from users daily list.
func DeleteTodayDailyItemFromUsersDailyList(uid bson.ObjectId, itemId bson.ObjectId) error {
  db, close, err := db.CloneDB()

  if err != nil {
    return err
  } else {
    defer close()
  }

  // find user's daily in today.
  dailyInfo, err := GetUserTodayDaily(uid)

  // find the data is in today's daily or not.
  include := false
  for _, i := range dailyInfo.DailyList {
    if i.Id == itemId {
      include = true
      break
    }
  }

  if !include {
    return errors.New("没有相关的日报内容")
  }

  // has related data, find and delete it.
  err = db.C(model.DailyCollection).UpdateId(dailyInfo.Id, bson.M{
    "$pull": bson.M{
      "dailyList": bson.M{
        "_id": itemId,
      },
    },
  })

  return err
}
