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
func UpdateUser(m bson.M) error {
  db, closer, err := mongo.CloneDB()

  if err != nil {
    return err
  } else {
    defer closer()
  }

  // username must be the only.
  count, err := db.C(model.UserCollection).Find(bson.M{"username": m["username"]}).Count()

  if err != nil {
    return errors.New(errno.ErrCreateUserFailed)
  }

  if count > 0 {
    return errors.New(errno.ErrSameUsername)
  }

  // nickname must be the only.
  count, err = db.C(model.UserCollection).Find(bson.M{"nickname": m["nickname"]}).Count()

  if err != nil {
    return errors.New(errno.ErrCreateUserFailed)
  }

  if count > 0 {
    return errors.New(errno.ErrSameNickname)
  }

  // department must be exist.
  count, err = db.C(model.DepartmentCollection).FindId(bson.ObjectIdHex(m["departmentId"].(string))).Count()

  if count == 0 {
    return errors.New(errno.ErrDepartmentNotFound)
  }

  // insert it.
  err = db.C(model.UserCollection).UpdateId(bson.ObjectIdHex(m["id"].(string)), bson.M{
    "$set": m,
  })

  if err != nil {
    return errors.New(errno.ErrUpdateUserFailed)
  }

  // refresh user count in department
  // RefreshGroupCount(group.Id)

  return nil
}

// user login by username and password
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

// Query user info by id.
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

  return gin.H{
    "department": data.GetResultOneWithRef(db),
  }, nil
}

// Query users list with skip and limit.
// it will find all of users when 'departmentId' is empty.
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

  // create condition sql
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
    list = append(list, r.GetResultOneWithRef(db))
  }

  return gin.H{
    "list": list,
    "count": count,
    "skip": skip,
    "limit": limit,
  }, nil
}
