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

// 创建用户
func CreateUser(data model.User, departmentId string) error {
  db, closer, err := mongo.CloneDB()

  if err != nil {
    return err
  } else {
    defer closer()
  }

  // 是否存在的标志
  data.Exist = true
  data.CreateTime = time.Now()

  // 用户状态
  data.Status = 1

  // check
  errgo.ErrorIfStringIsEmpty(data.NickName, errgo.ErrNicknameEmpty)
  errgo.ErrorIfLenLessThen(data.NickName, 2, errgo.ErrNicknameTooShort)
  errgo.ErrorIfLenMoreThen(data.NickName, 14, errgo.ErrNicknameTooLong)
  errgo.ErrorIfStringNotObjectId(departmentId, errgo.ErrDepartmentIdError)
  errgo.ErrorIfStringIsEmpty(data.Title, errgo.ErrUserTitleIsEmpty)
  errgo.ErrorIfLenMoreThen(data.Title, 14, errgo.ErrUserTitleTooLong)
  if data.Role != 1 && data.Role != 2 && data.Role != 3 {
    return errors.New(errgo.ErrUserRoleError)
  }
  errgo.ErrorIfStringIsEmpty(data.UserName, errgo.ErrUsernameEmpty)
  errgo.ErrorIfLenLessThen(data.UserName, 6, errgo.ErrUsernameTooShort)
  errgo.ErrorIfLenMoreThen(data.UserName, 14, errgo.ErrUsernameTooLong)
  errgo.ErrorIfStringIsEmpty(data.Password, errgo.ErrPasswordEmpty)

  if err = errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return err
  }

  // save department id
  data.Department = mgo.DBRef{
    Id:         bson.ObjectIdHex(departmentId),
    Collection: model.DepartmentCollection,
    Database:   conf.DBName,
  }

  // username must be the only.
  count, err := db.C(model.UserCollection).Find(model.User{
    UserName: data.UserName,
    Exist:    true,
  }).Count()

  if err != nil {
    return errors.New(errgo.ErrCreateUserFailed)
  }

  if count > 0 {
    return errors.New(errgo.ErrSameUsername)
  }

  // nickname must be the only.
  count, err = db.C(model.UserCollection).Find(model.User{
    NickName: data.NickName,
    Exist:    true,
  }).Count()

  if err != nil {
    return errors.New(errgo.ErrCreateUserFailed)
  }

  if count > 0 {
    return errors.New(errgo.ErrSameNickname)
  }

  // department must be exist.
  department := new(model.Department)
  db.FindRef(&data.Department).One(department)

  if department.Name == "" {
    return errors.New(errgo.ErrDepartmentNotFound)
  }

  // insert it.
  err = db.C(model.UserCollection).Insert(data)

  if err != nil {
    return errors.New(errgo.ErrCreateUserFailed)
  }

  // update count in department
  UpdateDepartmentsUserCount()

  return nil
}

// 更新用户
func UpdateUser(id string, data bson.M) error {
  db, closer, err := mongo.CloneDB()

  if err != nil {
    return err
  } else {
    defer closer()
  }

  // check
  errgo.ErrorIfStringNotObjectId(id, errgo.ErrUserIdError)

  if nickname, ok := data["nickname"]; ok {
    errgo.ErrorIfLenLessThen(nickname.(string), 2, errgo.ErrNicknameTooShort)
    errgo.ErrorIfLenMoreThen(nickname.(string), 14, errgo.ErrNicknameTooLong)
  }

  if departmentId, ok := data["department.$id"]; ok {
    errgo.ErrorIfStringNotObjectId(departmentId.(string), errgo.ErrDepartmentIdError)
  }

  if title, ok := data["title"]; ok {
    errgo.ErrorIfLenMoreThen(title.(string), 14, errgo.ErrUserTitleTooLong)
  }

  if role, ok := data["role"]; ok {
    role := role.(int)
    if role != 1 && role != 2 && role != 3 {
      return errors.New(errgo.ErrUserRoleError)
    }
  }

  if err = errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return err
  }

  // 姓名唯一
  if nickname, ok := data["nickname"]; ok {
    count, err := db.C(model.UserCollection).Find(bson.M{
      "nickname": nickname,
      "exist":    true,
      "_id": bson.M{
        "$ne": bson.ObjectIdHex(id),
      },
    }).Count()

    if err != nil {
      return errors.New(errgo.ErrUpdateUserFailed)
    }

    if count > 0 {
      return errors.New(errgo.ErrSameNickname)
    }
  }

  // 部门必须存在
  if departmentId, ok := data["department.$id"]; ok {
    ref := mgo.DBRef{
      Id:         bson.ObjectIdHex(departmentId.(string)),
      Collection: model.DepartmentCollection,
      Database:   conf.DBName,
    }

    count, err := db.FindRef(&ref).Select(bson.M{"exist": true}).Count()

    if err != nil {
      return errors.New(errgo.ErrUpdateUserFailed)
    }

    if count == 0 {
      return errors.New(errgo.ErrDepartmentNotFound)
    }

    data["department.$id"] = bson.ObjectIdHex(departmentId.(string))
  }


  // 更新数据
  data["exist"] = true
  err = db.C(model.UserCollection).UpdateId(bson.ObjectIdHex(id), bson.M{
    "$set": data,
  })

  if err != nil {
    return errors.New(errgo.ErrUpdateUserFailed)
  }

  // update count in department
  UpdateDepartmentsUserCount()

  return nil
}

// 根据id删除用户
func DelUserById(id string) error {
  db, closer, err := mongo.CloneDB()

  if err != nil {
    return err
  } else {
    defer closer()
  }

  // check
  errgo.ErrorIfStringNotObjectId(id, errgo.ErrUserIdError)

  if err = errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return err
  }

  // 删除
  err = db.C(model.UserCollection).UpdateId(bson.ObjectIdHex(id), bson.M{
    "$set": model.User{
      Exist: false,
    },
  })

  if err != nil {
    if err == mgo.ErrNotFound {
      return errors.New(errgo.ErrUserNotFound)
    }
    return err
  }

  err = UpdateDepartmentsUserCount()

  if err != nil {
    return errors.New(errgo.ErrDeleteUserFailed)
  }

  return nil
}

// 根据id查找用户信息
func GetUserInfoById(id string) (gin.H, error) {
  db, closer, err := mongo.CloneDB()

  if err != nil {
    return nil, err
  } else {
    defer closer()
  }

  // check
  errgo.ErrorIfStringNotObjectId(id, errgo.ErrUserIdError)

  if err = errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return nil, err
  }

  data := new(model.User)

  err = db.C(model.UserCollection).Find(bson.M{
    "_id":   bson.ObjectIdHex(id),
    "exist": true,
  }).One(&data)

  if err != nil {
    if err == mgo.ErrNotFound {
      return nil, errors.New(errgo.ErrUserNotFound)
    }
    return nil, err
  }

  return data.GetMap(db), nil
}

// 用户登录并返回用户id
func UserLogin(username string, password string) (id string, err error) {
  db, closer, err := mongo.CloneDB()

  if err != nil {
    return "", err
  } else {
    defer closer()
  }

  // check
  errgo.ErrorIfStringIsEmpty(username, errgo.ErrUsernameEmpty)
  errgo.ErrorIfStringIsEmpty(password, errgo.ErrPasswordEmpty)

  if err = errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return "", err
  }

  data := new(model.User)

  err = db.C(model.UserCollection).Find(bson.M{
    "username": username,
    "password": password,
    "exist":    true,
  }).One(&data)

  if err != nil {
    if err == mgo.ErrNotFound {
      return "", errors.New(errgo.ErrUsernameOrPasswordError)
    }
    return "", err
  } else {
    return data.Id.Hex(), nil
  }
}

// 获取用户列表
func GetUsersList(skip int, limit int, query bson.M) (gin.H, error) {
  db, closer, err := mongo.CloneDB()

  if err != nil {
    return nil, err
  } else {
    defer closer()
  }

  // check
  errgo.ErrorIfIntLessThen(skip, 0, errgo.ErrSkipRange)
  errgo.ErrorIfIntLessThen(limit, 1, errgo.ErrLimitRange)
  errgo.ErrorIfIntMoreThen(limit, 100, errgo.ErrLimitRange)

  if err = errgo.PopError(); err != nil {
    errgo.ClearErrorStack()
    return nil, err
  }

  data := new([]model.User)
  query["exist"] = true

  // find it
  err = db.C(model.UserCollection).Find(query).Skip(skip).Limit(limit).Sort("-createTime").All(data)

  if err != nil {
    if err == mgo.ErrNotFound {
      return nil, errors.New(errgo.ErrUserNotFound)
    }
    return nil, err
  }

  // get count
  count, err := db.C(model.UserCollection).Find(query).Count()

  if err != nil {
    return nil, errors.New(errgo.ErrUserNotFound)
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
