package service

import (
  "errors"
  `gopkg.in/mgo.v2`
  "gopkg.in/mgo.v2/bson"
  "time"
  "workerbook/db"
  "workerbook/model"
)

// Insert user info into database.
func CreateUser(data model.User) error {
  db, close, err := db.CloneDB()

  if err != nil {
    return err
  } else {
    defer close()
  }

  // check the data is error or not.
  if data.UserName == "" {
    return errors.New("用户名不能为空")
  } else if data.Password == "" {
    return errors.New("密码不能为空")
  } else if data.NickName == "" {
    return errors.New("昵称不能为空")
  } else if !bson.IsObjectIdHex(data.Gid) {
    return errors.New("分组号错误")
  }

  // supplement other data.
  if data.Role == 0 {
    data.Role = 1
  }
  data.CreateTime = time.Now()

  // username must be the only.
  count, err := db.C(model.UserCollection).Find(bson.M{"username": data.UserName}).Count()

  if err != nil {
    return err
  }

  if count > 0 {
    return errors.New("已存在相同的用户")
  }

  // group must be exist.
  group := model.Group{}
  db.C(model.GroupCollection).FindId(bson.ObjectIdHex(data.Gid)).One(&group)

  if group.Id == "" {
    return errors.New("找不到该分组")
  }

  data.GroupName = group.Name

  // create a new object id.
  if data.Id == "" {
    data.Id = bson.NewObjectId()
  }

  // insert it.
  data.Exist = true

  err = db.C(model.UserCollection).Insert(data)

  if err != nil {
    return err
  }

  // refresh group count
  RefreshGroupCount(group.Id)

  return nil
}

// user login by username and password
func UserLogin(username string, password string) (string, error) {
  db, close, err := db.CloneDB()

  if err != nil {
    return "", err
  } else {
    defer close()
  }

  data := model.UserResult{}

  err = db.C(model.UserCollection).Find(bson.M{
    "username": username,
    "password": password,
  }).One(&data)

  if err != nil {
    if err == mgo.ErrNotFound {
      return "", errors.New("用户名或密码错误")
    }
    return "", err
  } else {
    return data.Id.Hex(), nil
  }
}

// Query user info by id.
func GetUserInfoById(id bson.ObjectId) (model.UserResult, error) {
  db, close, err := db.CloneDB()

  data := model.UserResult{}

  if err != nil {
    return data, err
  } else {
    defer close()
  }

  err = db.C(model.UserCollection).FindId(id).One(&data)

  if err != nil {
    return data, err
  }

  return data, nil
}

// Query users list with skip and limit.
// it will find all of users when 'gid' is empty.
func GetUsersList(gid string, skip int, limit int) ([]model.UserResult, error) {
  db, close, err := db.CloneDB()

  if err != nil {
    return nil, err
  } else {
    defer close()
  }

  data := []model.UserResult{}

  if limit < 0 {
    limit = 0
  }

  // create condition sql
  sql := bson.M{}
  if gid != "" {
    sql["gid"] = gid
  }

  // find it
  err = db.C(model.UserCollection).Find(sql).Skip(skip).Limit(limit).All(&data)

  if err != nil {
    return nil, err
  }

  return data, nil
}
