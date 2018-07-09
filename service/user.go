package service

import (
  "errors"
  "github.com/gin-gonic/gin"
  "gopkg.in/mgo.v2"
  "gopkg.in/mgo.v2/bson"
  "time"
  "workerbook/errno"
  "workerbook/model"
  "workerbook/mongo"
)

// Insert user info into database.
func CreateUser(data model.User) error {
  db, closer, err := mongo.CloneDB()

  if err != nil {
    return err
  } else {
    defer closer()
  }

  // supplement other data.
  if data.Role != 0 && data.Role != 1 && data.Role != 2 {
    data.Role = 1
  }

  data.CreateTime = time.Now()

  // username must be the only.
  count, err := db.C(model.UserCollection).Find(bson.M{"username": data.UserName}).Count()

  if err != nil {
    return errors.New(errno.ErrCreateUserFailed)
  }

  if count > 0 {
    return errors.New(errno.ErrSameUsername)
  }

  // nickname must be the only.
  count, err = db.C(model.UserCollection).Find(bson.M{"nickname": data.NickName}).Count()

  if err != nil {
    return errors.New(errno.ErrCreateUserFailed)
  }

  if count > 0 {
    return errors.New(errno.ErrSameNickname)
  }

  // department must be exist.
  department := new(model.Department)
  db.FindRef(&data.Department).One(department)

  if department.Name == "" {
    return errors.New(errno.ErrDepartmentNotFound)
  }

  // set status
  data.Status = 1

  // insert it.
  err = db.C(model.UserCollection).Insert(data)

  if err != nil {
    return errors.New(errno.ErrCreateUserFailed)
  }

  // refresh user count in department
  // RefreshGroupCount(group.Id)

  return nil
}

// update user info.
func UpdateUser(id bson.ObjectId, m bson.M) error {
  db, closer, err := mongo.CloneDB()

  if err != nil {
    return err
  } else {
    defer closer()
  }

  // 帐号唯一
  count, err := db.C(model.UserCollection).Find(bson.M{"username": m["username"]}).Count()

  if err != nil {
    return errors.New(errno.ErrUpdateUserFailed)
  }

  if count > 0 {
    return errors.New(errno.ErrSameUsername)
  }

  // 姓名唯一
  count, err = db.C(model.UserCollection).Find(bson.M{
    "nickname": m["nickname"],
    "_id": bson.M{
      "$ne": id,
    },
  }).Count()

  if err != nil {
    return errors.New(errno.ErrUpdateUserFailed)
  }

  if count > 0 {
    return errors.New(errno.ErrSameNickname)
  }

  // 部门必选并必须存在
  ref := m["department"].(mgo.DBRef)
  count, err = db.FindRef(&ref).Count()

  if count == 0 {
    return errors.New(errno.ErrDepartmentNotFound)
  }

  // 更新数据
  err = db.C(model.UserCollection).UpdateId(id, bson.M{
    "$set": m,
  })

  if err != nil {
    return errors.New(errno.ErrUpdateUserFailed)
  }

  // refresh user count in department
  // RefreshGroupCount(group.Id)

  return nil
}

// 用户登录并返回用户id
func UserLogin(username string, password string) (id string, err error) {
  db, closer, err := mongo.CloneDB()

  if err != nil {
    return "", err
  } else {
    defer closer()
  }

  data := new(model.User)

  err = db.C(model.UserCollection).Find(bson.M{
    "username": username,
    "password": password,
  }).One(&data)

  if err != nil {
    if err == mgo.ErrNotFound {
      return "", errors.New(errno.ErrUsernameOrPasswordError)
    }
    return "", err
  } else {
    return data.Id.Hex(), nil
  }
}

// 根据id查找用户信息
func GetUserInfoById(id bson.ObjectId) (gin.H, error) {
  db, closer, err := mongo.CloneDB()

  data := new(model.User)

  if err != nil {
    return nil, err
  } else {
    defer closer()
  }

  err = db.C(model.UserCollection).FindId(id).One(&data)

  if err != nil {
    if err == mgo.ErrNotFound {
      return nil, errors.New(errno.ErrUserNotFound)
    }
    return nil, err
  }

  return data.GetMap(db), nil
}

// 获取用户列表
// 如果departmentId为空，查找所有用户
func GetUsersList(departmentId string, skip int, limit int) (gin.H, error) {
  db, closer, err := mongo.CloneDB()

  if err != nil {
    return nil, err
  } else {
    defer closer()
  }

  data := new([]model.User)

  if limit < 0 {
    limit = 0
  } else if limit > 100 {
    limit = 100
  }

  // 设置查找sql
  sql := bson.M{}
  if departmentId != "" {
    sql["departmentId"] = departmentId
  }

  // find it
  err = db.C(model.UserCollection).Find(sql).Skip(skip).Limit(limit).Sort("-createTime").All(data)

  if err != nil {
    if err == mgo.ErrNotFound {
      return nil, errors.New(errno.ErrUserNotFound)
    }
    return nil, err
  }

  // get count
  count, err := db.C(model.UserCollection).Find(sql).Count()

  if err != nil {
    return nil, errors.New(errno.ErrUserNotFound)
  }

  // result
  var list []gin.H

  for _, r := range *data {
    list = append(list, r.GetMap(db))
  }

  return gin.H{
    "list":  list,
    "count": count,
    "skip":  skip,
    "limit": limit,
  }, nil
}
